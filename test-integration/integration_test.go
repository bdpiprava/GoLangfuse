package integration_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/suite"

	langfuse "github.com/bdpiprava/GoLangfuse"
	"github.com/bdpiprava/GoLangfuse/config"
	"github.com/bdpiprava/GoLangfuse/types"
	"github.com/bdpiprava/easy-http/pkg/httpx"
)

// LangfuseIntegrationTestSuite is a test suite for Langfuse integration tests.
type LangfuseIntegrationTestSuite struct {
	suite.Suite
	subject langfuse.Langfuse
	cfg     *config.Langfuse
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(LangfuseIntegrationTestSuite))
}

func (s *LangfuseIntegrationTestSuite) SetupSuite() {
	err := godotenv.Load(".integration.env")
	s.Require().NoError(err, "Failed to load .integration.env file")

	s.cfg, err = loadEnv()
	s.Require().NoError(err, "Failed to load configuration from environment variables")

	s.subject = langfuse.New(s.cfg)
}

func (s *LangfuseIntegrationTestSuite) Test_SendTrace() {
	traceEvent := NewTestTraceEvent()
	expectedSession := BuildSession(traceEvent)

	s.subject.AddEvent(context.TODO(), traceEvent)

	s.Eventually(func() bool {
		tracePath := fmt.Sprintf("/api/public/traces/%s", traceEvent.ID.String())
		got, err := getEventType[types.TraceEvent](s.cfg, tracePath)
		if err != nil {
			s.T().Logf("Failed to get trace event: %v", err)
			return false
		}

		return s.Equal(traceEvent, got, cmpopts.IgnoreFields(types.TraceEvent{}, "Timestamp"))
	}, 10*time.Second, 100*time.Millisecond, "Trace event should be sent successfully")

	s.Eventually(func() bool {
		sessionPath := fmt.Sprintf("/api/public/sessions/%s", traceEvent.SessionID)
		got, err := getEventType[types.Session](s.cfg, sessionPath)
		if err != nil {
			s.T().Logf("Failed to get session: %v", err)
			return false
		}

		return s.Equal(expectedSession, got,
			cmpopts.IgnoreFields(types.Trace{}, "Timestamp", "CreatedAt", "UpdatedAt"),
			cmpopts.IgnoreFields(types.Session{}, "CreatedAt"),
		)
	}, 10*time.Second, 100*time.Millisecond, "Trace event should be sent successfully")
}

func (s *LangfuseIntegrationTestSuite) Test_SendGeneration() {
	generationEvent := NewTestGenerationEvent()

	s.subject.AddEvent(context.TODO(), generationEvent)

	s.Eventually(func() bool {
		observationPath := fmt.Sprintf("/api/public/observations/%s", generationEvent.ID.String())
		got, err := getEventType[types.GenerationEvent](s.cfg, observationPath)
		if err != nil {
			s.T().Logf("Failed to get generation event: %v", err)
			return false
		}

		return s.Equal(generationEvent, got,
			cmpopts.IgnoreFields(types.GenerationEvent{}, "StartTime", "CompletionStartTime", "EndTime"))
	}, 10*time.Second, 100*time.Millisecond, "Generation event should be sent successfully")
}

func (s *LangfuseIntegrationTestSuite) Test_SendSpan() {
	spanEvent := NewTestSpanEvent()

	s.subject.AddEvent(context.TODO(), spanEvent)

	s.Eventually(func() bool {
		observationPath := fmt.Sprintf("/api/public/observations/%s", spanEvent.ID.String())
		got, err := getEventType[types.SpanEvent](s.cfg, observationPath)
		if err != nil {
			s.T().Logf("Failed to get span event: %v", err)
			return false
		}

		return s.Equal(spanEvent, got,
			cmpopts.IgnoreFields(types.SpanEvent{}, "StartTime", "EndTime"))
	}, 10*time.Second, 100*time.Millisecond, "Span event should be sent successfully")
}

func (s *LangfuseIntegrationTestSuite) Test_SendScore() {
	scoreEvent := NewTestScoreEvent()

	s.subject.AddEvent(context.TODO(), scoreEvent)

	s.Eventually(func() bool {
		scorePath := fmt.Sprintf("/api/public/scores/%s", scoreEvent.ID.String())
		got, err := getEventType[types.ScoreEvent](s.cfg, scorePath)
		if err != nil {
			s.T().Logf("Failed to get score event: %v", err)
			return false
		}

		return s.Equal(scoreEvent, got)
	}, 10*time.Second, 100*time.Millisecond, "Score event should be sent successfully")
}

func (s *LangfuseIntegrationTestSuite) Equal(expected, got any, opts ...cmp.Option) bool {
	less := func(a, b string) bool { return a < b }
	opts = append(opts, cmpopts.SortSlices(less))

	diff := cmp.Diff(expected, got, opts...)
	if diff != "" {
		s.T().Logf("Trace event mismatch (-want +got):\n%s", diff)
		return false
	}

	return true
}

func getEventType[T any](cfg *config.Langfuse, path string) (*T, error) {
	r, err := httpx.GET[string](
		httpx.WithPath(path),
		httpx.WithBaseURL(cfg.URL),
		httpx.WithHeader("Authorization", "Basic "+basicAuth(cfg.PublicKey, cfg.SecretKey)),
		httpx.WithHeader("Content-Type", "application/json"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get event type: %w", err)
	}

	if r.StatusCode == 404 {
		return nil, fmt.Errorf("event not found")
	}

	var t T
	err = json.Unmarshal(r.RawBody, &t)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w, body: %s", err, string(r.RawBody))
	}
	return &t, err
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// loadEnv load config variables into config.Langfuse.
func loadEnv() (*config.Langfuse, error) {
	var cfg config.Langfuse
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, err
}
