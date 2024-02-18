package cache

import (
	"context"
	"time"
)

type CacheJsonPayload struct {
	Key        string
	Value      []byte
	Expiration time.Duration
}
type CacheIntegerPayload struct {
	Key        string
	Value      int
	Expiration time.Duration
}

type Service interface {
	StoreInteger(ctx context.Context, asset CacheIntegerPayload) error
	StoreJSON(ctx context.Context, assets []*CacheJsonPayload) error
	RetrieveJSON(ctx context.Context, keys []string) ([]*CacheJsonPayload, error)
	RetrieveInteger(ctx context.Context, key string) (int, error)
	ExistenceCheck(ctx context.Context, key string) (bool, error)
	ClearCacheByKeys(ctx context.Context, keys []string) error
	KeyIncrement(ctx context.Context, key string) (int64, error)
}
