package langfuse_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	langfuse "github.com/bdpiprava/GoLangfuse"
	"github.com/bdpiprava/GoLangfuse/config"
	"github.com/bdpiprava/GoLangfuse/mock"
	"github.com/bdpiprava/GoLangfuse/types"
)

func Test_Send(t *testing.T) {
	cfg := &config.Langfuse{
		URL: "http://localhost:3000",
	}
	httpClient := &http.Client{}
	eventID := uuid.MustParse("f8359e80-1ecd-471b-bf2a-49d2009a9179")
	newClient := langfuse.NewClient(cfg, httpClient)
	traceID := "f8359e80-1ecd-471b-bf2a-49d2009a9179"

	testCases := []struct {
		name         string
		eventToSend  types.LangfuseEvent
		expectations func(*testing.T, error)
	}{
		{
			name:        "when try to send custom event of unknown type should fail with error",
			eventToSend: &CustomType{},
			expectations: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "cannot process event of 'unknown' type")
			},
		},
		{
			name:         "when try to send trace event should result in success",
			eventToSend:  &types.TraceEvent{ID: &eventID},
			expectations: func(t *testing.T, err error) { assert.NoError(t, err) },
		},
		{
			name:         "when try to send generation event should result in success",
			eventToSend:  &types.GenerationEvent{ID: &eventID},
			expectations: func(t *testing.T, err error) { assert.NoError(t, err) },
		},
		{
			name:         "when try to send span event should result in success",
			eventToSend:  &types.SpanEvent{ID: &eventID},
			expectations: func(t *testing.T, err error) { assert.NoError(t, err) },
		},
		{
			name:         "when try to send score event should result in success",
			eventToSend:  &types.ScoreEvent{ID: &eventID, Name: "example", Value: 0.9, TraceID: &traceID},
			expectations: func(t *testing.T, err error) { assert.NoError(t, err) },
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			response := &http.Response{Body: io.NopCloser(strings.NewReader("{}"))}
			mockTransport := mock.AddMockTransport(t, httpClient)
			mockTransport.ExpectWith("POST", "http://localhost:3000/api/public/ingestion").Return(response, nil)

			err := newClient.Send(context.TODO(), test.eventToSend)

			test.expectations(t, err)
		})
	}
}

func Test_Send_ValidateStruct(t *testing.T) {
	cfg := &config.Langfuse{
		URL: "http://localhost:3000",
	}
	httpClient := &http.Client{}
	newClient := langfuse.NewClient(cfg, httpClient)
	traceID := "10000000-0000-0000-0000-000000000001"

	testCases := []struct {
		name         string
		eventToSend  types.LangfuseEvent
		expectations func(*testing.T, error)
	}{
		{
			name:        "when name for score event is not provided results in error",
			eventToSend: &types.ScoreEvent{TraceID: &traceID, Value: 0.3},
			expectations: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "ingestion event validation failed: name: non zero value required")
			},
		},
		{
			name:        "when value for score event is not provided results in error",
			eventToSend: &types.ScoreEvent{TraceID: &traceID, Name: "score"},
			expectations: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "ingestion event validation failed: value: non zero value required")
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			err := newClient.Send(context.TODO(), test.eventToSend)

			test.expectations(t, err)
		})
	}
}

func Test_Send_WithInvalidLangfuseURL_Fails(t *testing.T) {

	testCases := []struct {
		name         string
		config       *config.Langfuse
		expectations func(*testing.T, error)
	}{
		{
			name:         "when langfuse url is not provided should fail",
			config:       &config.Langfuse{URL: ""},
			expectations: func(t *testing.T, err error) { assert.Contains(t, err.Error(), "missing langfuse config") },
		},
		{
			name:         "when langfuse url is only whitespaces should fail",
			config:       &config.Langfuse{URL: "    \n \t \r  "},
			expectations: func(t *testing.T, err error) { assert.Contains(t, err.Error(), "missing langfuse config") },
		},
		{
			name:   "when invalid url syntax for langfuse url should fail",
			config: &config.Langfuse{URL: "@@localhost:3000"},
			expectations: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), `failed to build langfuse url using @@localhost:3000 and /api/public/ingestion`)
			},
		},
	}

	httpClient := &http.Client{}
	eventID := uuid.MustParse("f8359e80-1ecd-471b-bf2a-49d2009a9179")

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			newClient := langfuse.NewClient(test.config, httpClient)

			err := newClient.Send(context.TODO(), &types.TraceEvent{ID: &eventID})

			test.expectations(t, err)
		})
	}
}

func Test_Send_WithResponse(t *testing.T) {
	cfg := &config.Langfuse{
		URL:       "http://localhost:3000",
		PublicKey: "LangfusePublicKey",
		SecretKey: "LangfuseSecretKey",
	}
	httpClient := &http.Client{}
	eventID := uuid.MustParse("f8359e80-1ecd-471b-bf2a-49d2009a9179")
	newClient := langfuse.NewClient(cfg, httpClient)

	testCases := []struct {
		name          string
		response      *http.Response
		responseError error
		expectations  func(*testing.T, error)
	}{
		{
			name:          "when received error from langfuse",
			responseError: fmt.Errorf("forced error"),
			expectations: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), `request to langfuse failed: Post "http://localhost:3000/api/public/ingestion": forced error`)
			},
		},
		{
			name:          "when received empty response from langfuse",
			responseError: nil,
			response:      &http.Response{},
			expectations: func(t *testing.T, err error) {
				assert.Contains(t, err.Error(), "failed to parse response: unexpected end of JSON input")
			},
		},
		{
			name:          "when received success response from langfuse",
			responseError: nil,
			response:      &http.Response{Body: io.NopCloser(strings.NewReader("{}"))},
			expectations: func(t *testing.T, err error) {
				assert.Nil(t, err)
			},
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			mockTransport := mock.AddMockTransport(t, httpClient)
			mockTransport.ExpectWith("POST", "http://localhost:3000/api/public/ingestion").Return(test.response, test.responseError)

			err := newClient.Send(context.TODO(), &types.TraceEvent{ID: &eventID})

			test.expectations(t, err)
		})
	}
}

type CustomType struct{}

func (c *CustomType) GetID() *uuid.UUID {
	id := uuid.New()
	return &id
}

func (c *CustomType) SetID(*uuid.UUID) {}
