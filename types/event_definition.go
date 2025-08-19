package types

import "github.com/google/uuid"

// LangfuseEvent interface representing langfuse event
type LangfuseEvent interface {
	// GetID return an event ID
	GetID() *uuid.UUID

	// SetID set event ID
	SetID(id *uuid.UUID)

	Clone() LangfuseEvent
}
