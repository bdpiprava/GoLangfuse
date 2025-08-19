package types

import (
	"time"

	"github.com/google/uuid"
)

// TraceEvent model representing langfuse trace
// Fields:
//   - ID an uuid for trace, auto-generated if not provided
//   - Name trace name
//   - UserID a user id to map traces to individual users
//   - SessionID a session id to map traces to specific session
//   - Release map trace with release
//   - Version map trace with  version
//   - Metadata of the trace
//   - Tags attach tags to the trace
//   - Public trace visibility, public or private, defaults to private
//   - Input an input to LLM
//   - Output an output from LLM
//   - Environment the environment in which the trace was created, e.g. "production", "staging", etc.
//   - Timestamp the time at which the trace was created
//   - ExternalID external ID for mapping to other systems
type TraceEvent struct {
	ID          *uuid.UUID     `json:"id" valid:"-"`
	Name        string         `json:"name" valid:"required"`
	UserID      string         `json:"userId,omitempty" valid:"-"`
	SessionID   string         `json:"sessionId,omitempty" valid:"-"`
	Release     string         `json:"release,omitempty" valid:"-"`
	Version     string         `json:"version,omitempty" valid:"-"`
	Metadata    map[string]any `json:"metadata,omitempty" valid:"-"`
	Tags        []string       `json:"tags,omitempty" valid:"-"`
	Public      bool           `json:"public" valid:"-"`
	Input       any            `json:"input,omitempty" valid:"-"`
	Output      any            `json:"output,omitempty" valid:"-"`
	Environment string         `json:"environment,omitempty" valid:"-"`
	Timestamp   *time.Time     `json:"timestamp,omitempty" valid:"-"`
	ExternalID  string         `json:"externalId,omitempty" valid:"-"`
}

// GetID return an event ID
func (t *TraceEvent) GetID() *uuid.UUID {
	return t.ID
}

// SetID set event ID
func (t *TraceEvent) SetID(id *uuid.UUID) {
	t.ID = id
}

// Clone creates a deep copy of the TraceEvent
func (t *TraceEvent) Clone() LangfuseEvent {
	if t == nil {
		return nil
	}

	clone := &TraceEvent{
		Name:        t.Name,
		UserID:      t.UserID,
		SessionID:   t.SessionID,
		Release:     t.Release,
		Version:     t.Version,
		Public:      t.Public,
		Environment: t.Environment,
		ExternalID:  t.ExternalID,
	}

	// Deep copy pointer fields
	if t.ID != nil {
		id := *t.ID
		clone.ID = &id
	}

	if t.Timestamp != nil {
		timestamp := *t.Timestamp
		clone.Timestamp = &timestamp
	}

	// Deep copy slice
	if t.Tags != nil {
		clone.Tags = make([]string, len(t.Tags))
		copy(clone.Tags, t.Tags)
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

// TraceBuilder provides a fluent interface for building TraceEvent
type TraceBuilder struct {
	trace *TraceEvent
}

// NewTrace creates a new TraceBuilder
func NewTrace(name string) *TraceBuilder {
	now := time.Now().UTC()
	return &TraceBuilder{
		trace: &TraceEvent{
			Name:      name,
			Timestamp: &now,
			Public:    false,
		},
	}
}

// WithID sets the trace ID
func (b *TraceBuilder) WithID(id uuid.UUID) *TraceBuilder {
	b.trace.ID = &id
	return b
}

// WithUserID sets the user ID
func (b *TraceBuilder) WithUserID(userID string) *TraceBuilder {
	b.trace.UserID = userID
	return b
}

// WithSessionID sets the session ID
func (b *TraceBuilder) WithSessionID(sessionID string) *TraceBuilder {
	b.trace.SessionID = sessionID
	return b
}

// WithRelease sets the release
func (b *TraceBuilder) WithRelease(release string) *TraceBuilder {
	b.trace.Release = release
	return b
}

// WithVersion sets the version
func (b *TraceBuilder) WithVersion(version string) *TraceBuilder {
	b.trace.Version = version
	return b
}

// WithMetadata sets the metadata
func (b *TraceBuilder) WithMetadata(metadata map[string]any) *TraceBuilder {
	b.trace.Metadata = metadata
	return b
}

// WithTags sets the tags
func (b *TraceBuilder) WithTags(tags ...string) *TraceBuilder {
	b.trace.Tags = tags
	return b
}

// WithPublic sets the public flag
func (b *TraceBuilder) WithPublic(public bool) *TraceBuilder {
	b.trace.Public = public
	return b
}

// WithInput sets the input
func (b *TraceBuilder) WithInput(input any) *TraceBuilder {
	b.trace.Input = input
	return b
}

// WithOutput sets the output
func (b *TraceBuilder) WithOutput(output any) *TraceBuilder {
	b.trace.Output = output
	return b
}

// WithEnvironment sets the environment
func (b *TraceBuilder) WithEnvironment(environment string) *TraceBuilder {
	b.trace.Environment = environment
	return b
}

// WithTimestamp sets the timestamp
func (b *TraceBuilder) WithTimestamp(timestamp time.Time) *TraceBuilder {
	b.trace.Timestamp = &timestamp
	return b
}

// WithExternalID sets the external ID
func (b *TraceBuilder) WithExternalID(externalID string) *TraceBuilder {
	b.trace.ExternalID = externalID
	return b
}

// Build returns the built TraceEvent
func (b *TraceBuilder) Build() *TraceEvent {
	return b.trace
}
