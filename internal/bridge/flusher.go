package bridge

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/ogormans-deptstack/orca-temporal-bridge/internal/events"
	"go.temporal.io/sdk/client"
)

// EventFlusher batches events and sends them to Temporal workflows
type EventFlusher struct {
	temporal   client.Client
	eventChan  chan *events.BridgeEvent
	config     Config
	buffer     []*events.BridgeEvent
	bufferMu   chan struct{} // simple sync without full mutex
	eventCount int64
	flushCount int64
	errorCount int64
}

// NewEventFlusher creates a new event flusher
func NewEventFlusher(temporal client.Client, eventChan chan *events.BridgeEvent, config Config) *EventFlusher {
	return &EventFlusher{
		temporal:   temporal,
		eventChan:  eventChan,
		config:     config,
		buffer:     make([]*events.BridgeEvent, 0, config.BufferSize),
		bufferMu:   make(chan struct{}, 1),
		eventCount: 0,
		flushCount: 0,
		errorCount: 0,
	}
}

// Flush periodically flushes buffered events to Temporal
func (f *EventFlusher) Flush(ctx context.Context) error {
	ticker := time.NewTicker(f.config.FlushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Final flush on shutdown
			_ = f.flushBuffer(ctx)
			return ctx.Err()

		case event := <-f.eventChan:
			if event == nil {
				continue
			}

			// Add to buffer
			f.addToBuffer(event)

			// Check if buffer is full
			if f.bufferLen() >= f.config.BufferSize {
				if err := f.flushBuffer(ctx); err != nil {
					atomic.AddInt64(&f.errorCount, 1)
				}
			}

		case <-ticker.C:
			// Time-based flush
			if f.bufferLen() > 0 {
				if err := f.flushBuffer(ctx); err != nil {
					atomic.AddInt64(&f.errorCount, 1)
				}
			}
		}
	}
}

// addToBuffer adds an event to the buffer (thread-safe)
func (f *EventFlusher) addToBuffer(event *events.BridgeEvent) {
	f.bufferMu <- struct{}{}        // acquire
	defer func() { <-f.bufferMu }() // release

	f.buffer = append(f.buffer, event)
	atomic.AddInt64(&f.eventCount, 1)
}

// bufferLen returns buffer length (simplified, may have slight race)
func (f *EventFlusher) bufferLen() int {
	f.bufferMu <- struct{}{}
	defer func() { <-f.bufferMu }()
	return len(f.buffer)
}

// flushBuffer sends buffered events to Temporal
func (f *EventFlusher) flushBuffer(ctx context.Context) error {
	f.bufferMu <- struct{}{} // acquire
	if len(f.buffer) == 0 {
		<-f.bufferMu // release
		return nil
	}

	// Copy buffer for processing
	toFlush := make([]*events.BridgeEvent, len(f.buffer))
	copy(toFlush, f.buffer)
	f.buffer = f.buffer[:0] // clear buffer
	<-f.bufferMu            // release

	// Process each event
	for _, event := range toFlush {
		if err := f.sendToTemporal(ctx, event); err != nil {
			return err
		}
	}

	atomic.AddInt64(&f.flushCount, 1)
	return nil
}

// sendToTemporal sends a single event to Temporal via workflow
func (f *EventFlusher) sendToTemporal(ctx context.Context, event *events.BridgeEvent) error {
	// Determine workflow type based on event type
	workflowID := f.generateWorkflowID(event)
	workflowType := f.getWorkflowType(event)

	// Temporal workflow options
	options := client.StartWorkflowOptions{
		ID:          workflowID,
		TaskQueue:   f.config.TaskQueue,
		RetryPolicy: nil, // TODO: configure retry policy
	}

	// Start workflow (will be coordinator or worker based on event type)
	run, err := f.temporal.ExecuteWorkflow(ctx, options, workflowType, event)
	if err != nil {
		return fmt.Errorf("failed to start workflow %s: %w", workflowID, err)
	}

	// Log workflow started
	_ = run // use run to avoid unused variable

	return nil
}

// generateWorkflowID creates a unique workflow ID from event data
func (f *EventFlusher) generateWorkflowID(event *events.BridgeEvent) string {
	// Format: event_type:agent_id:timestamp
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s:%s:%d:%s",
		event.EventType,
		event.AgentID,
		timestamp,
		uuid.New().String()[:8],
	)
}

// getWorkflowType returns the Temporal workflow name for an event
func (f *EventFlusher) getWorkflowType(event *events.BridgeEvent) string {
	// For now, all events start a generic orchestration workflow
	// Later: specialized workflows per event type
	switch event.EventType {
	case events.EventCoordinatorStarted, events.EventCoordinatorCompleted:
		return "CoordinatorWorkflow"
	case events.EventWorkerSpawned, events.EventWorkerCompleted, events.EventWorkerFailed:
		return "WorkerWorkflow"
	default:
		return "EventWorkflow"
	}
}

// GetStats returns flusher statistics
func (f *EventFlusher) GetStats() map[string]interface{} {
	f.bufferMu <- struct{}{}
	bufferLen := len(f.buffer)
	<-f.bufferMu

	return map[string]interface{}{
		"events_processed": atomic.LoadInt64(&f.eventCount),
		"flushes":          atomic.LoadInt64(&f.flushCount),
		"errors":           atomic.LoadInt64(&f.errorCount),
		"buffer_size":      bufferLen,
	}
}
