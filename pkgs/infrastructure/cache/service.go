package cache

import (
	"context"
	"time"
)

type CachePayload struct {
	Key        string
	Value      []byte
	Expiration time.Duration
}

type Service interface {
	Store(ctx context.Context, assets []*CachePayload) error
	Retrieve(ctx context.Context, keys []string) ([]*CachePayload, error)
	ExistenceCheck(ctx context.Context, key string) (bool, error)
	ClearCacheByKeys(ctx context.Context, keys []string) error
	KeyIncrement(ctx context.Context, key string) (int64, error)
}
