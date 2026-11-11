package bridge

import (
	"context"
	"testing"
	"time"

	"github.com/ogormans-deptstack/orca-temporal-bridge/internal/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBridge_Lifecycle tests basic bridge start/stop
func TestBridge_Lifecycle(t *testing.T) {
	ctx := context.Background()
	mockProvider := NewMockEventProvider()
	defer func() { _ = mockProvider.Close() }()

	config := Config{
		TemporalHost:  "localhost",
		TemporalPort:  7233,
		Namespace:     "default",
		TaskQueue:     "test-queue",
		EventProvider: mockProvider,
		BufferSize:    100,
		FlushInterval: 100 * time.Millisecond,
	}

	bridge, err := NewBridge(ctx, config)
	require.NoError(t, err)
	assert.NotNil(t, bridge)

	// Initially not running
	assert.False(t, bridge.IsRunning())

	// Note: We can't actually start without a real Temporal server
	// In CI/CD, docker-compose would start it first
	// For now, just verify the bridge can be created
}

// TestBridge_Configuration tests config validation
func TestBridge_Configuration(t *testing.T) {
	ctx := context.Background()

	// Test missing provider
	_, err := NewBridge(ctx, Config{})
	assert.Error(t, err)
	assert.Equal(t, "EventProvider is required", err.Error())

	// Test valid config
	mockProvider := NewMockEventProvider()
	defer func() { _ = mockProvider.Close() }()

	config := Config{
		TemporalHost:  "localhost",
		TemporalPort:  7233,
		Namespace:     "default",
		TaskQueue:     "test-queue",
		EventProvider: mockProvider,
		BufferSize:    0, // Will be set to default
		FlushInterval: 0, // Will be set to default
	}

	bridge, err := NewBridge(ctx, config)
	require.NoError(t, err)
	assert.NotNil(t, bridge)

	// Verify defaults were set by NewBridge
	assert.Equal(t, 1000, bridge.config.BufferSize)
	assert.Equal(t, 100*time.Millisecond, bridge.config.FlushInterval)
}

// TestBridge_GetStatus tests status reporting
func TestBridge_GetStatus(t *testing.T) {
	ctx := context.Background()
	mockProvider := NewMockEventProvider()
	defer func() { _ = mockProvider.Close() }()

	config := Config{
		TemporalHost:  "localhost",
		TemporalPort:  7233,
		Namespace:     "default",
		TaskQueue:     "test-queue",
		EventProvider: mockProvider,
		BufferSize:    100,
	}

	bridge, err := NewBridge(ctx, config)
	require.NoError(t, err)

	status := bridge.GetStatus()
	assert.False(t, status["running"].(bool))
	assert.Equal(t, "localhost:7233", status["temporal_url"].(string))
	assert.Equal(t, "default", status["namespace"].(string))
	assert.Equal(t, "test-queue", status["task_queue"].(string))
}

// TestEventListener_Creation tests listener initialization
func TestEventListener_Creation(t *testing.T) {
	mockProvider := NewMockEventProvider()
	defer func() { _ = mockProvider.Close() }()

	eventChan := make(chan *events.BridgeEvent, 100)
	listener := NewEventListener(mockProvider, eventChan)

	assert.NotNil(t, listener)
	stats := listener.GetStats()
	assert.Equal(t, int64(0), stats["events_received"].(int64))
	assert.Equal(t, int64(0), stats["errors"].(int64))
}

// TestEventListener_Subscription tests event subscription
func TestEventListener_Subscription(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	mockProvider := NewMockEventProvider()
	defer func() { _ = mockProvider.Close() }()

	eventChan := make(chan *events.BridgeEvent, 100)
	listener := NewEventListener(mockProvider, eventChan)

	// Start listener in goroutine
	go func() {
		_ = listener.Listen(ctx)
	}()

	// Emit test event
	err := mockProvider.EmitCoordinatorStarted("wt-1", "coord-1", map[string]string{"version": "1.0"})
	require.NoError(t, err)

	// Receive event on channel
	select {
	case event := <-eventChan:
		require.NotNil(t, event)
		assert.Equal(t, events.EventCoordinatorStarted, event.EventType)
		assert.Equal(t, "coord-1", event.AgentID)
		assert.Equal(t, "wt-1", event.WorktreeID)
	case <-ctx.Done():
		t.Fatal("timeout waiting for event")
	}

	// Verify stats updated
	stats := listener.GetStats()
	assert.Equal(t, int64(1), stats["events_received"].(int64))
}

// TestEventFlusher_Creation tests flusher initialization
func TestEventFlusher_Creation(t *testing.T) {
	// Note: Can't test full flusher without real Temporal client
	// Just test object creation and stats

	config := Config{
		BufferSize:    100,
		FlushInterval: 100 * time.Millisecond,
	}

	eventChan := make(chan *events.BridgeEvent, 100)
	flusher := NewEventFlusher(nil, eventChan, config)

	assert.NotNil(t, flusher)
	stats := flusher.GetStats()
	assert.Equal(t, int64(0), stats["events_processed"].(int64))
	assert.Equal(t, int64(0), stats["flushes"].(int64))
	assert.Equal(t, int64(0), stats["errors"].(int64))
	assert.Equal(t, 0, stats["buffer_size"].(int))
}

// TestMockEventProvider_CoordinatorStarted tests mock provider event emission
func TestMockEventProvider_CoordinatorStarted(t *testing.T) {
	provider := NewMockEventProvider()
	defer func() { _ = provider.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	eventStream, err := provider.Subscribe(ctx)
	require.NoError(t, err)

	// Emit event
	err = provider.EmitCoordinatorStarted("wt-1", "coord-1", map[string]string{"key": "value"})
	require.NoError(t, err)

	// Receive event
	select {
	case event := <-eventStream:
		require.NotNil(t, event)
		assert.Equal(t, events.EventCoordinatorStarted, event.EventType)
		assert.Equal(t, "coord-1", event.AgentID)

		payload := event.Payload.(*events.AgentStartedPayload)
		assert.NotNil(t, payload)
		assert.Equal(t, "coord-1-workflow", payload.WorkflowID)
		assert.Equal(t, "value", payload.Metadata["key"])
	case <-ctx.Done():
		t.Fatal("timeout waiting for event")
	}
}

// TestMockEventProvider_WorkerSpawned tests worker spawned event
func TestMockEventProvider_WorkerSpawned(t *testing.T) {
	provider := NewMockEventProvider()
	defer func() { _ = provider.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	eventStream, err := provider.Subscribe(ctx)
	require.NoError(t, err)

	err = provider.EmitWorkerSpawned("wt-1", "worker-1", "parent-wf")
	require.NoError(t, err)

	select {
	case event := <-eventStream:
		require.NotNil(t, event)
		assert.Equal(t, events.EventWorkerSpawned, event.EventType)
		assert.Equal(t, "worker-1", event.AgentID)
		assert.Equal(t, "worker", event.AgentRole)

		payload := event.Payload.(*events.AgentStartedPayload)
		assert.Equal(t, "parent-wf", payload.WorkflowID)
	case <-ctx.Done():
		t.Fatal("timeout waiting for event")
	}
}

// TestMockEventProvider_AgentDecisionMade tests agent decision event
func TestMockEventProvider_AgentDecisionMade(t *testing.T) {
	provider := NewMockEventProvider()
	defer func() { _ = provider.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	eventStream, err := provider.Subscribe(ctx)
	require.NoError(t, err)

	err = provider.EmitAgentDecisionMade("agent-1", "What is 2+2?", "4", 150)
	require.NoError(t, err)

	select {
	case event := <-eventStream:
		require.NotNil(t, event)
		assert.Equal(t, events.EventAgentDecisionMade, event.EventType)

		payload := event.Payload.(*events.AgentDecisionPayload)
		assert.Equal(t, "What is 2+2?", payload.Prompt)
		assert.Equal(t, "4", payload.Response)
		assert.Equal(t, int64(150), payload.DurationMs)
	case <-ctx.Done():
		t.Fatal("timeout waiting for event")
	}
}

// TestMockEventProvider_WorkerCompleted tests worker completed event
func TestMockEventProvider_WorkerCompleted(t *testing.T) {
	provider := NewMockEventProvider()
	defer func() { _ = provider.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	eventStream, err := provider.Subscribe(ctx)
	require.NoError(t, err)

	err = provider.EmitWorkerCompleted("worker-1", true)
	require.NoError(t, err)

	select {
	case event := <-eventStream:
		require.NotNil(t, event)
		assert.Equal(t, events.EventWorkerCompleted, event.EventType)

		payload := event.Payload.(*events.AgentCompletedPayload)
		assert.True(t, payload.Success)
		assert.Equal(t, []byte("test result"), payload.Result)
	case <-ctx.Done():
		t.Fatal("timeout waiting for event")
	}
}

// TestMockEventProvider_WorkerFailed tests worker failed event
func TestMockEventProvider_WorkerFailed(t *testing.T) {
	provider := NewMockEventProvider()
	defer func() { _ = provider.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	eventStream, err := provider.Subscribe(ctx)
	require.NoError(t, err)

	err = provider.EmitWorkerFailed("worker-1", "timeout exceeded", 2)
	require.NoError(t, err)

	select {
	case event := <-eventStream:
		require.NotNil(t, event)
		assert.Equal(t, events.EventWorkerFailed, event.EventType)
		assert.Equal(t, "timeout exceeded", event.ErrorMessage)

		payload := event.Payload.(*events.AgentFailedPayload)
		assert.Equal(t, "timeout exceeded", payload.Error)
		assert.Equal(t, int32(2), payload.Attempt)
	case <-ctx.Done():
		t.Fatal("timeout waiting for event")
	}
}

// TestMockEventProvider_MultipleEvents tests emitting multiple events
func TestMockEventProvider_MultipleEvents(t *testing.T) {
	provider := NewMockEventProvider()
	defer func() { _ = provider.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	eventStream, err := provider.Subscribe(ctx)
	require.NoError(t, err)

	// Emit multiple events
	err = provider.EmitCoordinatorStarted("wt-1", "coord-1", nil)
	require.NoError(t, err)

	err = provider.EmitWorkerSpawned("wt-1", "worker-1", "coord-1-wf")
	require.NoError(t, err)

	err = provider.EmitAgentDecisionMade("worker-1", "test", "result", 100)
	require.NoError(t, err)

	// Receive all events
	for i := 0; i < 3; i++ {
		select {
		case event := <-eventStream:
			require.NotNil(t, event)
		case <-ctx.Done():
			t.Fatalf("timeout waiting for event %d", i+1)
		}
	}
}

// TestMockEventProvider_Close tests provider closing
func TestMockEventProvider_Close(t *testing.T) {
	provider := NewMockEventProvider()

	// Close should succeed
	err := provider.Close()
	assert.NoError(t, err)

	// Double close should succeed (idempotent)
	err = provider.Close()
	assert.NoError(t, err)

	// Emit after close should not error but drop event
	err = provider.EmitCoordinatorStarted("wt-1", "coord-1", nil)
	assert.NoError(t, err) // Silently drops
}
