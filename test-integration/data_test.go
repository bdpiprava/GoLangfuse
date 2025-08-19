package integration_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bdpiprava/GoLangfuse/types"
)

func TestNewTestTraceEvent(t *testing.T) {
	traceEvent := NewTestTraceEvent()

	require.NotNil(t, traceEvent)
	assert.NotNil(t, traceEvent.ID)
	assert.Equal(t, "SendTrace", traceEvent.Name)
	assert.NotEmpty(t, traceEvent.SessionID)
	assert.NotEmpty(t, traceEvent.UserID)
	assert.Equal(t, "Test input for SendTrace", traceEvent.Input)
	assert.Equal(t, "Test output for SendTrace", traceEvent.Output)
	assert.NotNil(t, traceEvent.Metadata)
	assert.Contains(t, traceEvent.Tags, "test")
	assert.Contains(t, traceEvent.Tags, "integration")
	assert.True(t, traceEvent.Public)
}

func TestNewTestGenerationEvent(t *testing.T) {
	generationEvent := NewTestGenerationEvent()

	require.NotNil(t, generationEvent)
	assert.NotNil(t, generationEvent.ID)
	assert.Equal(t, "TestGeneration", generationEvent.Name)
	assert.NotNil(t, generationEvent.TraceID)
	assert.NotNil(t, generationEvent.StartTime)
	assert.Equal(t, "gpt-3.5-turbo", generationEvent.Model)
	assert.Equal(t, "Test input for generation", generationEvent.Input)
	assert.Equal(t, "Test output for generation", generationEvent.Output)
	assert.NotNil(t, generationEvent.Metadata)
	assert.Equal(t, types.Default, generationEvent.Level)
	assert.Equal(t, 10, generationEvent.Usage.Input)
	assert.Equal(t, 5, generationEvent.Usage.Output)
	assert.Equal(t, 15, generationEvent.Usage.Total)
	assert.NotNil(t, generationEvent.ModelParameters)
	assert.Equal(t, 0, generationEvent.PromptVersion)
}

func TestNewTestSpanEvent(t *testing.T) {
	spanEvent := NewTestSpanEvent()

	require.NotNil(t, spanEvent)
	assert.NotNil(t, spanEvent.ID)
	assert.Equal(t, "TestSpan", spanEvent.Name)
	assert.NotNil(t, spanEvent.TraceID)
	assert.NotNil(t, spanEvent.StartTime)
	assert.Equal(t, "Test input for span", spanEvent.Input)
	assert.Equal(t, "Test output for span", spanEvent.Output)
	assert.NotNil(t, spanEvent.Metadata)
	assert.Equal(t, types.Default, spanEvent.Level)
	assert.Equal(t, "1.0", spanEvent.Version)
}

func TestNewTestScoreEvent(t *testing.T) {
	scoreEvent := NewTestScoreEvent()

	require.NotNil(t, scoreEvent)
	assert.NotNil(t, scoreEvent.ID)
	assert.Equal(t, "TestScore", scoreEvent.Name)
	assert.NotNil(t, scoreEvent.TraceID)
	assert.Equal(t, float32(0.85), scoreEvent.Value)
	assert.NotNil(t, scoreEvent.Comment)
	assert.Equal(t, "Test score comment", *scoreEvent.Comment)
}

func TestBuildSession(t *testing.T) {
	traceEvent := NewTestTraceEvent()
	session := BuildSession(traceEvent)

	require.NotNil(t, session)
	assert.Equal(t, traceEvent.SessionID, session.ID)
	assert.Equal(t, "test-project", session.ProjectID)
	assert.Len(t, session.Traces, 1)

	trace := session.Traces[0]
	assert.Equal(t, traceEvent.ID.String(), trace.ID)
	assert.Equal(t, traceEvent.Name, trace.Name)
	assert.Equal(t, traceEvent.UserID, trace.UserID)
	assert.Equal(t, traceEvent.SessionID, trace.SessionID)
	assert.Equal(t, "test-project", trace.ProjectID)
	assert.Equal(t, traceEvent.Public, trace.Public)
	assert.Equal(t, traceEvent.Tags, trace.Tags)
	assert.Equal(t, traceEvent.Input, trace.Input)
	assert.Equal(t, traceEvent.Output, trace.Output)
}

func TestBuildSessionEmpty(t *testing.T) {
	session := BuildSession()
	assert.Nil(t, session)
}

func TestEventTimestamps(t *testing.T) {
	before := time.Now()

	generationEvent := NewTestGenerationEvent()
	spanEvent := NewTestSpanEvent()

	after := time.Now()

	require.NotNil(t, generationEvent.StartTime)
	require.NotNil(t, spanEvent.StartTime)

	assert.True(t, generationEvent.StartTime.After(before) || generationEvent.StartTime.Equal(before))
	assert.True(t, generationEvent.StartTime.Before(after) || generationEvent.StartTime.Equal(after))

	assert.True(t, spanEvent.StartTime.After(before) || spanEvent.StartTime.Equal(before))
	assert.True(t, spanEvent.StartTime.Before(after) || spanEvent.StartTime.Equal(after))
}
