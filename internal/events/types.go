// Package events provides event schema and translation between ORCA and Temporal formats.
package events

import (
	"time"
)

// BridgeEvent is the unified event type representing both ORCA and Temporal events.
type BridgeEvent struct {
	ID            string      `json:"id"`
	Timestamp     time.Time   `json:"timestamp"`
	Source        string      `json:"source"` // "orca" or "temporal"
	EventType     string      `json:"event_type"`
	WorktreeID    string      `json:"worktree_id"`
	AgentID       string      `json:"agent_id"`
	AgentRole     string      `json:"agent_role"` // "coordinator" or "worker"
	PreviousState string      `json:"previous_state"`
	CurrentState  string      `json:"current_state"`
	Payload       interface{} `json:"payload"`
}

// EventType constants
const (
	EventCoordinatorStarted   = "coordinator_started"
	EventWorkerSpawned        = "worker_spawned"
	EventAgentDecisionMade    = "agent_decision_made"
	EventWorkerCompleted      = "worker_completed"
	EventWorkerFailed         = "worker_failed"
	EventCoordinatorCompleted = "coordinator_completed"
	EventAgentError           = "agent_error"
)

// AgentRole constants
const (
	RoleCoordinator = "coordinator"
	RoleWorker      = "worker"
)

// EventSource constants
const (
	SourceORCA     = "orca"
	SourceTemporal = "temporal"
)
