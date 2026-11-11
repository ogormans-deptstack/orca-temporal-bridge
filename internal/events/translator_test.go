package events

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTranslator(t *testing.T) {
	tr := NewTranslator()
	assert.NotNil(t, tr)
}

// ============================================================================
// EventCoordinatorStarted Tests
// ============================================================================

func TestTranslateOrcaEvent_CoordinatorStarted(t *testing.T) {
	tr := NewTranslator()
	mockEvent := &MockOrcaEvent{
		EventType:  EventCoordinatorStarted,
		WorktreeID: "wt-123",
		AgentID:    "coord-1",
		AgentRole:  "coordinator",
		WorkflowID: "wf-456",
		TaskID:     "task-789",
		Metadata: map[string]string{
			"version": "1.0",
			"env":     "test",
		},
	}

	event, err := tr.TranslateOrcaEvent(mockEvent.ToMap())
	require.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, EventCoordinatorStarted, event.EventType)
	assert.Equal(t, SourceORCA, event.Source)
	assert.Equal(t, "wt-123", event.WorktreeID)
	assert.Equal(t, "coord-1", event.AgentID)
	assert.Equal(t, "coordinator", event.AgentRole)

	payload := event.Payload.(*AgentStartedPayload)
	assert.Equal(t, "wf-456", payload.WorkflowID)
	assert.Equal(t, "task-789", payload.TaskID)
	assert.Equal(t, "1.0", payload.Metadata["version"])
	assert.Equal(t, "test", payload.Metadata["env"])
}

func TestTranslateOrcaEvent_CoordinatorStarted_MinimalFields(t *testing.T) {
	tr := NewTranslator()
	mockEvent := &MockOrcaEvent{
		EventType: EventCoordinatorStarted,
	}

	event, err := tr.TranslateOrcaEvent(mockEvent.ToMap())
	require.NoError(t, err)
	assert.Equal(t, EventCoordinatorStarted, event.EventType)

	payload := event.Payload.(*AgentStartedPayload)
	assert.Equal(t, "", payload.WorkflowID)
	assert.Equal(t, "", payload.TaskID)
}

// ============================================================================
// EventWorkerSpawned Tests
// ============================================================================

func TestTranslateOrcaEvent_WorkerSpawned(t *testing.T) {
	tr := NewTranslator()
	mockEvent := &MockOrcaEvent{
		EventType:  EventWorkerSpawned,
		WorktreeID: "wt-123",
		AgentID:    "worker-1",
		AgentRole:  "worker",
		WorkflowID: "parent-wf",
		TaskID:     "worker-task-1",
		Metadata: map[string]string{
			"index": "0",
			"type":  "research",
		},
	}

	event, err := tr.TranslateOrcaEvent(mockEvent.ToMap())
	require.NoError(t, err)
	assert.Equal(t, EventWorkerSpawned, event.EventType)
	assert.Equal(t, "worker-1", event.AgentID)
	assert.Equal(t, "worker", event.AgentRole)

	payload := event.Payload.(*AgentStartedPayload)
	assert.Equal(t, "parent-wf", payload.WorkflowID)
	assert.Equal(t, "worker-task-1", payload.TaskID)
	assert.Equal(t, "research", payload.Metadata["type"])
}

func TestTranslateOrcaEvent_WorkerSpawned_NoMetadata(t *testing.T) {
	tr := NewTranslator()
	mockEvent := &MockOrcaEvent{
		EventType:  EventWorkerSpawned,
		WorkflowID: "parent-wf",
		TaskID:     "worker-task-1",
	}

	event, err := tr.TranslateOrcaEvent(mockEvent.ToMap())
	require.NoError(t, err)

	payload := event.Payload.(*AgentStartedPayload)
	assert.NotNil(t, payload.Metadata)
	assert.Equal(t, 0, len(payload.Metadata))
}

// ============================================================================
// EventAgentDecisionMade Tests
// ============================================================================

