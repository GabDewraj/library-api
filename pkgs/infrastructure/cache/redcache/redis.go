package redcache

import (
	"context"

	"github.com/GabDewraj/library-api/pkgs/infrastructure/cache"
	"github.com/redis/go-redis/v9"
)

type service struct {
	rdb *redis.Client
}

func NewRedisCache(client *redis.Client) cache.Service {
	return &service{
		rdb: client,
	}
}

// ClearCacheByKeys implements Service.
func (s *service) ClearCacheByKeys(ctx context.Context, keys []string) error {
	return s.rdb.Del(ctx, keys...).Err()
}

func (s *service) ExistenceCheck(ctx context.Context, key string) (bool, error) {
	exists, err := s.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if exists == 1 {
		return true, nil
	}
	return false, nil
}

func (s *service) Store(ctx context.Context, payload []*cache.CachePayload) error {
	values := map[string]string{}
	for _, asset := range payload {
		values[asset.Key] = string(asset.Value)
	}
	// Set the values in the cache
	if err := s.rdb.MSet(ctx, values).Err(); err != nil {
		return err
	}

	// Set the expiration times
	for _, asset := range payload {
		if err := s.rdb.Expire(context.Background(),
			asset.Key, asset.Expiration).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) Retrieve(ctx context.Context, keys []string) ([]*cache.CachePayload, error) {

	// Retrieve the serialized strings
	jsonStrings, err := s.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	var data []*cache.CachePayload
	for index, jsonString := range jsonStrings {
		if jsonString != nil {
			retrievedAsset := &cache.CachePayload{
				Key:   keys[index],
				Value: []byte(jsonString.(string)),
			}
			data = append(data, retrievedAsset)
		}
	}
	return data, nil
}

// KeyIncrement implements cache.Service.
func (s *service) KeyIncrement(ctx context.Context, key string) (int64, error) {
	return s.rdb.Incr(ctx, key).Result()
}
