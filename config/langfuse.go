// Package config provides application config
package config

// Langfuse holds config for langfuse api calls
// URL langfuse server url
// PublicKey public key from langfuse project
// SecretKey secret key from langfuse project
type Langfuse struct {
	URL                    string `envconfig:"LANGFUSE_URL"`
	PublicKey              string `envconfig:"LANGFUSE_PUBLIC_KEY"`
	SecretKey              string `envconfig:"LANGFUSE_SECRET_KEY"`
	NumberOfEventProcessor int    `envconfig:"LANGFUSE_NUM_OF_EVENT_PROCESSOR" default:"1"`
}
