package events

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Translator converts between ORCA and Temporal event formats.
type Translator struct {
	// TODO: Add configuration if needed
}

// NewTranslator creates a new event translator.
func NewTranslator() *Translator {
	return &Translator{}
}

// TranslateOrcaEvent converts an ORCA event to a BridgeEvent.
// Handles all known ORCA event types.
func (t *Translator) TranslateOrcaEvent(orcaEvent interface{}) (*BridgeEvent, error) {
	if orcaEvent == nil {
		return nil, fmt.Errorf("ORCA event is nil")
	}

	// Try to extract common fields if it's a map
	eventMap, ok := orcaEvent.(map[string]interface{})
	if !ok {
		// If not a map, treat as opaque payload
		return &BridgeEvent{
			ID:        uuid.New().String(),
			Timestamp: time.Now(),
			Source:    SourceORCA,
			EventType: "unknown_orca_event",
			Payload:   orcaEvent,
		}, nil
	}

	// Extract fields from event map
	eventType := normalizeEventType(extractString(eventMap, "event_type", "type", "kind"))
	workTreeID := extractString(eventMap, "worktree_id", "worktree", "workspace_id")
	agentID := extractString(eventMap, "agent_id", "agent", "task_id")
	agentRole := extractString(eventMap, "agent_role", "role")
	previousState := extractString(eventMap, "previous_state", "old_state", "from")
	currentState := extractString(eventMap, "current_state", "new_state", "to")

	// Handle specific event types
	bridgeEvent := &BridgeEvent{
		ID:            uuid.New().String(),
		Timestamp:     time.Now(),
		Source:        SourceORCA,
		EventType:     eventType,
		WorktreeID:    workTreeID,
		AgentID:       agentID,
		AgentRole:     agentRole,
		PreviousState: previousState,
		CurrentState:  currentState,
		Payload:       orcaEvent,
	}

	// Translate payload based on event type
	if err := t.translateOrcaPayload(bridgeEvent, eventMap); err != nil {
		return nil, err
	}

	// Validate the resulting event
	if err := t.ValidateEvent(bridgeEvent); err != nil {
		return nil, fmt.Errorf("invalid translated ORCA event: %w", err)
	}

	return bridgeEvent, nil
}

// translateOrcaPayload translates event-specific payload data
func (t *Translator) translateOrcaPayload(event *BridgeEvent, eventMap map[string]interface{}) error {
	eventType := normalizeEventType(event.EventType)

	switch eventType {
	case EventCoordinatorStarted:
		event.EventType = EventCoordinatorStarted
		payload := &AgentStartedPayload{
			WorkflowID: extractString(eventMap, "workflow_id", "coordinator_id"),
			TaskID:     extractString(eventMap, "task_id", "coordinator_task_id"),
			Metadata:   extractMetadata(eventMap),
		}
		event.Payload = payload

	case EventWorkerSpawned:
		event.EventType = EventWorkerSpawned
		payload := &AgentStartedPayload{
			WorkflowID: extractString(eventMap, "workflow_id", "parent_workflow_id"),
			TaskID:     extractString(eventMap, "task_id", "worker_task_id"),
			Metadata:   extractMetadata(eventMap),
		}
		event.Payload = payload

	case EventAgentDecisionMade:
		event.EventType = EventAgentDecisionMade
		payload := &AgentDecisionPayload{
			DecisionID: extractString(eventMap, "decision_id", "event_id", "id"),
			Prompt:     extractString(eventMap, "prompt", "input", "request"),
			Response:   extractString(eventMap, "response", "output", "result"),
			DurationMs: extractInt64(eventMap, "duration_ms", "duration"),
		}
		event.Payload = payload

	case EventWorkerCompleted:
		event.EventType = EventWorkerCompleted
		resultBytes := extractBytes(eventMap, "result", "output", "data")
		payload := &AgentCompletedPayload{
			WorkflowID: extractString(eventMap, "workflow_id", "parent_workflow_id"),
			TaskID:     extractString(eventMap, "task_id", "worker_task_id"),
			Result:     resultBytes,
			Success:    extractBool(eventMap, "success", "completed", "ok"),
		}
		event.Payload = payload

	case EventWorkerFailed:
		event.EventType = EventWorkerFailed
		payload := &AgentFailedPayload{
			WorkflowID: extractString(eventMap, "workflow_id", "parent_workflow_id"),
			TaskID:     extractString(eventMap, "task_id", "worker_task_id"),
			Error:      extractString(eventMap, "error", "error_message", "reason"),
			Attempt:    int32(extractInt64(eventMap, "attempt", "retries", "retry_count")),
		}
		// Also set error fields on the main event
		event.ErrorMessage = payload.Error
		event.ErrorCode = extractString(eventMap, "error_code", "code", "status")
		event.Payload = payload

	case EventCoordinatorCompleted:
		event.EventType = EventCoordinatorCompleted
		resultBytes := extractBytes(eventMap, "result", "output", "data")
		payload := &AgentCompletedPayload{
			WorkflowID: extractString(eventMap, "workflow_id", "coordinator_id"),
			TaskID:     extractString(eventMap, "task_id", "coordinator_task_id"),
			Result:     resultBytes,
			Success:    extractBool(eventMap, "success", "completed", "ok"),
		}
		event.Payload = payload

	case EventAgentError:
		event.EventType = EventAgentError
		event.ErrorMessage = extractString(eventMap, "error", "error_message", "message")
		event.ErrorCode = extractString(eventMap, "error_code", "code", "status")
		// Payload remains the original event map

	default:
		// Unknown event type, keep payload as-is
		event.Payload = eventMap
	}

	return nil
}

