package redis

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

type Storable interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type Journal struct {
	rdb Storable
}

// NewJournal creates a new instance of the Journal.
func NewJournal(rdb Storable) *Journal {
	return &Journal{rdb: rdb}
}

// Stores pairs of POW challenges.
// Save stores a key-value pair in Redis with a TTL of 1 hour.
func (j *Journal) Save(ctx context.Context, key string, value string) error {
	ttl := time.Hour // 1 hour TTL
	err := j.rdb.Set(ctx, key, value, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

// Get retrieves the value associated with the given key.
func (j *Journal) Get(ctx context.Context, key string) (*string, error) {
	value, err := j.rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil // Key does not exist
	} else if err != nil {
		return nil, err // Some other error occurred
	}
	return &value, nil
}

// Delete removes a key-value pair from Redis.
func (j *Journal) Delete(ctx context.Context, key string) error {
	err := j.rdb.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}
