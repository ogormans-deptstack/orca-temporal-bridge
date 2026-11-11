# Real-World ORCA-Temporal Bridge Flow

**How events flow from ORCA agents to Temporal's durable history in production**

---

## Scenario: Multi-Agent Code Review

Three agents work together to review code:
- **Coordinator**: Manages the review workflow
- **Worker 1**: Analyzes code quality (golangci-lint)
- **Worker 2**: Scans for security issues (gosec)

All events are captured in Temporal's immutable history for audit, replay, and compliance.

---

## Timeline: T0 → T7 (Events in ORCA)

### T0: Coordinator Starts
```json
{
  "event_type": "coordinator_started",
  "agent_id": "coordinator-v1",
  "worktree_id": "main-review-2024",
  "workflow_id": "code-review-001"
}
```
→ Sent to EventListener

### T1: Worker 1 Spawned
```json
{
  "event_type": "worker_spawned",
  "agent_id": "worker-quality",
  "task_id": "analyze-quality"
}
```
→ Sent to EventListener

### T2: Worker 1 Decides
```json
{
  "event_type": "agent_decision_made",
  "agent_id": "worker-quality",
  "prompt": "Check code quality",
  "response": "3 issues found: unused variable line 42, missing error handling line 105",
  "duration_ms": 245
}
```
→ Sent to EventListener

### T3: Worker 2 Spawned
```json
{
  "event_type": "worker_spawned",
  "agent_id": "worker-security",
  "task_id": "scan-security"
}
```
→ Sent to EventListener

### T4: Worker 2 Decides
```json
{
  "event_type": "agent_decision_made",
  "agent_id": "worker-security",
  "prompt": "Scan for vulnerabilities",
  "response": "No critical issues. 1 low-risk SQL injection potential line 332",
  "duration_ms": 189
}
```
→ Sent to EventListener

### T5-T7: Workers Complete, Coordinator Aggregates
Events flow through the bridge until coordinator_completed is sent.

---

## How The Bridge Works

### Component 1: EventListener (Real-Time Capture)

**What it does:**
- Subscribes to ORCA's event stream (via EventProvider interface)
- Receives events in real-time as they occur
- Forwards them to the event channel (buffered queue)
- **Non-blocking**: Doesn't hold up ORCA's execution

**Code:**
```go
for event := range provider.Subscribe(ctx) {
    // Non-blocking: ORCA doesn't wait
    eventChan <- event  // Buffered channel
    eventCount++
}
```

**Result**: ORCA finishes instantly, events queued for async processing.

---

### Component 2: Event Channel (Buffer)

**What it does:**
- Acts as a queue between listener and flusher
- Decouples real-time ORCA events from Temporal writes
- **Default**: 1000 events or 100ms (whichever comes first)
- Non-blocking for ORCA, batching for Temporal

**In our scenario:**
- 8 events arrive in ~500ms
- Buffer accumulates: [E1, E2, E3, E4, E5, E6, E7, E8]
- At 100ms mark: Flush triggered (even with only 8 events)

---

### Component 3: EventFlusher (Translation & Batch Write)

**What it does:**
1. Batches events from the channel
2. Translates each event (ORCA → Temporal format)
3. Starts appropriate Temporal workflows
4. Logs to Temporal's immutable history

**Workflow Starting:**

| Time | Event Type | Workflow Type | Workflow ID |
|------|-----------|---------------|-------------|
| T0 | coordinator_started | CoordinatorWorkflow | coordinator_started:coordinator-v1:1731280400000:abc1 |
| T1 | worker_spawned | WorkerWorkflow | worker_spawned:worker-quality:... |
| T2 | agent_decision_made | EventWorkflow | agent_decision_made:worker-quality:... |
| T3 | worker_spawned | WorkerWorkflow | worker_spawned:worker-security:... |
| T4 | agent_decision_made | EventWorkflow | agent_decision_made:worker-security:... |
| T5 | worker_completed | WorkerWorkflow | worker_completed:worker-quality:... |
| T6 | worker_completed | WorkerWorkflow | worker_completed:worker-security:... |
| T7 | coordinator_completed | CoordinatorWorkflow | coordinator_completed:coordinator-v1:... |

---

## Phase 3: Temporal Storage (Durability)

### Workflow Execution Example

