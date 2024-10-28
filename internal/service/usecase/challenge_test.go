package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/littlebugger/pow-wow/internal/service/entity"
	"github.com/littlebugger/pow-wow/internal/service/usecase"
	"github.com/littlebugger/pow-wow/internal/service/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChallenger_MakeChallenge(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockMinter := mocks.NewMockMinter(t)
	mockRecordable := mocks.NewMockRecordable(t)
	challenger := usecase.NewChallenger(mockMinter, mockRecordable)

	// Define a sample challenge
	challenge := entity.Challenge{
		// Fill in challenge attributes as necessary
	}

	// Define scenarios
	t.Run("successful challenge creation and scoring", func(t *testing.T) {
		// Mock successful challenge generation
		mockMinter.On("GenerateChallenge").Return(challenge, nil).Once()

		// Mock successful scoring
		mockRecordable.On("Score", ctx, mock.AnythingOfType("uuid.UUID"), challenge).Return(nil).Once()

		// Execute
		returnedRemark, returnedChallenge, err := challenger.MakeChallenge(ctx)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, challenge, returnedChallenge)
		assert.NotEqual(t, uuid.Nil, returnedRemark)
		mockMinter.AssertExpectations(t)
		mockRecordable.AssertExpectations(t)
	})

	t.Run("failed to generate challenge", func(t *testing.T) {
		// Mock challenge generation failure
		mockMinter.On("GenerateChallenge").Return(entity.Challenge{}, errors.New("generation error")).Once()

		// Execute
		returnedRemark, returnedChallenge, err := challenger.MakeChallenge(ctx)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, entity.Challenge{}, returnedChallenge)
		assert.Equal(t, uuid.Nil, returnedRemark)
		assert.Contains(t, err.Error(), "failed to generate challenge")
		mockMinter.AssertExpectations(t)
	})

	t.Run("failed to score challenge", func(t *testing.T) {
		// Mock successful challenge generation
		mockMinter.On("GenerateChallenge").Return(challenge, nil).Once()

		// Mock scoring failure
		mockRecordable.On("Score", ctx, mock.AnythingOfType("uuid.UUID"), challenge).Return(errors.New("scoring error")).Once()

		// Execute
		returnedRemark, returnedChallenge, err := challenger.MakeChallenge(ctx)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, entity.Challenge{}, returnedChallenge)
		assert.Equal(t, uuid.Nil, returnedRemark)
		assert.Contains(t, err.Error(), "failed to score challenge")
		mockMinter.AssertExpectations(t)
		mockRecordable.AssertExpectations(t)
	})
}
