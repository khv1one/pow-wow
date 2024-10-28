package usecase_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/littlebugger/pow-wow/internal/service/entity"
	"github.com/littlebugger/pow-wow/internal/service/usecase"
	"github.com/littlebugger/pow-wow/internal/service/usecase/mocks"
	"github.com/stretchr/testify/assert"
)

func TestJournal_Score(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockStorage := mocks.NewMockTrackable(t)
	journal := usecase.NewJournal(mockStorage)

	// Define a sample challenge and remark UUID
	remark := uuid.New()
	challenge := entity.Challenge{
		// Fill in challenge attributes as necessary
	}

	// Mock successful storage save
	value, _ := json.Marshal(challenge)
	mockStorage.On("Save", ctx, remark.String(), string(value)).Return(nil).Once()

	// Execute
	err := journal.Score(ctx, remark, challenge)

	// Assertions
	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestJournal_Match(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockStorage := mocks.NewMockTrackable(t)
	journal := usecase.NewJournal(mockStorage)

	// Define a sample challenge and remark UUID
	remark := uuid.New()
	challenge := entity.Challenge{
		// Fill in challenge attributes as necessary
	}
	storedValue, _ := json.Marshal(challenge)

	// Define scenarios
	t.Run("successful match", func(t *testing.T) {
		// Mock successful get
		mockStorage.On("Get", ctx, remark.String()).Return(toStringPointer(string(storedValue)), nil).Once()

		// Execute
		result, err := journal.Match(ctx, remark)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, challenge, result)
		mockStorage.AssertExpectations(t)
	})

	t.Run("no such challenge", func(t *testing.T) {
		// Mock nil result
		mockStorage.On("Get", ctx, remark.String()).Return(nil, nil).Once()

		// Execute
		result, err := journal.Match(ctx, remark)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, entity.Challenge{}, result)
		assert.Contains(t, err.Error(), "no such challenge")
		mockStorage.AssertExpectations(t)
	})

	t.Run("corrupted challenge", func(t *testing.T) {
		// Mock corrupted data
		mockStorage.On("Get", ctx, remark.String()).Return(toStringPointer("corrupted"), nil).Once()

		// Execute
		result, err := journal.Match(ctx, remark)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, entity.Challenge{}, result)
		assert.Contains(t, err.Error(), "challenge was corrupted")
		mockStorage.AssertExpectations(t)
	})
}

func TestJournal_MarkSolved(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockStorage := mocks.NewMockTrackable(t)
	journal := usecase.NewJournal(mockStorage)

	// Define a sample remark UUID
	remark := uuid.New()

	// Mock successful delete
	mockStorage.On("Delete", ctx, remark.String()).Return(nil).Once()

	// Execute
	err := journal.MarkSolved(ctx, remark)

	// Assertions
	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func toStringPointer(s string) *string {
	return &s
}
