// Package bridge provides the public SDK for integrating the ORCA-Temporal bridge into other applications.
package bridge

import (
	"context"
	"fmt"
)

// Client is the public API for using the ORCA-Temporal bridge.
type Client interface {
	// OnOrcaEvent processes an incoming ORCA event
	OnOrcaEvent(ctx context.Context, event interface{}) error

	// GetWorkflowHistory retrieves the execution history for a workflow
	GetWorkflowHistory(ctx context.Context, workflowID string) ([]interface{}, error)

	// ReplayWorkflow replays agent decisions from workflow history
	ReplayWorkflow(ctx context.Context, workflowID string) ([]interface{}, error)

	// Status returns the current bridge status
	Status(ctx context.Context) (Status, error)

	// Close closes the client connection
	Close(ctx context.Context) error
}

// Status represents the current bridge status.
type Status struct {
	Running            bool
	TemporalConnected  bool
	EventsReceived     int64
	EventsPersisted    int64
	Errors             int64
	LastEventTimestamp int64
}

// NewClient creates a new bridge client.
func NewClient(ctx context.Context, addr string) (Client, error) {
	// TODO: Implement client factory
	// - Parse address (host:port)
	// - Connect to Temporal
	// - Initialize internal bridge

	return nil, fmt.Errorf("not yet implemented")
}

// Config holds client configuration.
type Config struct {
	TemporalHost    string
	TemporalPort    string
	Namespace       string
	BufferSize      int
	FlushIntervalMs int
}