func TestTranslateOrcaEvent_AgentDecisionMade(t *testing.T) {
	tr := NewTranslator()
	mockEvent := &MockOrcaEvent{
		EventType:  EventAgentDecisionMade,
		WorktreeID: "wt-123",
		AgentID:    "worker-1",
		DecisionID: "decision-1",
		Prompt:     "What is 2+2?",
		Response:   "4",
		DurationMs: 1500,
	}

	event, err := tr.TranslateOrcaEvent(mockEvent.ToMap())
	require.NoError(t, err)
	assert.Equal(t, EventAgentDecisionMade, event.EventType)

	payload := event.Payload.(*AgentDecisionPayload)
	assert.Equal(t, "decision-1", payload.DecisionID)
	assert.Equal(t, "What is 2+2?", payload.Prompt)
	assert.Equal(t, "4", payload.Response)
	assert.Equal(t, int64(1500), payload.DurationMs)
}

func TestTranslateOrcaEvent_AgentDecisionMade_VariantEventType(t *testing.T) {
	tr := NewTranslator()
	mockEvent := &MockOrcaEvent{
		EventType:  "decision_made", // variant spelling
		DecisionID: "decision-2",
		Prompt:     "Test prompt",
		Response:   "Test response",
		DurationMs: 1000,
	}

	event, err := tr.TranslateOrcaEvent(mockEvent.ToMap())
	require.NoError(t, err)
	assert.Equal(t, EventAgentDecisionMade, event.EventType)

	payload := event.Payload.(*AgentDecisionPayload)
	assert.Equal(t, "decision-2", payload.DecisionID)
}

func TestTranslateOrcaEvent_AgentDecisionMade_ZeroDuration(t *testing.T) {
	tr := NewTranslator()
	mockEvent := &MockOrcaEvent{
		EventType:  EventAgentDecisionMade,
		DecisionID: "decision-3",
		Prompt:     "prompt",
		Response:   "response",
	}

	event, err := tr.TranslateOrcaEvent(mockEvent.ToMap())
	require.NoError(t, err)

	payload := event.Payload.(*AgentDecisionPayload)
	assert.Equal(t, int64(0), payload.DurationMs)
}

// ============================================================================
// EventWorkerCompleted Tests
// ============================================================================

func TestTranslateOrcaEvent_WorkerCompleted(t *testing.T) {
	tr := NewTranslator()
	mockEvent := &MockOrcaEvent{
		EventType:  EventWorkerCompleted,
		WorktreeID: "wt-123",
		AgentID:    "worker-1",
		WorkflowID: "parent-wf",
		TaskID:     "worker-task-1",
		Result:     []byte("success result"),
		Success:    true,
	}

	event, err := tr.TranslateOrcaEvent(mockEvent.ToMap())
	require.NoError(t, err)
	assert.Equal(t, EventWorkerCompleted, event.EventType)

	payload := event.Payload.(*AgentCompletedPayload)
	assert.Equal(t, "parent-wf", payload.WorkflowID)
	assert.Equal(t, "worker-task-1", payload.TaskID)
	assert.Equal(t, []byte("success result"), payload.Result)
	assert.True(t, payload.Success)
}

func TestTranslateOrcaEvent_WorkerCompleted_Variant(t *testing.T) {
	tr := NewTranslator()
	mockEvent := &MockOrcaEvent{
		EventType:  "completed", // variant
		WorkflowID: "parent-wf",
		TaskID:     "worker-task-1",
		Result:     []byte("data"),
		Success:    true,
	}

	event, err := tr.TranslateOrcaEvent(mockEvent.ToMap())
	require.NoError(t, err)
	assert.Equal(t, EventWorkerCompleted, event.EventType)
}

func TestTranslateOrcaEvent_WorkerCompleted_NoResult(t *testing.T) {
	tr := NewTranslator()
	mockEvent := &MockOrcaEvent{
		EventType:  EventWorkerCompleted,
		WorkflowID: "parent-wf",
		TaskID:     "worker-task-1",
		Success:    true,
	}

	event, err := tr.TranslateOrcaEvent(mockEvent.ToMap())
	require.NoError(t, err)

	payload := event.Payload.(*AgentCompletedPayload)
	assert.Nil(t, payload.Result)
	assert.True(t, payload.Success)
}

