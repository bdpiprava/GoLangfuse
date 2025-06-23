package langfuse_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	langfuse "github.com/bdpiprava/GoLangfuse"
	"github.com/bdpiprava/GoLangfuse/config"
	"github.com/bdpiprava/GoLangfuse/mock"
	"github.com/bdpiprava/GoLangfuse/types"
)

func Test_AddEvent_ShouldCallClientToSendEvents(t *testing.T) {
	cfg := &config.Langfuse{
		URL:                    "http://localhost:3000",
		PublicKey:              "LangfusePublicKey",
		SecretKey:              "LangfuseSecretKey",
		NumberOfEventProcessor: 1,
	}
	httpClient := &http.Client{}
	eventID := uuid.MustParse("f8359e80-1ecd-471b-bf2a-49d2009a9179")
	mockTransport := mock.AddMockTransport(t, httpClient)

	resp := &http.Response{Body: io.NopCloser(strings.NewReader("{}"))}
	mockTransport.ExpectWith("POST", "http://localhost:3000/api/public/ingestion").Return(resp, nil)
	subject := langfuse.NewWithClient(cfg, httpClient)

	subject.AddEvent(context.TODO(), &types.TraceEvent{ID: &eventID, Name: "LLM"})

	assert.Eventually(t, func() bool {
		return mockTransport.AllExpectationMet()
	}, time.Second*10, time.Millisecond*100)

	body, err := io.ReadAll(mockTransport.RecordedRequests()[0].Body)
	assert.Nil(t, err)

	assert.Contains(t, string(body), `"id":"f8359e80-1ecd-471b-bf2a-49d2009a9179"`)
	assert.Contains(t, string(body), `"type":"trace-create"`)
	assert.Contains(t, string(body), `"name":"LLM"`)
	assert.Contains(t, string(body), `"public":false`)
}
