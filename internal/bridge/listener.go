package bridge

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/ogormans-deptstack/orca-temporal-bridge/internal/events"
)

// EventListener subscribes to ORCA events and sends them for processing
type EventListener struct {
	provider   EventProvider
	eventChan  chan *events.BridgeEvent
	eventCount int64
	errorCount int64
}

// NewEventListener creates a new event listener
func NewEventListener(provider EventProvider, eventChan chan *events.BridgeEvent) *EventListener {
	return &EventListener{
		provider:   provider,
		eventChan:  eventChan,
		eventCount: 0,
		errorCount: 0,
	}
}

// Listen subscribes to events and forwards them
func (l *EventListener) Listen(ctx context.Context) error {
	eventStream, err := l.provider.Subscribe(ctx)
	if err != nil {
		return fmt.Errorf("failed to subscribe to events: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case event, ok := <-eventStream:
			if !ok {
				// Stream closed
				return nil
			}

			if event == nil {
				atomic.AddInt64(&l.errorCount, 1)
				continue
			}

			// Forward to flusher (non-blocking with buffer)
			select {
			case l.eventChan <- event:
				atomic.AddInt64(&l.eventCount, 1)
			case <-ctx.Done():
				return ctx.Err()
			default:
				// Buffer full, drop event (production: could implement backpressure)
				atomic.AddInt64(&l.errorCount, 1)
			}
		}
	}
}

// GetStats returns listener statistics
func (l *EventListener) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"events_received": atomic.LoadInt64(&l.eventCount),
		"errors":          atomic.LoadInt64(&l.errorCount),
	}
}
