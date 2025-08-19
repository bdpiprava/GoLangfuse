// Package integration_test provides utilities for creating test data structures
package integration_test

import (
	"time"

	"github.com/google/uuid"

	"github.com/bdpiprava/GoLangfuse/types"
)

// NewTestTraceEvent returns a sample TraceEvent for testing purposes
func NewTestTraceEvent() *types.TraceEvent {
	sessionID := uuid.New().String()
	userID := uuid.New().String()
	traceID := uuid.New()
	return &types.TraceEvent{
		ID:        &traceID,
		Name:      "SendTrace",
		SessionID: sessionID,
		UserID:    userID,
		Input:     "Test input for SendTrace",
		Output:    "Test output for SendTrace",
		Metadata: map[string]any{
			"key": "value",
		},
		Tags:   []string{"test", "integration"},
		Public: true,
	}
}

// NewTestGenerationEvent returns a sample GenerationEvent for testing purposes
func NewTestGenerationEvent() *types.GenerationEvent {
	traceID := uuid.New()
	generationID := uuid.New()
	now := time.Now().UTC()

	return &types.GenerationEvent{
		ID:        &generationID,
		Name:      "TestGeneration",
		TraceID:   &traceID,
		StartTime: &now,
		Model:     "gpt-3.5-turbo",
		Input:     "Test input for generation",
		Output:    "Test output for generation",
		Metadata: map[string]any{
			"test":    "generation",
			"version": float64(1),
		},
		Level: types.Default,
		Usage: types.Usage{
			Input:  10,
			Output: 5,
			Total:  15,
			Unit:   "TOKENS",
		},
		ModelParameters: map[string]any{
			"temperature": 0.7,
			"max_tokens":  float64(100),
		},
		PromptName:    "",
		PromptVersion: 0,
	}
}

// NewTestSpanEvent returns a sample SpanEvent for testing purposes
func NewTestSpanEvent() *types.SpanEvent {
	traceID := uuid.New()
	spanID := uuid.New()
	now := time.Now().UTC()

	return &types.SpanEvent{
		ID:        &spanID,
		Name:      "TestSpan",
		TraceID:   &traceID,
		StartTime: &now,
		Input:     "Test input for span",
		Output:    "Test output for span",
		Metadata: map[string]any{
			"test":      "span",
			"operation": "test_operation",
		},
		Level:   types.Default,
		Version: "1.0",
	}
}

// NewTestScoreEvent returns a sample ScoreEvent for testing purposes
func NewTestScoreEvent() *types.ScoreEvent {
	scoreID := uuid.New()
	traceID := uuid.New().String()
	comment := "Test score comment"

	return &types.ScoreEvent{
		ID:      &scoreID,
		Name:    "TestScore",
		TraceID: &traceID,
		Value:   0.85,
		Comment: &comment,
	}
}

// BuildSession constructs a Session from a slice of TraceEvents.
func BuildSession(traceEvents ...*types.TraceEvent) *types.Session {
	if len(traceEvents) == 0 {
		return nil
	}

	sessionID := traceEvents[0].SessionID
	traces := make([]types.Trace, len(traceEvents))
	for i, traceEvent := range traceEvents {
		traces[i] = types.Trace{
			ID:        traceEvent.ID.String(),
			Name:      traceEvent.Name,
			UserID:    traceEvent.UserID,
			Metadata:  traceEvent.Metadata,
			ProjectID: "test-project",
			Public:    traceEvent.Public,
			Tags:      traceEvent.Tags,
			Input:     traceEvent.Input,
			Output:    traceEvent.Output,
			SessionID: sessionID,
		}
	}

	return &types.Session{
		ID:        sessionID,
		CreatedAt: time.Now(),
		ProjectID: "test-project",
		Traces:    traces,
	}
}
