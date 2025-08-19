package types

import (
	"github.com/google/uuid"
)

// ScoreEvent Create a score attached to a trace (and optionally an observation).
// Fields:
//   - ID The id of the score can be set, otherwise a random id is generated. Spans are upserted on id.
//   - Name identifier of the span. Useful for sorting/filtering in the UI.
//   - TraceID trace id where this span needs to be created.
//   - SessionID the id of the session to which the score should be attached.
//   - ObservationID the id of the observation to which the score should be attached.
//   - Value the value of the score. Can be any number, often standardized to 0..1.
//   - Comment Additional context/explanation of the score.
//   - DatasetRunID the id of the dataset run to which the score should be attached.
//   - Environment the environment in which the trace was created, e.g. "production", "staging", etc.
//   - Metadata of the span. it is merged when being updated via the API.
type ScoreEvent struct {
	ID            *uuid.UUID     `json:"id"`
	Name          string         `json:"name" valid:"required"`
	TraceID       *string        `json:"traceId"`
	SessionID     *string        `json:"sessionId,omitempty"`
	ObservationID *string        `json:"observationId,omitempty"`
	Value         float32        `json:"value" valid:"required"`
	Comment       *string        `json:"comment,omitempty"`
	DatasetRunID  *string        `json:"datasetRunId,omitempty"`
	Environment   *string        `json:"environment,omitempty"`
	Metadata      map[string]any `json:"metadata,omitempty"`
}

// GetID return an event ID
func (t *ScoreEvent) GetID() *uuid.UUID {
	return t.ID
}

// SetID set event ID
func (t *ScoreEvent) SetID(id *uuid.UUID) {
	t.ID = id
}

// Clone creates a deep copy of the ScoreEvent
func (t *ScoreEvent) Clone() LangfuseEvent {
	if t == nil {
		return nil
	}

	clone := &ScoreEvent{
		Name:  t.Name,
		Value: t.Value,
	}

	// Deep copy pointer fields
	if t.ID != nil {
		id := *t.ID
		clone.ID = &id
	}

	if t.TraceID != nil {
		traceID := *t.TraceID
		clone.TraceID = &traceID
	}

	if t.SessionID != nil {
		sessionID := *t.SessionID
		clone.SessionID = &sessionID
	}

	if t.ObservationID != nil {
		observationID := *t.ObservationID
		clone.ObservationID = &observationID
	}

	if t.Comment != nil {
		comment := *t.Comment
		clone.Comment = &comment
	}

	if t.DatasetRunID != nil {
		datasetRunID := *t.DatasetRunID
		clone.DatasetRunID = &datasetRunID
	}

	if t.Environment != nil {
		environment := *t.Environment
		clone.Environment = &environment
	}

	// Deep copy map
	if t.Metadata != nil {
		clone.Metadata = make(map[string]any, len(t.Metadata))
		for k, v := range t.Metadata {
			clone.Metadata[k] = deepCopyAny(v)
		}
	}

	return clone
}
