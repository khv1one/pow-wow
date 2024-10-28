package redis_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	rdb "github.com/littlebugger/pow-wow/internal/service/repository/redis"
)

// MockStorable is a mock implementation of the Storable interface
type MockStorable struct {
	mock.Mock
}

func (m *MockStorable) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockStorable) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockStorable) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return args.Get(0).(*redis.IntCmd)
}

func TestJournal_Save(t *testing.T) {
	ctx := context.Background()
	mockRedis := new(MockStorable)
	journal := rdb.NewJournal(mockRedis)

	// Define the key and value
	key := "test-key"
	value := "test-value"
	ttl := time.Hour

	// Mock successful Set
	mockRedis.On("Set", ctx, key, value, ttl).Return(redis.NewStatusCmd(ctx)).Once()

	// Execute
	err := journal.Save(ctx, key, value)

	// Assertions
	assert.NoError(t, err)
	mockRedis.AssertExpectations(t)
}

func TestJournal_Get(t *testing.T) {
	ctx := context.Background()
	mockRedis := new(MockStorable)
	journal := rdb.NewJournal(mockRedis)

	// Define the key and value
	key := "test-key"
	value := "test-value"

	// Define scenarios
	t.Run("successful retrieval", func(t *testing.T) {
		// Mock successful Get
		stringCmd := redis.NewStringCmd(ctx)
		stringCmd.SetVal(value)
		mockRedis.On("Get", ctx, key).Return(stringCmd).Once()

		// Execute
		result, err := journal.Get(ctx, key)

		// Assertions
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, &value, result)
		mockRedis.AssertExpectations(t)
	})

	t.Run("key does not exist", func(t *testing.T) {
		// Mock key not found (redis.Nil)
		stringCmd := redis.NewStringCmd(ctx)
		stringCmd.SetErr(redis.Nil)
		mockRedis.On("Get", ctx, key).Return(stringCmd).Once()

		// Execute
		result, err := journal.Get(ctx, key)

		// Assertions
		assert.NoError(t, err)
		assert.Nil(t, result)
		mockRedis.AssertExpectations(t)
	})

	t.Run("redis error", func(t *testing.T) {
		// Mock Redis error
		stringCmd := redis.NewStringCmd(ctx)
		stringCmd.SetErr(errors.New("redis error"))
		mockRedis.On("Get", ctx, key).Return(stringCmd).Once()

		// Execute
		result, err := journal.Get(ctx, key)

		// Assertions
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRedis.AssertExpectations(t)
	})
}

func TestJournal_Delete(t *testing.T) {
	ctx := context.Background()
	mockRedis := new(MockStorable)
	journal := rdb.NewJournal(mockRedis)

	// Define the key
	key := "test-key"

	// Mock successful Del
	intCmd := redis.NewIntCmd(ctx)
	intCmd.SetVal(1) // Return 1 as the number of deleted keys
	mockRedis.On("Del", ctx, []string{key}).Return(intCmd).Once()

	// Execute
	err := journal.Delete(ctx, key)

	// Assertions
	assert.NoError(t, err)
	mockRedis.AssertExpectations(t)
}
