package bridge

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ogormans-deptstack/orca-temporal-bridge/internal/events"
	"go.temporal.io/sdk/client"
)

// EventProvider abstracts event source (ORCA, mock, etc)
type EventProvider interface {
	Subscribe(ctx context.Context) (<-chan *events.BridgeEvent, error)
	Close() error
}

// TemporalClient abstracts Temporal operations
type TemporalClient interface {
	ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error)
	GetWorkflow(ctx context.Context, workflowID string) client.WorkflowRun
	Close() error
}

// Config holds bridge configuration
type Config struct {
	TemporalHost  string
	TemporalPort  int
	Namespace     string
	TaskQueue     string
	EventProvider EventProvider // Can be real ORCA or mock
	BufferSize    int           // events before flush
	FlushInterval time.Duration // max time before flush
}

// Bridge connects ORCA events to Temporal workflows
type Bridge struct {
	config     Config
	provider   EventProvider
	temporal   client.Client
	listener   *EventListener
	flusher    *EventFlusher
	ctx        context.Context
	cancel     context.CancelFunc
	mu         sync.RWMutex
	running    bool
	eventChan  chan *events.BridgeEvent
	translator *events.Translator
}

// NewBridge creates a new bridge
func NewBridge(ctx context.Context, config Config) (*Bridge, error) {
	if config.EventProvider == nil {
		return nil, fmt.Errorf("EventProvider is required")
	}

	if config.BufferSize == 0 {
		config.BufferSize = 1000 // default: 1K events
	}
	if config.FlushInterval == 0 {
		config.FlushInterval = 100 * time.Millisecond // default: 100ms
	}

	bctx, cancel := context.WithCancel(ctx)

	return &Bridge{
		config:     config,
		provider:   config.EventProvider,
		ctx:        bctx,
		cancel:     cancel,
		eventChan:  make(chan *events.BridgeEvent, config.BufferSize),
		translator: events.NewTranslator(),
		running:    false,
	}, nil
}

// Start initializes bridge listeners and flushers
func (b *Bridge) Start() error {
	b.mu.Lock()
	if b.running {
		b.mu.Unlock()
		return fmt.Errorf("bridge is already running")
	}
	b.mu.Unlock()

	// Connect to Temporal
	c, err := client.Dial(client.Options{
		HostPort:  fmt.Sprintf("%s:%d", b.config.TemporalHost, b.config.TemporalPort),
		Namespace: b.config.Namespace,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to Temporal: %w", err)
	}
	b.temporal = c

	// Create event listener (subscribes to provider)
	b.listener = NewEventListener(b.provider, b.eventChan)

	// Create event flusher (batches and sends to Temporal)
	b.flusher = NewEventFlusher(b.temporal, b.eventChan, b.config)

	// Start listener goroutine
	go func() {
		if err := b.listener.Listen(b.ctx); err != nil {
			fmt.Printf("listener error: %v\n", err)
		}
	}()

	// Start flusher goroutine
	go func() {
		if err := b.flusher.Flush(b.ctx); err != nil {
			fmt.Printf("flusher error: %v\n", err)
		}
	}()

	b.mu.Lock()
	b.running = true
	b.mu.Unlock()

	return nil
}

// Stop gracefully shuts down the bridge
func (b *Bridge) Stop() error {
	b.mu.Lock()
	if !b.running {
		b.mu.Unlock()
		return fmt.Errorf("bridge is not running")
	}
	b.mu.Unlock()

	// Signal shutdown
	b.cancel()

	// Give components time to flush
	time.Sleep(500 * time.Millisecond)

	// Close resources
	if b.provider != nil {
		if err := b.provider.Close(); err != nil {
			return fmt.Errorf("failed to close provider: %w", err)
		}
	}

	if b.temporal != nil {
		b.temporal.Close()
	}

	b.mu.Lock()
	b.running = false
	b.mu.Unlock()

	return nil
}

// IsRunning returns whether the bridge is running
func (b *Bridge) IsRunning() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.running
}

// GetStatus returns bridge status information
func (b *Bridge) GetStatus() map[string]interface{} {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return map[string]interface{}{
		"running":      b.running,
		"buffer_size":  len(b.eventChan),
		"max_buffer":   b.config.BufferSize,
		"temporal_url": fmt.Sprintf("%s:%d", b.config.TemporalHost, b.config.TemporalPort),
		"namespace":    b.config.Namespace,
		"task_queue":   b.config.TaskQueue,
	}
}