func TestTranslateOrcaEvent_WorkerCompleted_FailedTask(t *testing.T) {
	tr := NewTranslator()
	mockEvent := &MockOrcaEvent{
		EventType:  EventWorkerCompleted,
		WorkflowID: "parent-wf",
		TaskID:     "worker-task-1",
		Result:     []byte("partial result"),
		Success:    false,
	}

	event, err := tr.TranslateOrcaEvent(mockEvent.ToMap())
	require.NoError(t, err)

	payload := event.Payload.(*AgentCompletedPayload)
	assert.False(t, payload.Success)
}

// ============================================================================
// EventWorkerFailed Tests
// ============================================================================

func TestTranslateOrcaEvent_WorkerFailed(t *testing.T) {
	tr := NewTranslator()
	mockEvent := &MockOrcaEvent{
		EventType:  EventWorkerFailed,
		WorktreeID: "wt-123",
		AgentID:    "worker-1",
		WorkflowID: "parent-wf",
		TaskID:     "worker-task-1",
		Error:      "timeout exceeded",
		ErrorCode:  "ERR_TIMEOUT",
		Attempt:    2,
	}

	event, err := tr.TranslateOrcaEvent(mockEvent.ToMap())
	require.NoError(t, err)
	assert.Equal(t, EventWorkerFailed, event.EventType)
	assert.Equal(t, "timeout exceeded", event.ErrorMessage)
	assert.Equal(t, "ERR_TIMEOUT", event.ErrorCode)

	payload := event.Payload.(*AgentFailedPayload)
	assert.Equal(t, "parent-wf", payload.WorkflowID)
	assert.Equal(t, "worker-task-1", payload.TaskID)
	assert.Equal(t, "timeout exceeded", payload.Error)
	assert.Equal(t, int32(2), payload.Attempt)
}

func TestTranslateOrcaEvent_WorkerFailed_Variant(t *testing.T) {
	tr := NewTranslator()
	mockEvent := &MockOrcaEvent{
		EventType:  "failed", // variant
		WorkflowID: "parent-wf",
		TaskID:     "worker-task-1",
		Error:      "some error",
		Attempt:    1,
	}

	event, err := tr.TranslateOrcaEvent(mockEvent.ToMap())
	require.NoError(t, err)
	assert.Equal(t, EventWorkerFailed, event.EventType)
}

func TestTranslateOrcaEvent_WorkerFailed_NoErrorCode(t *testing.T) {
	tr := NewTranslator()
	mockEvent := &MockOrcaEvent{
		EventType:  EventWorkerFailed,
		WorkflowID: "parent-wf",
		TaskID:     "worker-task-1",
		Error:      "generic error",
		Attempt:    1,
	}

	event, err := tr.TranslateOrcaEvent(mockEvent.ToMap())
	require.NoError(t, err)

	assert.Equal(t, "generic error", event.ErrorMessage)
	assert.Equal(t, "", event.ErrorCode)
}

func TestTranslateOrcaEvent_WorkerFailed_MultipleRetries(t *testing.T) {
	tr := NewTranslator()
	mockEvent := &MockOrcaEvent{
		EventType:  EventWorkerFailed,
		WorkflowID: "parent-wf",
		TaskID:     "worker-task-1",
		Error:      "error",
		Attempt:    5,
	}

	event, err := tr.TranslateOrcaEvent(mockEvent.ToMap())
	require.NoError(t, err)

	payload := event.Payload.(*AgentFailedPayload)
	assert.Equal(t, int32(5), payload.Attempt)
}

// ============================================================================
// EventCoordinatorCompleted Tests
// ============================================================================

func TestTranslateOrcaEvent_CoordinatorCompleted(t *testing.T) {
	tr := NewTranslator()
	mockEvent := &MockOrcaEvent{
		EventType:  EventCoordinatorCompleted,
		WorktreeID: "wt-123",
		AgentID:    "coord-1",
		WorkflowID: "wf-456",
		TaskID:     "coord-task-1",
		Result:     []byte("final result"),
		Success:    true,
	}

	event, err := tr.TranslateOrcaEvent(mockEvent.ToMap())
	require.NoError(t, err)
	assert.Equal(t, EventCoordinatorCompleted, event.EventType)

	payload := event.Payload.(*AgentCompletedPayload)
	assert.Equal(t, "wf-456", payload.WorkflowID)
	assert.Equal(t, "coord-task-1", payload.TaskID)
	assert.Equal(t, []byte("final result"), payload.Result)
	assert.True(t, payload.Success)
}

