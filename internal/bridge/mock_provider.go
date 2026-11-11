package bridge

import (
	"context"
	"sync"
	"time"

	"github.com/ogormans-deptstack/orca-temporal-bridge/internal/events"
)

// MockEventProvider emits test events for E2E testing
type MockEventProvider struct {
	events chan *events.BridgeEvent
	closed bool
	mu     sync.Mutex
}

// NewMockEventProvider creates a new mock event provider
func NewMockEventProvider() *MockEventProvider {
	return &MockEventProvider{
		events: make(chan *events.BridgeEvent, 1000),
		closed: false,
	}
}

// Subscribe returns the event channel (implements EventProvider interface)
func (m *MockEventProvider) Subscribe(ctx context.Context) (<-chan *events.BridgeEvent, error) {
	// Return channel that respects context cancellation
	resultChan := make(chan *events.BridgeEvent)
	go func() {
		defer close(resultChan)
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-m.events:
				if !ok {
					return
				}
				select {
				case resultChan <- event:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return resultChan, nil
}

// Close closes the event channel
func (m *MockEventProvider) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.closed {
		close(m.events)
		m.closed = true
	}
	return nil
}

// Emit sends a test event
func (m *MockEventProvider) Emit(event *events.BridgeEvent) error {
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return nil // silently drop if closed
	}
	m.mu.Unlock()

	select {
	case m.events <- event:
		return nil
	default:
		return nil // drop if buffer full (non-blocking)
	}
}

// EmitCoordinatorStarted is a convenience method for testing
func (m *MockEventProvider) EmitCoordinatorStarted(worktreeID, agentID string, metadata map[string]string) error {
	event := &events.BridgeEvent{
		ID:         "test-" + agentID,
		Timestamp:  time.Now(),
		Source:     events.SourceORCA,
		EventType:  events.EventCoordinatorStarted,
		WorktreeID: worktreeID,
		AgentID:    agentID,
		AgentRole:  "coordinator",
		Payload: &events.AgentStartedPayload{
			WorkflowID: agentID + "-workflow",
			TaskID:     agentID + "-task",
			Metadata:   metadata,
		},
	}
	return m.Emit(event)
}

// EmitWorkerSpawned is a convenience method for testing
func (m *MockEventProvider) EmitWorkerSpawned(worktreeID, agentID, parentWorkflowID string) error {
	event := &events.BridgeEvent{
		ID:         "test-" + agentID,
		Timestamp:  time.Now(),
		Source:     events.SourceORCA,
		EventType:  events.EventWorkerSpawned,
		WorktreeID: worktreeID,
		AgentID:    agentID,
		AgentRole:  "worker",
		Payload: &events.AgentStartedPayload{
			WorkflowID: parentWorkflowID,
			TaskID:     agentID + "-task",
			Metadata:   make(map[string]string),
		},
	}
	return m.Emit(event)
}

// EmitAgentDecisionMade is a convenience method for testing
func (m *MockEventProvider) EmitAgentDecisionMade(agentID, prompt, response string, durationMs int64) error {
	event := &events.BridgeEvent{
		ID:        "test-" + agentID,
		Timestamp: time.Now(),
		Source:    events.SourceORCA,
		EventType: events.EventAgentDecisionMade,
		AgentID:   agentID,
		Payload: &events.AgentDecisionPayload{
			DecisionID: agentID + "-decision",
			Prompt:     prompt,
			Response:   response,
			DurationMs: durationMs,
		},
	}
	return m.Emit(event)
}

// EmitWorkerCompleted is a convenience method for testing
func (m *MockEventProvider) EmitWorkerCompleted(agentID string, success bool) error {
	event := &events.BridgeEvent{
		ID:        "test-" + agentID,
		Timestamp: time.Now(),
		Source:    events.SourceORCA,
		EventType: events.EventWorkerCompleted,
		AgentID:   agentID,
		Payload: &events.AgentCompletedPayload{
			WorkflowID: agentID + "-workflow",
			TaskID:     agentID + "-task",
			Result:     []byte("test result"),
			Success:    success,
		},
	}
	return m.Emit(event)
}

// EmitWorkerFailed is a convenience method for testing
func (m *MockEventProvider) EmitWorkerFailed(agentID, errorMsg string, attempt int32) error {
	event := &events.BridgeEvent{
		ID:           "test-" + agentID,
		Timestamp:    time.Now(),
		Source:       events.SourceORCA,
		EventType:    events.EventWorkerFailed,
		AgentID:      agentID,
		ErrorMessage: errorMsg,
		Payload: &events.AgentFailedPayload{
			WorkflowID: agentID + "-workflow",
			TaskID:     agentID + "-task",
			Error:      errorMsg,
			Attempt:    attempt,
		},
	}
	return m.Emit(event)
}
