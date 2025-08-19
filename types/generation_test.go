package types

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerationEvent_Clone(t *testing.T) {
	testID := uuid.New()
	testTraceID := uuid.New()
	testParentID := uuid.New()
	testStartTime := time.Now().UTC()
	testCompletionStartTime := testStartTime.Add(time.Millisecond * 100)
	testEndTime := testStartTime.Add(time.Second)

	tests := []struct {
		name       string
		input      *GenerationEvent
		want       *GenerationEvent
		validateFn func(t *testing.T, original, clone *GenerationEvent)
	}{
		{
			name:  "nil event returns nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "empty event",
			input: &GenerationEvent{},
			want:  &GenerationEvent{},
		},
		{
			name: "event with all fields populated",
			input: &GenerationEvent{
				ID:                  &testID,
				Name:                "test-generation",
				TraceID:             &testTraceID,
				StartTime:           &testStartTime,
				CompletionStartTime: &testCompletionStartTime,
				EndTime:             &testEndTime,
				Metadata:            map[string]any{"key": "value", "config": map[string]any{"temp": 0.7}},
				Model:               "gpt-4",
				Input:               map[string]any{"prompt": "test prompt"},
				Output:              map[string]any{"completion": "test completion"},
				Level:               Default,
				StatusMessage:       "success",
				ParentObservationID: &testParentID,
				Version:             "1.0",
				ModelParameters:     map[string]any{"temperature": 0.7, "max_tokens": 1000},
				Usage: Usage{
					Input:            100,
					Output:           50,
					Total:            150,
					Unit:             Tokens,
					PromptTokens:     100,
					CompletionTokens: 50,
					TotalTokens:      150,
				},
				PromptVersion: 1,
				PromptName:    "test-prompt",
			},
			want: &GenerationEvent{
				ID:                  &testID,
				Name:                "test-generation",
				TraceID:             &testTraceID,
				StartTime:           &testStartTime,
				CompletionStartTime: &testCompletionStartTime,
				EndTime:             &testEndTime,
				Metadata:            map[string]any{"key": "value", "config": map[string]any{"temp": 0.7}},
				Model:               "gpt-4",
				Input:               map[string]any{"prompt": "test prompt"},
				Output:              map[string]any{"completion": "test completion"},
				Level:               Default,
				StatusMessage:       "success",
				ParentObservationID: &testParentID,
				Version:             "1.0",
				ModelParameters:     map[string]any{"temperature": 0.7, "max_tokens": 1000},
				Usage: Usage{
					Input:            100,
					Output:           50,
					Total:            150,
					Unit:             Tokens,
					PromptTokens:     100,
					CompletionTokens: 50,
					TotalTokens:      150,
				},
				PromptVersion: 1,
				PromptName:    "test-prompt",
			},
			validateFn: func(t *testing.T, original, clone *GenerationEvent) {
				// Verify deep copy by modifying original and ensuring clone is unaffected
				original.Metadata["key"] = modifiedText
				assert.Equal(t, "value", clone.Metadata["key"], "Metadata should be deep copied")

				original.ModelParameters["temperature"] = 1.0
				assert.InEpsilon(t, 0.7, clone.ModelParameters["temperature"], 0.0001, "ModelParameters should be deep copied")

				// Verify pointers are different
				assert.NotSame(t, original.ID, clone.ID, "ID pointers should be different")
				assert.NotSame(t, original.TraceID, clone.TraceID, "TraceID pointers should be different")
				assert.NotSame(t, original.StartTime, clone.StartTime, "StartTime pointers should be different")
				originalMetadata := original.Metadata
				cloneMetadata := clone.Metadata
				assert.NotSame(t, &originalMetadata, &cloneMetadata, "Metadata maps should be different")

				originalModelParams := original.ModelParameters
				cloneModelParams := clone.ModelParameters
				assert.NotSame(t, &originalModelParams, &cloneModelParams, "ModelParameters maps should be different")
			},
		},
		{
			name: "event with nil pointer fields",
			input: &GenerationEvent{
				Name:     "test-generation",
				ID:       nil,
				TraceID:  nil,
				Metadata: nil,
			},
			want: &GenerationEvent{
				Name:     "test-generation",
				ID:       nil,
				TraceID:  nil,
				Metadata: nil,
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
			clonedGeneration, ok := got.(*GenerationEvent)
			require.True(t, ok, "Clone should return *GenerationEvent")

			assert.Equal(t, tt.want, clonedGeneration)

			if tt.validateFn != nil {
				tt.validateFn(t, tt.input, clonedGeneration)
			}
		})
	}
}

func TestGenerationEvent_Clone_ComplexInputOutput(t *testing.T) {
	complexInput := []map[string]any{
		{
			"role":    "user",
			"content": "Hello",
			"metadata": map[string]any{
				"timestamp": "2023-01-01",
			},
		},
		{
			"role":    "assistant",
			"content": "Hi there!",
		},
	}

	complexOutput := map[string]any{
		"choices": []any{
			map[string]any{
				"message": map[string]any{
					"role":    "assistant",
					"content": "Response",
				},
			},
		},
		"usage": map[string]any{
			"prompt_tokens":     10,
			"completion_tokens": 5,
		},
	}

	subject := &GenerationEvent{
		Name:   "complex-generation",
		Input:  complexInput,
		Output: complexOutput,
		Metadata: map[string]any{
			"experiment": "test",
			"params":     []string{"a", "b", "c"},
		},
	}

	got := subject.Clone()
	require.NotNil(t, got)

	clonedGeneration, ok := got.(*GenerationEvent)
	require.True(t, ok)

	// Verify deep copy by modifying original
	inputSlice := subject.Input.([]map[string]any)
	inputSlice[0]["role"] = modifiedText

	outputMap := subject.Output.(map[string]any)
	choices := outputMap["choices"].([]any)
	firstChoice := choices[0].(map[string]any)
	firstChoice["message"].(map[string]any)["content"] = modifiedText

	subject.Metadata["experiment"] = modifiedText

	// Ensure clone is unaffected
	clonedInput := clonedGeneration.Input.([]map[string]any)
	assert.Equal(t, "user", clonedInput[0]["role"])

	clonedOutput := clonedGeneration.Output.(map[string]any)
	clonedChoices := clonedOutput["choices"].([]any)
	clonedFirstChoice := clonedChoices[0].(map[string]any)
	assert.Equal(t, "Response", clonedFirstChoice["message"].(map[string]any)["content"])

	assert.Equal(t, "test", clonedGeneration.Metadata["experiment"])
}
