// Package logger provides a simple logging interface for applications.
package logger

import (
	"context"
	"maps"

	"github.com/sirupsen/logrus"
)

var baseLogger = logrus.New()
var fieldsCtxKey any

// SetFieldsCtxKey sets the key to use for storing fields in the context.
func SetFieldsCtxKey(key any) {
	fieldsCtxKey = key
}

// FromContext returns an logger with all values from context loaded on it.
func FromContext(ctx context.Context) logrus.FieldLogger {
	log := baseLogger
	if fieldsCtxKey == nil {
		return log
	}

	return log.WithFields(getFields(ctx, fieldsCtxKey))
}

func getFields(ctx context.Context, key any) logrus.Fields {
	messageFields, ok := ctx.Value(key).(logrus.Fields)
	if !ok {
		return logrus.Fields{}
	}

	return messageFields
}

// WithFields returns a new context with the given fields, merged with the existing fields.
func WithFields(ctx context.Context, key any, fields logrus.Fields) context.Context {
	maps.Copy(fields, getFields(ctx, key))
	return context.WithValue(ctx, key, fields)
}
