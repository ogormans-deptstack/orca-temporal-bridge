package bridge_test

import (
	"context"
	"testing"

	"github.com/ogormans-deptstack/orca-temporal-bridge/internal/bridge"
)

func TestBridgeCreation(t *testing.T) {
	config := bridge.Config{
		TemporalHost: "localhost",
		TemporalPort: "7233",
		Namespace:    "default",
	}

	ctx := context.Background()
	b, err := bridge.NewBridge(ctx, config)
	if err != nil {
		t.Fatalf("Failed to create bridge: %v", err)
	}

	if b == nil {
		t.Fatal("Bridge is nil")
	}

	t.Log("✓ Bridge created successfully")
}

// TODO: Add more comprehensive tests:
// - TestBridgeStart
// - TestBridgeStop
// - TestEventHandling
// - TestEventFlushing
// - TestErrorRecovery
// - TestConcurrency
