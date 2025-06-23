// Package langfuse provides an interface to send ingestion events to Langfuse in an asynchronous manner.
package langfuse

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/bdpiprava/GoLangfuse/config"
	"github.com/bdpiprava/GoLangfuse/logger"
	"github.com/bdpiprava/GoLangfuse/types"
)

// Langfuse an interface to send ingestion events to langfuse in async manner
// Event is added to the queue and then processor is sending it to the langfuse
type Langfuse interface {
	// AddEvent adds event to the channel and returns the event unique ID, generating one if missing.
	AddEvent(ctx context.Context, event types.LangfuseEvent) *uuid.UUID
}

type eventChanItem struct {
	ctx   context.Context
	event types.LangfuseEvent
}

type langfuseService struct {
	client       Client
	eventChannel chan eventChanItem
}

// New initialise new Langfuse instance for given config with background event processors
func New(config *config.Langfuse) Langfuse {
	return NewWithClient(config, http.DefaultClient)
}

// NewWithClient initialise new Langfuse instance with background event processors
func NewWithClient(config *config.Langfuse, customHttpClient *http.Client) Langfuse {
	eventManager := &langfuseService{
		client:       NewClient(config, customHttpClient),
		eventChannel: make(chan eventChanItem, 512),
	}
	eventManager.startEventProcessors(config.NumberOfEventProcessor)
	return eventManager
}

// AddEvent adds event to the channel and returns the event unique ID, generating one if missing.
func (l *langfuseService) AddEvent(ctx context.Context, event types.LangfuseEvent) *uuid.UUID {
	ensureEventID(event)
	l.eventChannel <- eventChanItem{ctx: ctx, event: event}
	return event.GetID()
}

// startEventProcessors start the background event processors
func (l *langfuseService) startEventProcessors(count int) {
	for i := 0; i < count; i++ {
		go func() {
			for item := range l.eventChannel {
				l.send(item.ctx, item.event)
			}
		}()
	}
}

// send sends a types.LangfuseEvent to the Langfuse and logs any issues.
func (l *langfuseService) send(ctx context.Context, event types.LangfuseEvent) {
	log := logger.FromContext(ctx)
	log.Debugf("sending event to langfuse %v", event)

	err := l.client.Send(ctx, event)
	if err != nil {
		log.WithError(err).Errorf("failed to send event %v", event)
		return
	}
}

// ensureEventID ensures that the IngestionEvent has a unique ID, generating one if missing.
func ensureEventID(ingestionEvent types.LangfuseEvent) {
	if ingestionEvent.GetID() != nil {
		return
	}
	newID := uuid.New()
	ingestionEvent.SetID(&newID)
}
