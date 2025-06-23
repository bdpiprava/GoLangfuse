// Package langfuse provides http client for langfuse APIs
package langfuse

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"

	"github.com/bdpiprava/GoLangfuse/config"
	"github.com/bdpiprava/GoLangfuse/logger"
	"github.com/bdpiprava/GoLangfuse/types"
)

// Client a client interface for sending events to a Langfuse server
type Client interface {
	// Send sends ingestion event to langfuse using rest API
	Send(ctx context.Context, event types.LangfuseEvent) error
}

type client struct {
	client *http.Client
	config *config.Langfuse
}

// NewClient initialise new langfuse api client
func NewClient(
	config *config.Langfuse,
	httpClient *http.Client,
) Client {
	return &client{
		client: httpClient,
		config: config,
	}
}

// Send sends ingestion event to langfuse using rest API
func (c client) Send(ctx context.Context, ingestionEvent types.LangfuseEvent) error {
	log := logger.FromContext(ctx)
	if strings.TrimSpace(c.config.URL) == "" {
		log.Warn("langfuse config is not provided. no action is taken")
		return fmt.Errorf("missing langfuse config")
	}

	eventType := getEventType(ingestionEvent)
	if eventType == "unknown" {
		log.Errorf("cannot process event of 'unknown' type")
		return errors.Errorf("cannot process event of 'unknown' type")
	}

	if _, err := govalidator.ValidateStruct(ingestionEvent); err != nil {
		log.WithError(err).Errorf("ingestion event validation failed")
		return errors.Wrapf(err, "ingestion event validation failed")
	}

	request := &ingestionRequest{
		Batch: []event{
			{
				ID:        ingestionEvent.GetID().String(),
				Body:      ingestionEvent,
				Type:      eventType,
				Timestamp: time.Now(),
			},
		},
	}

	resp, err := c.sendEvent(ctx, request)
	if err != nil {
		return err
	}

	if len(resp.Errors) > 0 {
		log.Errorf("request to langfues returned errors in response %v", resp.Errors)
		return errors.Errorf("request to langfues returned errors in response %v", resp.Errors)
	}

	return err
}

// sendEvent send and ingestion event to langfuse
func (c client) sendEvent(ctx context.Context, request *ingestionRequest) (*ingestionResponse, error) {
	log := logger.FromContext(ctx)
	apiPath, err := url.JoinPath(c.config.URL, "/api/public/ingestion")
	if err != nil {
		log.WithError(err).Errorf("failed to build langfuse url using %s and /api/public/ingestion", c.config.URL)
		return nil, errors.Wrapf(err, "failed to build langfuse url using %s and /api/public/ingestion", c.config.URL)
	}

	payload, err := json.Marshal(request)
	if err != nil {
		log.WithError(err).Error("failed to marshal request payload")
		return nil, errors.Wrapf(err, "failed to marshal request payload")
	}

	httpRequest, err := http.NewRequest(http.MethodPost, apiPath, bytes.NewBuffer(payload))
	if err != nil {
		log.WithError(err).Error("failed to create langfuse request")
		return nil, errors.Wrapf(err, "failed to create langfuse request")
	}

	httpRequest.SetBasicAuth(c.config.PublicKey, c.config.SecretKey)
	httpRequest.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpRequest)
	if err != nil {
		log.WithError(err).Error("request to langfuse failed")
		return nil, errors.Wrapf(err, "request to langfuse failed")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("failed to read response")
		return nil, errors.Wrapf(err, "failed to read response")
	}

	var response ingestionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.WithError(err).Error("failed to parse response")
		return nil, errors.Wrapf(err, "failed to parse response")
	}
	return &response, nil
}

func getEventType(ingestionEvent types.LangfuseEvent) string {
	switch ingestionEvent.(type) {
	case *types.TraceEvent:
		return "trace-create"
	case *types.GenerationEvent:
		return "generation-create"
	case *types.SpanEvent:
		return "span-create"
	case *types.ScoreEvent:
		return "score-create"
	}
	return "unknown"
}
