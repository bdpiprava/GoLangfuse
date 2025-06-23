package types

import (
	"github.com/google/uuid"
)

// TraceEvent model representing langfuse trace
type TraceEvent struct {
	ID          *uuid.UUID     `json:"id"`                     // ID an uuid for trace, auto-generated if not provided
	Name        string         `json:"name"`                   // Name trace name
	UserID      string         `json:"userId,omitempty"`       // UserID a user id to map traces to individual users
	SessionID   string         `json:"sessionId,omitempty"`    // SessionID a session id to map traces to specific session
	Release     string         `json:"release,omitempty"`      // Release map trace with release
	Version     string         `json:"version,omitempty"`      // Version map trace with  version
	Metadata    map[string]any `json:"metadata,omitempty"`     // Metadata metadata for the trace
	Tags        []string       `json:"tags,omitempty"`         // Tags attach tags to the trace
	Public      bool           `json:"public" default:"false"` // Public trace visibility, public or private, defaults to privet
	Input       any            `json:"input,omitempty"`        // Input an input to LLM
	Output      any            `json:"output,omitempty"`       // Output an output from LLM
	Environment string         `json:"environment,omitempty"`  // Environment the environment in which the trace was created, e.g. "production", "staging", etc.
}

// GetID return an event ID
func (t *TraceEvent) GetID() *uuid.UUID {
	return t.ID
}

// SetID set event ID
func (t *TraceEvent) SetID(id *uuid.UUID) {
	t.ID = id
}
