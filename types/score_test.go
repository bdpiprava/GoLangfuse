package types

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScoreEvent_Clone(t *testing.T) {
	testID := uuid.New()
	testTraceID := "trace-123"
	testSessionID := "session-456"
	testObservationID := "observation-789"
	testComment := "Great response quality"
	testDatasetRunID := "run-abc"
	testEnvironment := "production"

	tests := []struct {
		name       string
		input      *ScoreEvent
		want       *ScoreEvent
		validateFn func(t *testing.T, original, clone *ScoreEvent)
	}{
		{
			name:  "nil event returns nil",
			input: nil,
			want:  nil,
		},
		{
			name:  "empty event",
			input: &ScoreEvent{},
			want:  &ScoreEvent{},
		},
		{
			name: "event with all fields populated",
			input: &ScoreEvent{
				ID:            &testID,
				Name:          "quality_score",
				TraceID:       &testTraceID,
				SessionID:     &testSessionID,
				ObservationID: &testObservationID,
				Value:         0.95,
				Comment:       &testComment,
				DatasetRunID:  &testDatasetRunID,
				Environment:   &testEnvironment,
				Metadata: map[string]any{
					"evaluator": "human",
					"criteria":  []string{"accuracy", "relevance", "completeness"},
					"details": map[string]any{
						"accuracy":     0.9,
						"relevance":    1.0,
						"completeness": 0.85,
					},
				},
			},
			want: &ScoreEvent{
				ID:            &testID,
				Name:          "quality_score",
				TraceID:       &testTraceID,
				SessionID:     &testSessionID,
				ObservationID: &testObservationID,
				Value:         0.95,
				Comment:       &testComment,
				DatasetRunID:  &testDatasetRunID,
				Environment:   &testEnvironment,
				Metadata: map[string]any{
					"evaluator": "human",
					"criteria":  []string{"accuracy", "relevance", "completeness"},
					"details": map[string]any{
						"accuracy":     0.9,
						"relevance":    1.0,
						"completeness": 0.85,
					},
				},
			},
			validateFn: func(t *testing.T, original, clone *ScoreEvent) {
				// Verify deep copy by modifying original and ensuring clone is unaffected
				original.Metadata["evaluator"] = modifiedText
				assert.Equal(t, "human", clone.Metadata["evaluator"], "Metadata should be deep copied")

				criteria := original.Metadata["criteria"].([]string)
				criteria[0] = modifiedText
				clonedCriteria := clone.Metadata["criteria"].([]string)
				assert.Equal(t, "accuracy", clonedCriteria[0], "Nested slice should be deep copied")

				// Verify pointers are different
				assert.NotSame(t, original.ID, clone.ID, "ID pointers should be different")
				assert.NotSame(t, original.TraceID, clone.TraceID, "TraceID pointers should be different")
				assert.NotSame(t, original.SessionID, clone.SessionID, "SessionID pointers should be different")
				assert.NotSame(t, original.ObservationID, clone.ObservationID, "ObservationID pointers should be different")
				assert.NotSame(t, original.Comment, clone.Comment, "Comment pointers should be different")
				assert.NotSame(t, original.DatasetRunID, clone.DatasetRunID, "DatasetRunID pointers should be different")
				assert.NotSame(t, original.Environment, clone.Environment, "Environment pointers should be different")
				originalMetadata := original.Metadata
				cloneMetadata := clone.Metadata
				assert.NotSame(t, &originalMetadata, &cloneMetadata, "Metadata maps should be different")
			},
		},
		{
			name: "event with nil pointer fields",
			input: &ScoreEvent{
				Name:  "test-score",
				Value: 0.8,
				ID:    nil,
			},
			want: &ScoreEvent{
				Name:  "test-score",
				Value: 0.8,
				ID:    nil,
			},
		},
		{
			name: "event with minimal required fields",
			input: &ScoreEvent{
				Name:  "simple_score",
				Value: 1.0,
			},
			want: &ScoreEvent{
				Name:  "simple_score",
				Value: 1.0,
			},
		},
		{
			name: "event with zero value",
			input: &ScoreEvent{
				Name:  "zero_score",
				Value: 0.0,
			},
			want: &ScoreEvent{
				Name:  "zero_score",
				Value: 0.0,
			},
		},
		{
			name: "event with negative value",
			input: &ScoreEvent{
				Name:  "penalty_score",
				Value: -0.5,
			},
			want: &ScoreEvent{
				Name:  "penalty_score",
				Value: -0.5,
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
			clonedScore, ok := got.(*ScoreEvent)
			require.True(t, ok, "Clone should return *ScoreEvent")

			assert.Equal(t, tt.want, clonedScore)

			if tt.validateFn != nil {
				tt.validateFn(t, tt.input, clonedScore)
			}
		})
	}
}

func TestScoreEvent_Clone_ComplexMetadata(t *testing.T) {
	complexMetadata := map[string]any{
		"evaluation": map[string]any{
			"model":     "gpt-4",
			"version":   "2024-01",
			"timestamp": "2024-01-15T10:30:00Z",
			"metrics": map[string]any{
				"scores": []any{
					map[string]any{"dimension": "accuracy", "value": 0.95},
					map[string]any{"dimension": "relevance", "value": 0.88},
					map[string]any{"dimension": "completeness", "value": 0.92},
				},
				"aggregated": map[string]any{
					"mean":   0.916,
					"median": 0.92,
					"std":    0.029,
				},
			},
		},
		"context": map[string]any{
			"task_type":         "question_answering",
			"domain":            "customer_support",
			"difficulty":        "medium",
			"reference_answers": []string{"answer1", "answer2"},
		},
	}

	subject := &ScoreEvent{
		Name:     "comprehensive_evaluation",
		Value:    0.916,
		Metadata: complexMetadata,
	}

	got := subject.Clone()
	require.NotNil(t, got)

	clonedScore, ok := got.(*ScoreEvent)
	require.True(t, ok)

	// Verify deep copy by modifying original
	evaluation := subject.Metadata["evaluation"].(map[string]any)
	evaluation["model"] = modifiedText

	metrics := evaluation["metrics"].(map[string]any)
	scores := metrics["scores"].([]any)
	firstScore := scores[0].(map[string]any)
	firstScore["value"] = 0.5

	context := subject.Metadata["context"].(map[string]any)
	refAnswers := context["reference_answers"].([]string)
	refAnswers[0] = modifiedText

	// Ensure clone is unaffected
	clonedEvaluation := clonedScore.Metadata["evaluation"].(map[string]any)
	assert.Equal(t, "gpt-4", clonedEvaluation["model"])

	clonedMetrics := clonedEvaluation["metrics"].(map[string]any)
	clonedScores := clonedMetrics["scores"].([]any)
	clonedFirstScore := clonedScores[0].(map[string]any)
	assert.InEpsilon(t, 0.95, clonedFirstScore["value"], 0.0001)

	clonedContext := clonedScore.Metadata["context"].(map[string]any)
	clonedRefAnswers := clonedContext["reference_answers"].([]string)
	assert.Equal(t, "answer1", clonedRefAnswers[0])
}
