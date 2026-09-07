package ports

import (
	"context"
	"time"
)

// Cache is the port for the key-value store used by application use cases.
// Implementations: pkg/cache.RedisCache (structural typing, no adapter needed).
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
}
