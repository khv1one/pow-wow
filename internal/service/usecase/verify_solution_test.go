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
)

func TestOverseer_Oversee(t *testing.T) {
	// Create context and mock instances
	ctx := context.Background()
	mockVerifier := mocks.NewMockVerifier(t)
	mockRegistry := mocks.NewMockRegistry(t)
	mockWisdomBearer := mocks.NewMockWisdomBearer(t)

	// Create the overseer with mocks
	overseer := usecase.NewOverseer(mockVerifier, mockWisdomBearer, mockRegistry)

	// Define a sample challenge and remark UUID
	sampleChallenge := entity.Challenge{
		// Fill in challenge attributes as necessary
	}
	remark := uuid.New()
	solution := "correct_solution"

	// Define successful scenario behavior
	t.Run("successful verification and wisdom retrieval", func(t *testing.T) {
		// Set up mock expectations
		mockRegistry.On("Match", ctx, remark).Return(sampleChallenge, nil).Once()
		mockVerifier.On("VerifySolution", sampleChallenge, solution).Return(true, nil).Once()
		mockRegistry.On("MarkSolved", ctx, remark).Return(nil).Once()
		mockWisdomBearer.On("ExpandWisdom", ctx).Return("A wise quote", nil).Once()

		// Execute Oversee
		quote, err := overseer.Oversee(ctx, remark, solution)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, "A wise quote", quote)

		// Assert that all expectations were met
		mockVerifier.AssertExpectations(t)
		mockRegistry.AssertExpectations(t)
		mockWisdomBearer.AssertExpectations(t)
	})

	// Define failure scenario when challenge does not exist
	t.Run("challenge not found", func(t *testing.T) {
		// Set up mock expectations
		mockRegistry.On("Match", ctx, remark).Return(entity.Challenge{}, errors.New("not found")).Once()

		// Execute Oversee
		quote, err := overseer.Oversee(ctx, remark, solution)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, "", quote)
		assert.Contains(t, err.Error(), "failed to verify solution")

		// Assert that all expectations were met
		mockRegistry.AssertExpectations(t)
	})

	// Define failure scenario when solution is incorrect
	t.Run("incorrect solution", func(t *testing.T) {
		// Set up mock expectations
		mockRegistry.On("Match", ctx, remark).Return(sampleChallenge, nil).Once()
		mockVerifier.On("VerifySolution", sampleChallenge, solution).Return(false, nil).Once()

		// Execute Oversee
		quote, err := overseer.Oversee(ctx, remark, solution)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, "", quote)
		assert.Contains(t, err.Error(), "icorrect solution presented")

		// Assert that all expectations were met
		mockVerifier.AssertExpectations(t)
		mockRegistry.AssertExpectations(t)
	})

	// Define failure scenario when marking as solved fails
	t.Run("marking as solved fails", func(t *testing.T) {
		// Set up mock expectations
		mockRegistry.On("Match", ctx, remark).Return(sampleChallenge, nil).Once()
		mockVerifier.On("VerifySolution", sampleChallenge, solution).Return(true, nil).Once()
		mockRegistry.On("MarkSolved", ctx, remark).Return(errors.New("marking failed")).Once()

		// Execute Oversee
		quote, err := overseer.Oversee(ctx, remark, solution)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, "", quote)
		assert.Contains(t, err.Error(), "failed to mark solved solution")

		// Assert that all expectations were met
		mockVerifier.AssertExpectations(t)
		mockRegistry.AssertExpectations(t)
	})
}
