package types

import (
	"time"

	"github.com/google/uuid"
)

// SpanEvent A span represents durations of units of work in a trace.
// Usually, you want to add a span nested within a trace. Optionally you can nest it within another observation by
// providing a parent_observation_id. If no trace_id is provided, a new trace is created just for this span.
// Fields:
//   - ID The id of the span can be set, otherwise a random id is generated. Spans are upserted on id
//   - TraceID trace id where this span needs to be created
//   - ParentObservationID the ID of the parent observation, if applicable.
//   - Name identifier of the span. Useful for sorting/filtering in the UI.
//   - StartTime the time at which the span started, defaults to the current time.
//   - EndTime the time at which the span ended.
//   - Metadata of the span. It is merged when being updated via the API.
//   - Level the level of the generation. Used for sorting/filtering of traces with elevated error levels and for highlighting in the UI.
//   - StatusMessage the additional field for context of the event. E.g. the error message of an error event.
//   - Input the input to the span. Can be any JSON object.
//   - Output the output to the span. Can be any JSON object.
//   - Version the version of the span type. Used to understand how changes to the span type affect metrics. Useful in debugging.
//   - Environment the environment in which the trace was created, e.g. "production", "staging", etc.
type SpanEvent struct {
	ID                  *uuid.UUID     `json:"id"`
	TraceID             *uuid.UUID     `json:"traceId"`
	ParentObservationID *uuid.UUID     `json:"parentObservationId,omitempty"`
	Name                string         `json:"name"`
	StartTime           *time.Time     `json:"startTime,omitempty"`
	EndTime             *time.Time     `json:"endTime,omitempty"`
	Metadata            map[string]any `json:"metadata,omitempty"`
	Level               Level          `json:"level,omitempty"`
	StatusMessage       string         `json:"statusMessage,omitempty"`
	Input               any            `json:"input,omitempty"`
	Output              any            `json:"output,omitempty"`
	Version             string         `json:"version,omitempty"`
	Environment         string         `json:"environment,omitempty"`
}

// GetID return an event ID
func (t *SpanEvent) GetID() *uuid.UUID {
	return t.ID
}

// SetID set event ID
func (t *SpanEvent) SetID(id *uuid.UUID) {
	t.ID = id
}

// Clone creates a deep copy of the SpanEvent
func (t *SpanEvent) Clone() LangfuseEvent {
	if t == nil {
		return nil
	}

	clone := &SpanEvent{
		Name:          t.Name,
		Level:         t.Level,
		StatusMessage: t.StatusMessage,
		Version:       t.Version,
		Environment:   t.Environment,
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

	if t.ParentObservationID != nil {
		parentID := *t.ParentObservationID
		clone.ParentObservationID = &parentID
	}

	if t.StartTime != nil {
		startTime := *t.StartTime
		clone.StartTime = &startTime
	}

	if t.EndTime != nil {
		endTime := *t.EndTime
		clone.EndTime = &endTime
	}

	// Deep copy map
	if t.Metadata != nil {
		clone.Metadata = make(map[string]any, len(t.Metadata))
		for k, v := range t.Metadata {
			clone.Metadata[k] = deepCopyAny(v)
		}
	}

	// Deep copy any fields
	clone.Input = deepCopyAny(t.Input)
	clone.Output = deepCopyAny(t.Output)

	return clone
}

// Error set Level to error and EndTime with status message
func (t *SpanEvent) Error(statusMessage string) *SpanEvent {
	t.StatusMessage = statusMessage
	t.Level = Error
	return t.End()
}

// End set end time to now
func (t *SpanEvent) End() *SpanEvent {
	now := time.Now().UTC()
	t.EndTime = &now
	return t
}

// Start set start time to now
func (t *SpanEvent) Start() *SpanEvent {
	now := time.Now().UTC()
	t.StartTime = &now
	return t
}