// normalizeEventType maps common event type variants to canonical constants
func normalizeEventType(eventType string) string {
	switch eventType {
	// Coordinator started
	case EventCoordinatorStarted:
		return EventCoordinatorStarted

	// Worker spawned
	case EventWorkerSpawned:
		return EventWorkerSpawned

	// Agent decision made
	case EventAgentDecisionMade, "decision_made":
		return EventAgentDecisionMade

	// Worker completed
	case EventWorkerCompleted, "completed":
		return EventWorkerCompleted

	// Worker failed
	case EventWorkerFailed, "failed":
		return EventWorkerFailed

	// Coordinator completed
	case EventCoordinatorCompleted:
		return EventCoordinatorCompleted

	// Agent error
	case EventAgentError, "error":
		return EventAgentError

	default:
		return eventType
	}
}

// TranslateTemporalEvent converts a Temporal event to a BridgeEvent.
func (t *Translator) TranslateTemporalEvent(temporalEvent interface{}) (*BridgeEvent, error) {
	if temporalEvent == nil {
		return nil, fmt.Errorf("Temporal event is nil")
	}

	// Temporal events will be more structured (Protocol Buffer)
	// For now, treat as opaque and parse later
	return &BridgeEvent{
		ID:        uuid.New().String(),
		Timestamp: time.Now(),
		Source:    SourceTemporal,
		EventType: "unknown_temporal_event",
		Payload:   temporalEvent,
	}, nil
}

// ValidateEvent checks if a BridgeEvent is valid.
func (t *Translator) ValidateEvent(event *BridgeEvent) error {
	if event == nil {
		return fmt.Errorf("event is nil")
	}

	if event.ID == "" {
		return fmt.Errorf("event ID is required")
	}

	if event.Source != SourceORCA && event.Source != SourceTemporal {
		return fmt.Errorf("invalid event source: %s (must be '%s' or '%s')", event.Source, SourceORCA, SourceTemporal)
	}

	if event.EventType == "" {
		return fmt.Errorf("event type is required")
	}

	if event.Timestamp.IsZero() {
		return fmt.Errorf("event timestamp is required")
	}

	return nil
}

// Helper functions for extracting values from maps

func extractString(m map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if v, ok := m[key]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return ""
}

func extractInt64(m map[string]interface{}, keys ...string) int64 {
	for _, key := range keys {
		if v, ok := m[key]; ok {
			switch val := v.(type) {
			case int64:
				return val
			case int:
				return int64(val)
			case int32:
				return int64(val)
			case float64:
				return int64(val)
			}
		}
	}
	return 0
}

func extractBool(m map[string]interface{}, keys ...string) bool {
	for _, key := range keys {
		if v, ok := m[key]; ok {
			if b, ok := v.(bool); ok {
				return b
			}
		}
	}
	return false
}

func extractBytes(m map[string]interface{}, keys ...string) []byte {
	for _, key := range keys {
		if v, ok := m[key]; ok {
			switch val := v.(type) {
			case []byte:
				return val
			case string:
				return []byte(val)
			}
		}
	}
	return nil
}

func extractMetadata(m map[string]interface{}) map[string]string {
	metadata := make(map[string]string)
	if md, ok := m["metadata"]; ok {
		// Handle both map[string]interface{} and map[string]string
		if mdMap, ok := md.(map[string]interface{}); ok {
			for k, v := range mdMap {
				if s, ok := v.(string); ok {
					metadata[k] = s
				} else {
					metadata[k] = fmt.Sprintf("%v", v)
				}
			}
		} else if mdMap, ok := md.(map[string]string); ok {
			for k, v := range mdMap {
				metadata[k] = v
			}
		}
	}
	return metadata
}
