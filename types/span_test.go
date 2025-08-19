package types

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpanEvent_Clone(t *testing.T) {
	testID := uuid.New()
	testTraceID := uuid.New()
	testParentID := uuid.New()
	testStartTime := time.Now().UTC()
	testEndTime := testStartTime.Add(time.Second)

	tests := []struct {
		name       string
		input      *SpanEvent
		want       *SpanEvent
		validateFn func(t *testing.T, original, clone *SpanEvent)
	}{
		{
			name:  "nil event returns nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "empty event",
			input: &SpanEvent{},
			want:  &SpanEvent{},
		},
		{
			name: "event with all fields populated",
			input: &SpanEvent{
				ID:                  &testID,
				TraceID:             &testTraceID,
				ParentObservationID: &testParentID,
				Name:                "test-span",
				StartTime:           &testStartTime,
				EndTime:             &testEndTime,
				Metadata:            map[string]any{"operation": "database", "table": "users"},
				Level:               Default,
				StatusMessage:       "completed",
				Input:               map[string]any{"query": "SELECT * FROM users"},
				Output:              map[string]any{"rows": 10, "duration": "100ms"},
				Version:             "1.0",
				Environment:         "production",
			},
			want: &SpanEvent{
				ID:                  &testID,
				TraceID:             &testTraceID,
				ParentObservationID: &testParentID,
				Name:                "test-span",
				StartTime:           &testStartTime,
				EndTime:             &testEndTime,
				Metadata:            map[string]any{"operation": "database", "table": "users"},
				Level:               Default,
				StatusMessage:       "completed",
				Input:               map[string]any{"query": "SELECT * FROM users"},
				Output:              map[string]any{"rows": 10, "duration": "100ms"},
				Version:             "1.0",
				Environment:         "production",
			},
			validateFn: func(t *testing.T, original, clone *SpanEvent) {
				// Verify deep copy by modifying original and ensuring clone is unaffected
				original.Metadata["operation"] = modifiedText
				assert.Equal(t, "database", clone.Metadata["operation"], "Metadata should be deep copied")

				// Verify pointers are different
				assert.NotSame(t, original.ID, clone.ID, "ID pointers should be different")
				assert.NotSame(t, original.TraceID, clone.TraceID, "TraceID pointers should be different")
				assert.NotSame(t, original.ParentObservationID, clone.ParentObservationID, "ParentObservationID pointers should be different")
				assert.NotSame(t, original.StartTime, clone.StartTime, "StartTime pointers should be different")
				assert.NotSame(t, original.EndTime, clone.EndTime, "EndTime pointers should be different")
				originalMetadata := original.Metadata
				cloneMetadata := clone.Metadata
				assert.NotSame(t, &originalMetadata, &cloneMetadata, "Metadata maps should be different")
			},
		},
		{
			name: "event with nil pointer fields",
			input: &SpanEvent{
				Name:     "test-span",
				ID:       nil,
				TraceID:  nil,
				Metadata: nil,
			},
			want: &SpanEvent{
				Name:     "test-span",
				ID:       nil,
				TraceID:  nil,
				Metadata: nil,
			},
		},
		{
			name: "event with error level",
			input: &SpanEvent{
				Name:          "error-span",
				Level:         Error,
				StatusMessage: "Database connection failed",
			},
			want: &SpanEvent{
				Name:          "error-span",
				Level:         Error,
				StatusMessage: "Database connection failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.Clone()

			if tt.want == nil {
				assert.Nil(t, got)
				return
			}

			require.NotNil(t, got)
			clonedSpan, ok := got.(*SpanEvent)
			require.True(t, ok, "Clone should return *SpanEvent")

			assert.Equal(t, tt.want, clonedSpan)

			if tt.validateFn != nil {
				tt.validateFn(t, tt.input, clonedSpan)
			}
		})
	}
}

func TestSpanEvent_Clone_NestedStructures(t *testing.T) {
	complexMetadata := map[string]any{
		"database": map[string]any{
			"connection": map[string]any{
				"host":     "localhost",
				"port":     5432,
				"database": "testdb",
			},
			"query": map[string]any{
				"sql":        "SELECT * FROM users WHERE id = ?",
				"parameters": []any{123, "active"},
			},
		},
		"performance": map[string]any{
			"metrics": []any{
				map[string]any{"name": "duration", "value": 150.5},
				map[string]any{"name": "rows_affected", "value": 1},
			},
		},
	}

	subject := &SpanEvent{
		Name:     "complex-span",
		Metadata: complexMetadata,
		Input: map[string]any{
			"operation": "database_query",
			"context": map[string]any{
				"user_id":    "123",
				"session_id": "abc-def",
			},
		},
		Output: []map[string]any{
			{"id": 1, "name": "User 1"},
			{"id": 2, "name": "User 2"},
		},
	}

	got := subject.Clone()
	require.NotNil(t, got)

	clonedSpan, ok := got.(*SpanEvent)
	require.True(t, ok)

	// Verify deep copy by modifying original
	dbConfig := subject.Metadata["database"].(map[string]any)
	connection := dbConfig["connection"].(map[string]any)
	connection["host"] = modifiedText

	queryConfig := dbConfig["query"].(map[string]any)
	params := queryConfig["parameters"].([]any)
	params[0] = 999

	inputMap := subject.Input.(map[string]any)
	context := inputMap["context"].(map[string]any)
	context["user_id"] = modifiedText

	outputSlice := subject.Output.([]map[string]any)
	outputSlice[0]["name"] = "Modified User"

	// Ensure clone is unaffected
	clonedDbConfig := clonedSpan.Metadata["database"].(map[string]any)
	clonedConnection := clonedDbConfig["connection"].(map[string]any)
	assert.Equal(t, "localhost", clonedConnection["host"])

	clonedQueryConfig := clonedDbConfig["query"].(map[string]any)
	clonedParams := clonedQueryConfig["parameters"].([]any)
	assert.Equal(t, 123, clonedParams[0])

	clonedInput := clonedSpan.Input.(map[string]any)
	clonedContext := clonedInput["context"].(map[string]any)
	assert.Equal(t, "123", clonedContext["user_id"])

	clonedOutput := clonedSpan.Output.([]map[string]any)
	assert.Equal(t, "User 1", clonedOutput[0]["name"])
}
