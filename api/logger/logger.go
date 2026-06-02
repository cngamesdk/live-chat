package logger

import (
	"context"

	"go.uber.org/zap"
)

type ctxKeyType struct{}

var ctxKey ctxKeyType

var baseLogger *zap.Logger

func Init(l *zap.Logger) {
	baseLogger = l
}

// L returns the base logger without any context.
func L() *zap.Logger {
	return baseLogger
}

// Ctx returns a logger with trace_id from context, or base logger if absent.
func Ctx(ctx context.Context) *zap.Logger {
	if ctx == nil || baseLogger == nil {
		return baseLogger
	}
	traceID, _ := ctx.Value(ctxKey).(string)
	if traceID == "" {
		return baseLogger
	}
	return baseLogger.With(zap.String("trace_id", traceID))
}

// WithTraceID stores traceID in context for later log correlation.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, ctxKey, traceID)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	Ctx(ctx).Info(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	Ctx(ctx).Error(msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	Ctx(ctx).Warn(msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	Ctx(ctx).Fatal(msg, fields...)
}
