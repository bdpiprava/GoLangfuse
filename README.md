# 🚀 GoLangfuse

[![Go Report Card](https://goreportcard.com/badge/github.com/bdpiprava/GoLangfuse)](https://goreportcard.com/report/github.com/bdpiprava/GoLangfuse)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Reference](https://pkg.go.dev/badge/github.com/bdpiprava/GoLangfuse.svg)](https://pkg.go.dev/github.com/bdpiprava/GoLangfuse)

**Supercharge your Go LLM applications with seamless observability!** GoLangfuse is a lightning-fast, asynchronous client library for the Langfuse API that lets you track, analyze, and optimize your AI applications with minimal overhead.

## 🔍 About Langfuse

[Langfuse](https://langfuse.com) is an open-source observability platform built specifically for LLM applications. Think of it as your AI app's flight recorder and control panel all in one! Monitor performance, detect issues, evaluate quality, and optimize your LLM applications with powerful analytics.

## ✨ Why GoLangfuse?

Built for Go developers who demand both performance and simplicity, GoLangfuse offers:

- **⚡ Fire-and-Forget Events** - Non-blocking, asynchronous event tracking that won't slow down your application
- **🧵 Background Processing** - Configurable number of worker goroutines to handle event processing
- **🎯 Complete Coverage** - Track every aspect of your LLM interactions:
  - 📝 Traces - Log complete user sessions and requests
  - ⏱️ Spans - Measure and track units of work
  - 🤖 Generations - Capture LLM inputs, outputs, and performance
  - 📊 Scores - Evaluate and track quality metrics
- **🛠️ Developer-Friendly** - Automatic ID generation, customizable HTTP client, built-in validation
- **🧩 Easy Integration** - Drop-in solution for any Go application using LLMs

## 📦 Installation

One line and you're ready to go:

```bash
go get github.com/bdpiprava/GoLangfuse
```

## 🚀 Quick Start

### Basic Setup

```go
package main

import (
	"context"

	"github.com/bdpiprava/GoLangfuse"
	"github.com/bdpiprava/GoLangfuse/config"
	"github.com/bdpiprava/GoLangfuse/types"
)

func main() {
	// Create a Langfuse client in one line! ⚡
	langfuseClient := langfuse.New(
		&config.Langfuse{
			URL:                    "https://api.langfuse.com", // Or your self-hosted instance
			PublicKey:              "YOUR_PUBLIC_KEY",
			SecretKey:              "YOUR_SECRET_KEY",
			NumberOfEventProcessor: 2, // Parallel workers for event processing
		},
	)

	// Create and track a trace for user interaction
	traceEvent := &types.TraceEvent{
		Name:      "user-question",
		UserID:    "user-123",
		SessionID: "session-456",
		Tags:      []string{"production", "experiment-1"},
		Metadata: map[string]interface{}{
			"source":  "mobile-app",
			"version": "2.1.3",
		},
		Input: map[string]interface{}{
			"prompt": "Tell me a joke about programming",
		},
		Output: map[string]interface{}{
			"completion": "Why do programmers prefer dark mode? Because light attracts bugs!",
		},
	}

	// Send it off asynchronously - keep your application responsive!
	ctx := context.Background()
	eventID := langfuseClient.AddEvent(ctx, traceEvent)
	
	// That's it! Your event is being processed in the background
	// while your application continues running at full speed
}
```

## 🎭 Real-world Use Cases

### Tracking Complete User Interactions

```go
// Start by creating a trace for the entire user session
traceEvent := &types.TraceEvent{
	Name:      "customer-support-session",
	UserID:    userId,
	SessionID: sessionId,
	Tags:      []string{"support", "tier-1"},
}
traceID := langfuseClient.AddEvent(ctx, traceEvent)

// Now track the LLM generation - record exactly what happens
genEvent := &types.GenerationEvent{
	TraceID: traceID.String(),
	Name:    "initial-response",
	Model:   "gpt-4",
	Input: map[string]interface{}{
		"prompt": "Customer is asking about refund policy. Provide a helpful response.",
		"temperature": 0.7,
		"max_tokens": 500,
	},
	Output: map[string]interface{}{
		"completion": modelResponse,
		"finish_reason": "stop",
	},
	StartTime: startTime,
	EndTime:   endTime,
	// Track costs and token usage automatically
	Usage: types.Usage{
		Input:  150,
		Output: 420,
		Total:  570,
		Unit:   types.Tokens,
		InputCost: 0.0003,
		OutputCost: 0.0024,
		TotalCost: 0.0027,
	},
}
langfuseClient.AddEvent(ctx, genEvent)

// Record quality scores to measure and improve over time
scoreEvent := &types.ScoreEvent{
	TraceID: traceID.String(),
	Name:    "helpfulness",
	Value:   0.95,
	Comment: "Response addressed customer concern completely and accurately",
}
langfuseClient.AddEvent(ctx, scoreEvent)
```

### Performance Optimization with Custom HTTP Client

Need to fine-tune network performance? No problem:

```go
// Create a custom HTTP client with optimized settings
customClient := &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     90 * time.Second,
	},
}

// Inject your custom client into Langfuse
langfuseClient := langfuse.NewWithClient(
	&config.Langfuse{
		URL:                    "https://api.langfuse.com",
		PublicKey:              "YOUR_PUBLIC_KEY",
		SecretKey:              "YOUR_SECRET_KEY",
		NumberOfEventProcessor: 4, // Increase for high-volume applications
	},
	customClient,
)
```

## 📊 Event Types Guide

### Trace - The Big Picture

Traces represent entire user interactions or requests - the perfect container for your LLM application flows:

```go
traceEvent := &types.TraceEvent{
	Name:        "product-recommendation",
	UserID:      "customer-456",
	SessionID:   "browsing-789",
	Tags:        []string{"ecommerce", "recommendations"},
	Metadata:    map[string]any{"category": "electronics", "items_viewed": 7},
	Public:      false, // Control data visibility
	Input:       userBrowsingHistory,
	Output:      recommendedProducts,
	Environment: "production",
}
```

### Span - Measure Process Steps

Track each step of your application's processing pipeline:

```go
spanEvent := &types.SpanEvent{
	TraceID:   traceID.String(),
	Name:      "product-matching-algorithm",
	StartTime: startTime,
	EndTime:   endTime,
	Input:     userPreferences,
	Output:    matchedProducts,
	Metadata:  map[string]any{"algorithm_version": "v3.2", "confidence": 0.87},
	// Track errors when they happen
	Level:         types.Error,
	StatusMessage: "Failed to retrieve product inventory data",
}
```

### Generation - LLM Interaction Details

Zero in on individual LLM API calls:

```go
genEvent := &types.GenerationEvent{
	TraceID:   traceID.String(),
	Name:      "product-description",
	Model:     "gpt-4",
	StartTime: startTime,
	EndTime:   endTime,
	Input:     "Generate a compelling description for this smartphone",
	Output:    generatedDescription,
	Metadata:  map[string]any{"prompt_tokens": 57, "completion_tokens": 210},
	// Track streaming performance
	CompletionStartTime: streamStartTime,
}
```

### Score - Quality Metrics

Create objective or subjective measures of your system's performance:

```go
scoreEvent := &types.ScoreEvent{
	TraceID: traceID.String(),
	Name:    "conversion-rate",
	Value:   0.123, // 12.3% conversion
	Comment: "Product description led to purchase",
}
```

## 🧰 Development

### Prerequisites

- Go 1.24 or higher

### Running Tests

Run the comprehensive test suite with:

```bash
make tests
```

### Building

Build the project with:

```bash
make build
```

## 🤝 Contributing

Join us in making LLM observability in Go even better! Contributions are enthusiastically welcomed:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

---

Built with ❤️ for the Go and LLM communities. Happy observing!