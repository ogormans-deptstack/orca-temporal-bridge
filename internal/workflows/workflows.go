// Package workflows defines Temporal workflow definitions for agent orchestration.
package workflows

import (
	"go.temporal.io/sdk/workflow"
)

// CoordinatorWorkflow is the parent workflow that orchestrates worker agents.
// It:
// - Receives task definitions
// - Spawns child workflows for each worker agent
// - Waits for all workers to complete
// - Handles signals (pause, resume, cancel)
func CoordinatorWorkflow(ctx workflow.Context, taskDef interface{}) (interface{}, error) {
	// TODO: Implement coordinator workflow
	// - Start child workflows for workers
	// - Wait for completion
	// - Handle signals
	// - Return aggregated results

	return nil, nil
}

// WorkerWorkflow is a child workflow that executes a single agent task.
// It:
// - Receives task parameters
// - Records agent decision events
// - Completes with result or error
func WorkerWorkflow(ctx workflow.Context, taskDef interface{}) (interface{}, error) {
	// TODO: Implement worker workflow
	// - Wait for task start signal
	// - Execute task (via activity or direct logic)
	// - Record decision events
	// - Return result

	return nil, nil
}

// RegisterWorkflows registers all workflows with Temporal worker.
// This is a placeholder for Sprint 1.3 — actual worker registration
// will be done when the Temporal worker is initialized in main.
func RegisterWorkflows() {
	// TODO: Register workflows with worker
	// r.RegisterWorkflow(CoordinatorWorkflow)
	// r.RegisterWorkflow(WorkerWorkflow)
}
