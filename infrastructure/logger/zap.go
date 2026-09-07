package logger

import (
	"usermanagement-api/domain/ports"

	"go.uber.org/zap"
)

// zapPortLogger satisfies ports.Logger on top of zap.
type zapPortLogger struct {
	l *zap.Logger
}

func NewPortLogger(l *zap.Logger) ports.Logger {
	return &zapPortLogger{l: l}
}

func (z *zapPortLogger) Info(msg string, args ...any) {
	z.l.Sugar().Infow(msg, args...)
}

func (z *zapPortLogger) Warn(msg string, args ...any) {
	z.l.Sugar().Warnw(msg, args...)
}

func (z *zapPortLogger) Error(msg string, args ...any) {
	z.l.Sugar().Errorw(msg, args...)
}
