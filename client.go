// Package langfuse provides http client for langfuse APIs
package langfuse

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"

	"github.com/bdpiprava/GoLangfuse/config"
	"github.com/bdpiprava/GoLangfuse/logger"
	"github.com/bdpiprava/GoLangfuse/types"
)

const (
	eventTypeUnknown = "unknown"

	retryBackoffBase     = 2    // Base for exponential backoff calculation
	compressionThreshold = 1024 // Compress payload if > 1KB
	httpClientErrorStart = 400  // HTTP client error status codes start
)

// Client a client interface for sending events to a Langfuse server
type Client interface {
	// Send sends ingestion event to langfuse using rest API
	Send(ctx context.Context, event types.LangfuseEvent) error
	// SendBatch sends multiple events in a single batch to langfuse
	SendBatch(ctx context.Context, events []types.LangfuseEvent) error
}

type client struct {
	client *http.Client
	config *config.Langfuse
}

// NewOptimizedHTTPClient creates an HTTP client optimized for Langfuse API calls
func NewOptimizedHTTPClient(cfg *config.Langfuse) *http.Client {
	return &http.Client{
		Timeout: cfg.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:        cfg.MaxIdleConns,
			MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
			IdleConnTimeout:     cfg.IdleConnTimeout,
			DisableKeepAlives:   false,
			DisableCompression:  false,
		},
	}
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
		return ErrMissingURL
	}

	eventType := getEventType(ingestionEvent)
	if eventType == eventTypeUnknown {
		log.Errorf("cannot process event of 'unknown' type")
		return ErrUnknownEventType
	}

	if _, err := govalidator.ValidateStruct(ingestionEvent); err != nil {
		log.WithError(err).Errorf("ingestion event validation failed")
		return ErrEventValidation.WithCause(err)
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

	resp, err := c.sendEventWithRetry(ctx, request)
	if err != nil {
		return err
	}

	if len(resp.Errors) > 0 {
		log.Errorf("request to langfuse returned errors in response %v", resp.Errors)
		return ErrAPIServerError.WithDetails(map[string]any{
			"api_errors": resp.Errors,
		})
	}

	return nil
}

// SendBatch sends multiple events in a single batch to langfuse
func (c client) SendBatch(ctx context.Context, events []types.LangfuseEvent) error {
	log := logger.FromContext(ctx)
	if strings.TrimSpace(c.config.URL) == "" {
		log.Warn("langfuse config is not provided. no action is taken")
		return ErrMissingURL
	}

	if len(events) == 0 {
		return nil // Nothing to send
	}

	// Validate all events first
	var batchEvents []event
	for i, ingestionEvent := range events {
		eventType := getEventType(ingestionEvent)
		if eventType == eventTypeUnknown {
			log.Errorf("cannot process event of 'unknown' type")
			return ErrUnknownEventType.WithDetails(map[string]any{
				"event_index": i,
			})
		}

		if _, err := govalidator.ValidateStruct(ingestionEvent); err != nil {
			log.WithError(err).Errorf("ingestion event validation failed")
			return ErrEventValidation.WithCause(err).WithDetails(map[string]any{
				"event_index": i,
			})
		}

		batchEvents = append(batchEvents, event{
			ID:        ingestionEvent.GetID().String(),
			Body:      ingestionEvent,
			Type:      eventType,
			Timestamp: time.Now(),
		})
	}

	request := &ingestionRequest{
		Batch: batchEvents,
	}

	resp, err := c.sendEventWithRetry(ctx, request)
	if err != nil {
		return err
	}

	if len(resp.Errors) > 0 {
		log.Errorf("request to langfuse returned errors in response %v", resp.Errors)
		return ErrBatchProcessing.WithDetails(map[string]any{
			"api_errors": resp.Errors,
			"batch_size": len(events),
		})
	}

	return nil
}

// sendEventWithRetry sends an ingestion event to langfuse with retry logic
func (c client) sendEventWithRetry(ctx context.Context, request *ingestionRequest) (*ingestionResponse, error) {
	var lastErr error

	for i := 0; i <= c.config.MaxRetries; i++ {
		if i > 0 {
			// Calculate exponential backoff delay
			delay := time.Duration(math.Pow(retryBackoffBase, float64(i-1))) * c.config.RetryDelay
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		resp, err := c.sendEvent(ctx, request)
		if err == nil {
			return resp, nil
		}

		lastErr = err
		log := logger.FromContext(ctx)

		// Don't retry on client errors (4xx) or non-retryable errors
		var langfuseErr *Error
		if errors.As(err, &langfuseErr) {
			if !langfuseErr.IsRetryable() {
				log.WithError(err).Errorf("non-retryable error, not retrying: %v", err)
				break
			}
		}

		if i < c.config.MaxRetries {
			log.WithError(err).Warnf("request failed, retrying (attempt %d/%d)", i+1, c.config.MaxRetries)
		}
	}

	return nil, ErrRequestFailed.WithCause(lastErr).WithDetails(map[string]any{
		"max_retries": c.config.MaxRetries,
	})
}

// sendEvent send and ingestion event to langfuse
func (c client) sendEvent(ctx context.Context, request *ingestionRequest) (*ingestionResponse, error) {
	log := logger.FromContext(ctx)
	apiPath, err := url.JoinPath(c.config.URL, "/api/public/ingestion")
	if err != nil {
		log.WithError(err).Errorf("failed to build langfuse url using %s and /api/public/ingestion", c.config.URL)
		return nil, ErrInvalidConfig.WithCause(err).WithDetails(map[string]any{
			"url": c.config.URL,
		})
	}

	payload, err := json.Marshal(request)
	if err != nil {
		log.WithError(err).Error("failed to marshal request payload")
		return nil, ErrEventProcessing.WithCause(err)
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, apiPath, bytes.NewBuffer(payload))
	if err != nil {
		log.WithError(err).Error("failed to create langfuse request")
		return nil, ErrRequestFailed.WithCause(err)
	}

	httpRequest.SetBasicAuth(c.config.PublicKey, c.config.SecretKey)
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Accept-Encoding", "gzip")

	resp, err := c.client.Do(httpRequest)
	if err != nil {
		log.WithError(err).Error("request to langfuse failed")
		return nil, ErrConnectionFailed.WithCause(err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	// Handle HTTP errors
	if resp.StatusCode >= httpClientErrorStart {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, NewHTTPError(resp.StatusCode, string(bodyBytes))
	}

	// Handle compressed response
	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			log.WithError(err).Error("failed to create gzip reader")
			return nil, ErrEventProcessing.WithCause(err).WithDetails(map[string]any{
				"operation": "decompression",
			})
		}
		defer gzReader.Close()
		reader = gzReader
	}

	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		log.WithError(err).Error("failed to read response")
		return nil, ErrNetworkTimeout.WithCause(err)
	}

	var response ingestionResponse
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		log.WithError(err).Error("failed to parse response")
		return nil, ErrEventProcessing.WithCause(err).WithDetails(map[string]any{
			"operation":     "json_unmarshal",
			"response_body": string(bodyBytes),
		})
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
	return eventTypeUnknown
}