func TestTranslateOrcaEvent_CoordinatorCompleted_Failed(t *testing.T) {
	tr := NewTranslator()
	mockEvent := &MockOrcaEvent{
		EventType:  EventCoordinatorCompleted,
		WorkflowID: "wf-456",
		TaskID:     "coord-task-1",
		Result:     nil,
		Success:    false,
	}

	event, err := tr.TranslateOrcaEvent(mockEvent.ToMap())
	require.NoError(t, err)

	payload := event.Payload.(*AgentCompletedPayload)
	assert.False(t, payload.Success)
}

// ============================================================================
// EventAgentError Tests
// ============================================================================

func TestTranslateOrcaEvent_AgentError(t *testing.T) {
	tr := NewTranslator()
	mockEvent := &MockOrcaEvent{
		EventType: EventAgentError,
		AgentID:   "worker-1",
		Error:     "agent crashed",
		ErrorCode: "ERR_CRASH",
	}

	event, err := tr.TranslateOrcaEvent(mockEvent.ToMap())
	require.NoError(t, err)
	assert.Equal(t, EventAgentError, event.EventType)
	assert.Equal(t, "agent crashed", event.ErrorMessage)
	assert.Equal(t, "ERR_CRASH", event.ErrorCode)
}

func TestTranslateOrcaEvent_AgentError_Variant(t *testing.T) {
	tr := NewTranslator()
	mockEvent := &MockOrcaEvent{
		EventType: "error", // variant
		Error:     "some error",
	}

	event, err := tr.TranslateOrcaEvent(mockEvent.ToMap())
	require.NoError(t, err)
	assert.Equal(t, EventAgentError, event.EventType)
}

// ============================================================================
// Unknown Event Type Tests
// ============================================================================

func TestTranslateOrcaEvent_UnknownEventType(t *testing.T) {
	tr := NewTranslator()
	mockEvent := &MockOrcaEvent{
		EventType: "unknown_event_type",
		AgentID:   "worker-1",
	}

	event, err := tr.TranslateOrcaEvent(mockEvent.ToMap())
	require.NoError(t, err)
	assert.Equal(t, "unknown_event_type", event.EventType)
	// Payload should be the original map
	assert.IsType(t, map[string]interface{}{}, event.Payload)
}

// ============================================================================
// Nil and Edge Case Tests
// ============================================================================

func TestTranslateOrcaEvent_NilEvent(t *testing.T) {
	tr := NewTranslator()
	event, err := tr.TranslateOrcaEvent(nil)
	assert.Error(t, err)
	assert.Nil(t, event)
	assert.Equal(t, "ORCA event is nil", err.Error())
}

func TestTranslateOrcaEvent_NonMapEvent(t *testing.T) {
	tr := NewTranslator()
	event, err := tr.TranslateOrcaEvent("string event")
	require.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, SourceORCA, event.Source)
	assert.Equal(t, "unknown_orca_event", event.EventType)
	assert.Equal(t, "string event", event.Payload)
}

func TestTranslateOrcaEvent_IntEvent(t *testing.T) {
	tr := NewTranslator()
	event, err := tr.TranslateOrcaEvent(42)
	require.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, 42, event.Payload)
}

// ============================================================================
// TranslateTemporalEvent Tests
// ============================================================================

func TestTranslateTemporalEvent(t *testing.T) {
	tr := NewTranslator()
	temporalEvent := map[string]interface{}{"workflow_id": "wf-123"}

	event, err := tr.TranslateTemporalEvent(temporalEvent)
	require.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, SourceTemporal, event.Source)
	assert.Equal(t, "unknown_temporal_event", event.EventType)
	assert.Equal(t, temporalEvent, event.Payload)
}

