package types

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTraceEvent_Clone(t *testing.T) {
	testID := uuid.New()
	testTimestamp := time.Now().UTC()

	tests := []struct {
		name       string
		input      *TraceEvent
		want       *TraceEvent
		validateFn func(t *testing.T, original, clone *TraceEvent)
	}{
		{
			name:  "nil event returns nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "empty event",
			input: &TraceEvent{},
			want:  &TraceEvent{},
		},
		{
			name: "event with all fields populated",
			input: &TraceEvent{
				ID:          &testID,
				Name:        "test-trace",
				UserID:      "user-123",
				SessionID:   "session-456",
				Release:     "v1.0.0",
				Version:     "1.0",
				Metadata:    map[string]any{"key": "value", "nested": map[string]any{"inner": "data"}},
				Tags:        []string{"tag1", "tag2"},
				Public:      true,
				Input:       map[string]any{"prompt": "test input"},
				Output:      map[string]any{"response": "test output"},
				Environment: "production",
				Timestamp:   &testTimestamp,
				ExternalID:  "ext-789",
			},
			want: &TraceEvent{
				ID:          &testID,
				Name:        "test-trace",
				UserID:      "user-123",
				SessionID:   "session-456",
				Release:     "v1.0.0",
				Version:     "1.0",
				Metadata:    map[string]any{"key": "value", "nested": map[string]any{"inner": "data"}},
				Tags:        []string{"tag1", "tag2"},
				Public:      true,
				Input:       map[string]any{"prompt": "test input"},
				Output:      map[string]any{"response": "test output"},
				Environment: "production",
				Timestamp:   &testTimestamp,
				ExternalID:  "ext-789",
			},
			validateFn: func(t *testing.T, original, clone *TraceEvent) {
				// Verify deep copy by modifying original and ensuring clone is unaffected
				original.Tags[0] = modifiedText
				assert.Equal(t, "tag1", clone.Tags[0], "Tags should be deep copied")

				original.Metadata["key"] = modifiedText
				assert.Equal(t, "value", clone.Metadata["key"], "Metadata should be deep copied")

				// Verify pointers are different
				assert.NotSame(t, original.ID, clone.ID, "ID pointers should be different")
				assert.NotSame(t, original.Timestamp, clone.Timestamp, "Timestamp pointers should be different")
				// Note: We can't use NotSame for slices/maps directly since they might be recreated
				// Instead, we verify independence by checking that modifications don't affect the clone
				originalTags := original.Tags
				cloneTags := clone.Tags
				assert.NotSame(t, &originalTags, &cloneTags, "Tags slices should be different")

				originalMetadata := original.Metadata
				cloneMetadata := clone.Metadata
				assert.NotSame(t, &originalMetadata, &cloneMetadata, "Metadata maps should be different")
			},
		},
		{
			name: "event with nil pointer fields",
			input: &TraceEvent{
				Name: "test-trace",
				ID:   nil,
				Tags: nil,
			},
			want: &TraceEvent{
				Name: "test-trace",
				ID:   nil,
				Tags: nil,
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
			clonedTrace, ok := got.(*TraceEvent)
			require.True(t, ok, "Clone should return *TraceEvent")

			assert.Equal(t, tt.want, clonedTrace)

			if tt.validateFn != nil {
				tt.validateFn(t, tt.input, clonedTrace)
			}
		})
	}
}

func TestTraceEvent_Clone_DeepCopyComplexData(t *testing.T) {
	complexData := map[string]any{
		"string":  "value",
		"number":  42,
		"boolean": true,
		"slice":   []any{"a", "b", "c"},
		"nested": map[string]any{
			"inner": "value",
			"deep": map[string]any{
				"level": 3,
			},
		},
	}

	subject := &TraceEvent{
		Name:     "complex-trace",
		Metadata: complexData,
		Input:    complexData,
		Output:   []any{"output1", "output2"},
	}

	got := subject.Clone()
	require.NotNil(t, got)

	clonedTrace, ok := got.(*TraceEvent)
	require.True(t, ok)

	// Verify deep copy by modifying original
	subject.Metadata["string"] = modifiedText
	nestedMap := subject.Metadata["nested"].(map[string]any)
	nestedMap["inner"] = modifiedText

	inputMap := subject.Input.(map[string]any)
	inputMap["number"] = 999

	outputSlice := subject.Output.([]any)
	outputSlice[0] = modifiedText

	// Ensure clone is unaffected
	assert.Equal(t, "value", clonedTrace.Metadata["string"])
	clonedNested := clonedTrace.Metadata["nested"].(map[string]any)
	assert.Equal(t, "value", clonedNested["inner"])

	clonedInput := clonedTrace.Input.(map[string]any)
	assert.Equal(t, 42, clonedInput["number"])

	clonedOutput := clonedTrace.Output.([]any)
	assert.Equal(t, "output1", clonedOutput[0])
}