```
Workflow ID: agent_decision_made:worker-quality:1731280401000:xyz9
Workflow Type: EventWorkflow
Status: Completed
Started: 2024-11-10T12:00:01Z
Completed: 2024-11-10T12:00:01.245Z

Input:
{
  "id": "uuid-1234",
  "timestamp": "2024-11-10T12:00:01Z",
  "source": "orca",
  "event_type": "agent_decision_made",
  "agent_id": "worker-quality",
  "payload": {
    "decision_id": "decision-quality-001",
    "prompt": "Check code quality",
    "response": "3 issues found...",
    "duration_ms": 245
  }
}
```

### Searchable Fields (Indexed)

```
WorkflowID: agent_decision_made:worker-quality:...
WorkflowType: EventWorkflow
Status: Completed
StartTime: 2024-11-10T12:00:01Z
SearchAttribute[agent_id]: worker-quality
SearchAttribute[event_type]: agent_decision_made
```

---

## Phase 4: Audit & Replay

### Scenario 1: Team Review

**Question**: "What did the security worker decide?"

**Temporal Query**:
```sql
WorkflowType = "EventWorkflow"
AND SearchAttribute["agent_id"] = "worker-security"
```

**Result**:
```
Workflow: agent_decision_made:worker-security:1731280402000:xyz10
Decision: "No critical issues. 1 low-risk potential"
Duration: 189ms
Timestamp: 2024-11-10T12:00:02Z
```

**Audit Trail**:
- ✅ What: Security scan results
- ✅ Who: worker-security
- ✅ When: exact timestamp
- ✅ How Long: 189ms
- ✅ Input: complete prompt

---

### Scenario 2: Replay from Step N

**Problem**: Something went wrong after T3. Replay to debug.

**In Temporal UI**:
1. Find the workflow at T3 (worker_spawned for worker-security)
2. Click "Replay from here"
3. Temporal re-executes from that point with same inputs
4. Compare results to identify what changed

**Benefits**:
- Debug without re-running ORCA
- Test code changes safely
- Prove auditability

---

## Key Architecture Principles

### 1. Non-Blocking

```
ORCA:     [E1] [E2] [E3] [E4] [E5] [E6] [E7] DONE
Bridge:    |    |    |    |    |    |    |   (async)
Temporal:                   [Batch Write]
```

ORCA finishes before Temporal writes complete → Real-time UX.

### 2. Batching

```
Without Bridge:
8 events → 8 gRPC calls → High latency

With Bridge:
8 events → Buffer 100ms → 1 gRPC call → Lower latency
```

### 3. Typed Payloads

```go
// Good: Structured data
payload := &AgentDecisionPayload{
    Prompt:     "Check code",
    Response:   "3 issues",
    DurationMs: 245,
}
```

Enables validation, indexing, and safety.

---

## Real-World Use Cases

### Use Case 1: Post-Mortem
"What exactly happened in this workflow?"
→ Query Temporal, get complete timeline with all decisions.

### Use Case 2: Compliance
"Prove what happened and when."
→ Export immutable audit trail with timestamps.

### Use Case 3: Performance Analysis
"Which workers are bottlenecks?"
→ Aggregate duration data, identify slow components.

### Use Case 4: Continuous Improvement
"Test updated logic."
→ Replay from step N, validate changes safely.

---

## Data Flow Summary

```
T0-T7: ORCA generates 8 events
  ↓
EventListener (capture in real-time)
  ↓
Event Channel (buffer 100ms or 1K)
  ↓
EventFlusher (translate + batch)
  ↓
Execute Workflows (gRPC to Temporal)
  ↓
Temporal Server (immutable history)
  ↓
Query + Replay (post-execution analysis)
```

---

## Testing Locally

```bash
# Start Temporal
make run-temporal

# Run integration tests
make test-with-temporal

# View Temporal UI
open http://localhost:8233

# Query workflows
temporal workflow list --address localhost:7233
temporal workflow describe --workflow-id "coordinator_started:..."
```

---

## Summary

The ORCA-Temporal Bridge transforms a real-time multi-agent system into an auditable, replayable, compliant workflow engine:

- **ORCA** handles local development with instant feedback
- **Bridge** captures events without blocking
- **Temporal** stores durable history for audit, replay, and analysis

This is how enterprises build production-grade multi-agent systems.