func TestTranslateTemporalEvent_NilEvent(t *testing.T) {
	tr := NewTranslator()
	event, err := tr.TranslateTemporalEvent(nil)
	assert.Error(t, err)
	assert.Nil(t, event)
	assert.Equal(t, "temporal event is nil", err.Error())
}

// ============================================================================
// ValidateEvent Tests
// ============================================================================

func TestValidateEvent_Valid(t *testing.T) {
	event := &BridgeEvent{
		ID:        "evt-123",
		Timestamp: time.Now(),
		Source:    SourceORCA,
		EventType: EventCoordinatorStarted,
	}

	tr := NewTranslator()
	err := tr.ValidateEvent(event)
	assert.NoError(t, err)
}

func TestValidateEvent_NilEvent(t *testing.T) {
	tr := NewTranslator()
	err := tr.ValidateEvent(nil)
	assert.Error(t, err)
	assert.Equal(t, "event is nil", err.Error())
}

func TestValidateEvent_NoID(t *testing.T) {
	event := &BridgeEvent{
		Timestamp: time.Now(),
		Source:    SourceORCA,
		EventType: EventCoordinatorStarted,
	}

	tr := NewTranslator()
	err := tr.ValidateEvent(event)
	assert.Error(t, err)
	assert.Equal(t, "event ID is required", err.Error())
}

func TestValidateEvent_InvalidSource(t *testing.T) {
	event := &BridgeEvent{
		ID:        "evt-123",
		Timestamp: time.Now(),
		Source:    "invalid",
		EventType: EventCoordinatorStarted,
	}

	tr := NewTranslator()
	err := tr.ValidateEvent(event)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid event source")
}

func TestValidateEvent_NoEventType(t *testing.T) {
	event := &BridgeEvent{
		ID:        "evt-123",
		Timestamp: time.Now(),
		Source:    SourceORCA,
	}

	tr := NewTranslator()
	err := tr.ValidateEvent(event)
	assert.Error(t, err)
	assert.Equal(t, "event type is required", err.Error())
}

func TestValidateEvent_NoTimestamp(t *testing.T) {
	event := &BridgeEvent{
		ID:        "evt-123",
		Source:    SourceORCA,
		EventType: EventCoordinatorStarted,
	}

	tr := NewTranslator()
	err := tr.ValidateEvent(event)
	assert.Error(t, err)
	assert.Equal(t, "event timestamp is required", err.Error())
}

// ============================================================================
// State Transition Tests
// ============================================================================

func TestTranslateOrcaEvent_StateTransition(t *testing.T) {
	tr := NewTranslator()
	mockEvent := &MockOrcaEvent{
		EventType:     EventWorkerCompleted,
		PreviousState: "running",
		CurrentState:  "completed",
		WorkflowID:    "wf-1",
		TaskID:        "task-1",
		Success:       true,
	}

	event, err := tr.TranslateOrcaEvent(mockEvent.ToMap())
	require.NoError(t, err)
	assert.Equal(t, "running", event.PreviousState)
	assert.Equal(t, "completed", event.CurrentState)
}

// ============================================================================
// Helper Functions Tests
// ============================================================================

func TestExtractString_MultipleKeys(t *testing.T) {
	m := map[string]interface{}{"alt_key": "value"}

	// Test with fallback keys
	result := extractString(m, "primary_key", "alt_key", "default_key")
	assert.Equal(t, "value", result)
}

func TestExtractInt64_TypeConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected int64
	}{
		{"int64", int64(100), 100},
		{"int", int(50), 50},
		{"float64", float64(75.5), 75},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := map[string]interface{}{"value": tt.input}
			result := extractInt64(m, "value")
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractBytes_StringConversion(t *testing.T) {
	m := map[string]interface{}{"data": "test string"}
	result := extractBytes(m, "data")
	assert.Equal(t, []byte("test string"), result)
}

func TestExtractMetadata_NestedMap(t *testing.T) {
	m := map[string]interface{}{
		"metadata": map[string]interface{}{
			"key1": "value1",
			"key2": 123, // Should be converted to string
		},
	}

	result := extractMetadata(m)
	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, "123", result["key2"])
}
