package events

import "time"

// Event source constants
const (
	SourceORCA     = "orca"
	SourceTemporal = "temporal"
)

// ORCA event type constants
const (
	EventCoordinatorStarted   = "coordinator_started"
	EventWorkerSpawned        = "worker_spawned"
	EventAgentDecisionMade    = "agent_decision_made"
	EventWorkerCompleted      = "worker_completed"
	EventWorkerFailed         = "worker_failed"
	EventCoordinatorCompleted = "coordinator_completed"
	EventAgentError           = "agent_error"
)

// BridgeEvent is the unified event type for ORCA and Temporal events.
type BridgeEvent struct {
	ID            string      `json:"id"`
	Timestamp     time.Time   `json:"timestamp"`
	Source        string      `json:"source"`
	EventType     string      `json:"event_type"`
	WorktreeID    string      `json:"worktree_id,omitempty"`
	AgentID       string      `json:"agent_id,omitempty"`
	AgentRole     string      `json:"agent_role,omitempty"`
	PreviousState string      `json:"previous_state,omitempty"`
	CurrentState  string      `json:"current_state,omitempty"`
	Payload       interface{} `json:"payload,omitempty"`
	ErrorMessage  string      `json:"error_message,omitempty"`
	ErrorCode     string      `json:"error_code,omitempty"`
}

// AgentStartedPayload: Coordinator or worker agent started
type AgentStartedPayload struct {
	WorkflowID string            `json:"workflow_id"`
	TaskID     string            `json:"task_id"`
	Metadata   map[string]string `json:"metadata"`
}

// AgentDecisionPayload: Agent made a decision (prompt/response)
type AgentDecisionPayload struct {
	DecisionID string `json:"decision_id"`
	Prompt     string `json:"prompt"`
	Response   string `json:"response"`
	DurationMs int64  `json:"duration_ms"`
}

// AgentCompletedPayload: Agent task completed
type AgentCompletedPayload struct {
	WorkflowID string `json:"workflow_id"`
	TaskID     string `json:"task_id"`
	Result     []byte `json:"result"`
	Success    bool   `json:"success"`
}

// AgentFailedPayload: Agent task failed
type AgentFailedPayload struct {
	WorkflowID string `json:"workflow_id"`
	TaskID     string `json:"task_id"`
	Error      string `json:"error"`
	Attempt    int32  `json:"attempt"`
}

// MockOrcaEvent is a helper for testing (simulates ORCA event structure)
type MockOrcaEvent struct {
	EventType     string            `json:"event_type"`
	WorktreeID    string            `json:"worktree_id,omitempty"`
	AgentID       string            `json:"agent_id,omitempty"`
	AgentRole     string            `json:"agent_role,omitempty"`
	PreviousState string            `json:"previous_state,omitempty"`
	CurrentState  string            `json:"current_state,omitempty"`
	WorkflowID    string            `json:"workflow_id,omitempty"`
	TaskID        string            `json:"task_id,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	DecisionID    string            `json:"decision_id,omitempty"`
	Prompt        string            `json:"prompt,omitempty"`
	Response      string            `json:"response,omitempty"`
	DurationMs    int64             `json:"duration_ms,omitempty"`
	Result        []byte            `json:"result,omitempty"`
	Success       bool              `json:"success,omitempty"`
	Error         string            `json:"error,omitempty"`
	ErrorCode     string            `json:"error_code,omitempty"`
	Attempt       int32             `json:"attempt,omitempty"`
}

// ToMap converts MockOrcaEvent to map[string]interface{} for testing translator
func (m *MockOrcaEvent) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
	result["event_type"] = m.EventType
	if m.WorktreeID != "" {
		result["worktree_id"] = m.WorktreeID
	}
	if m.AgentID != "" {
		result["agent_id"] = m.AgentID
	}
	if m.AgentRole != "" {
		result["agent_role"] = m.AgentRole
	}
	if m.PreviousState != "" {
		result["previous_state"] = m.PreviousState
	}
	if m.CurrentState != "" {
		result["current_state"] = m.CurrentState
	}
	if m.WorkflowID != "" {
		result["workflow_id"] = m.WorkflowID
	}
	if m.TaskID != "" {
		result["task_id"] = m.TaskID
	}
	if len(m.Metadata) > 0 {
		result["metadata"] = m.Metadata
	}
	if m.DecisionID != "" {
		result["decision_id"] = m.DecisionID
	}
	if m.Prompt != "" {
		result["prompt"] = m.Prompt
	}
	if m.Response != "" {
		result["response"] = m.Response
	}
	if m.DurationMs > 0 {
		result["duration_ms"] = m.DurationMs
	}
	if len(m.Result) > 0 {
		result["result"] = m.Result
	}
	if m.Success {
		result["success"] = m.Success
	}
	if m.Error != "" {
		result["error"] = m.Error
	}
	if m.ErrorCode != "" {
		result["error_code"] = m.ErrorCode
	}
	if m.Attempt > 0 {
		result["attempt"] = m.Attempt
	}
	return result
}
