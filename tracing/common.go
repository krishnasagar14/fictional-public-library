package tracing

import (
	"context"
	"fictional-public-library/literals"
	"fictional-public-library/logging"
	"github.com/sirupsen/logrus"
)

func GetTracedLogEntry(ctx context.Context) *logrus.Entry {
	traceID := ctx.Value(literals.TraceIDContextKey)
	fields := logrus.Fields{}
	fields[literals.TraceIDContextKey] = traceID
	return logging.Log.WithFields(fields)
}
