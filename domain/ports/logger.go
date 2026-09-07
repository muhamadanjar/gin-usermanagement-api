package ports

// Logger is the minimal logging port used by application use cases.
// Implementations: infrastructure/logger (zap adapter).
type Logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}
