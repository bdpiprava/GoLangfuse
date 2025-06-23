package types

import (
	"github.com/google/uuid"
)

// ScoreEvent Create a score attached to a trace (and optionally an observation).
type ScoreEvent struct {
	ID            *uuid.UUID     `json:"id"`                      // ID The id of the score can be set, otherwise a random id is generated. Spans are upserted on id.
	Name          string         `json:"name" valid:"required"`   // Name identifier of the span. Useful for sorting/filtering in the UI.
	TraceID       *string        `json:"traceId"`                 // TraceID trace id where this span needs to be created.
	SessionID     *string        `json:"sessionId,omitempty"`     // SessionID the id of the session to which the score should be attached.
	ObservationID *string        `json:"observationId,omitempty"` // ObservationID the id of the observation to which the score should be attached.
	Value         float32        `json:"value" valid:"required"`  // Value the value of the score. Can be any number, often standardized to 0..1.
	Comment       *string        `json:"comment,omitempty"`       // Comment Additional context/explanation of the score.
	DatasetRunID  *string        `json:"datasetRunId,omitempty"`  // DatasetRunID the id of the dataset run to which the score should be attached.
	Environment   *string        `json:"environment,omitempty"`   // Environment the environment in which the trace was created, e.g. "production", "staging", etc.
	Metadata      map[string]any `json:"metadata,omitempty"`      // Metadata additional metadata of the span. Can be any JSON object. Metadata is merged when being updated via the API.
}

// GetID return an event ID
func (t *ScoreEvent) GetID() *uuid.UUID {
	return t.ID
}

// SetID set event ID
func (t *ScoreEvent) SetID(id *uuid.UUID) {
	t.ID = id
}
