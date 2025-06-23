package langfuse

import (
	"time"

	"github.com/google/uuid"

	"github.com/bdpiprava/GoLangfuse/types"
)

// This file contains all types required to frame a request but does not require it to be exposed

// event an ingestion event to add trace, span, generation or score to langfuse
type event struct {
	ID        string              `json:"id"` // ID an event id
	Type      string              `json:"type"`
	Timestamp time.Time           `json:"timestamp"`
	Metadata  map[string]any      `json:"metadata,omitempty"`
	Body      types.LangfuseEvent `json:"body"`
}

// ingestionRequest langfuse ingestion request, only for client internal use
type ingestionRequest struct {
	Batch []event `json:"batch"`
}

// success langfuse response for the success cases
type success struct {
	ID     uuid.UUID `json:"id"`
	Status int       `json:"status"`
}

// eventError an error specific to event
type eventError struct {
	ID      uuid.UUID `json:"id"`
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Error   string    `json:"error"`
}

// ingestionResponse api call response from langfuse
type ingestionResponse struct {
	Successes []success    `json:"successes"`
	Errors    []eventError `json:"errors"`
}
