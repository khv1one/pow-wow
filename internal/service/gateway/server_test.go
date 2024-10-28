package gateway_test

import (
	"bytes"
	"errors"
	pow "github.com/littlebugger/pow-wow/deps/api"
	"strconv"
	"testing"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/littlebugger/pow-wow/internal/service/entity"
	"github.com/littlebugger/pow-wow/internal/service/gateway"
	"github.com/littlebugger/pow-wow/internal/service/gateway/mocks"
)

func TestServer_GetChallenge(t *testing.T) {
	// Setup
	e := echo.New()
	mockChallenger := new(mocks.MockChallenger)
	server := &gateway.Server{
		Challenger: mockChallenger,
	}

	// Define a sample challenge
	remark := uuid.New()
	challenge := entity.Challenge{
		Task:       "sample-task",
		Difficulty: 3,
	}

	// Define scenarios
	t.Run("successful challenge generation", func(t *testing.T) {
		// Mock successful challenge generation
		mockChallenger.On("MakeChallenge", mock.Anything).Return(remark, challenge, nil).Once()

		// Create a request and response recorder
		req := httptest.NewRequest(http.MethodGet, "/challenge", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Execute the handler
		if assert.NoError(t, server.GetChallenge(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, remark.String(), rec.Header().Get("X-Remark"))

			// Parse the response body
			var response map[string]string
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, challenge.Task, response["challenge"])
			assert.Equal(t, strconv.Itoa(challenge.Difficulty), response["difficulty"])
		}

		mockChallenger.AssertExpectations(t)
	})

	t.Run("challenge generation fails", func(t *testing.T) {
		// Mock failure in challenge generation
		mockChallenger.On("MakeChallenge", mock.Anything).Return(uuid.Nil, entity.Challenge{}, errors.New("generation failed")).Once()

		// Create a request and response recorder
		req := httptest.NewRequest(http.MethodGet, "/challenge", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Execute the handler
		if assert.NoError(t, server.GetChallenge(c)) {
			assert.Equal(t, http.StatusInternalServerError, rec.Code)

			// Parse the response body
			var response map[string]string
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Challenge emit failed", response["error"])
		}

		mockChallenger.AssertExpectations(t)
	})
}

func TestServer_VerifySolution(t *testing.T) {
	// Setup
	e := echo.New()
	mockSupervisor := new(mocks.MockSupervisor)
	server := &gateway.Server{
		Supervisor: mockSupervisor,
	}

	// Define a sample remark and nonce
	remark := uuid.New()
	nonce := "correct-nonce"

	// Define scenarios
	t.Run("successful solution verification", func(t *testing.T) {
		// Mock successful verification
		mockSupervisor.On("Oversee", mock.Anything, remark, nonce).Return("A wise quote", nil).Once()

		// Create a request and response recorder
		reqBody := `{"nonce":"correct-nonce"}`
		req := httptest.NewRequest(http.MethodPost, "/verify", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Remark", remark.String())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Execute the handler
		if assert.NoError(t, server.VerifySolution(c, pow.VerifySolutionParams{XRemark: remark})) {
			assert.Equal(t, http.StatusOK, rec.Code)

			// Parse the response body
			var response map[string]string
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "A wise quote", response["quote"])
		}

		mockSupervisor.AssertExpectations(t)
	})

	t.Run("invalid PoW solution", func(t *testing.T) {
		// Mock invalid solution
		mockSupervisor.On("Oversee", mock.Anything, remark, "incorrect-nonce").Return("", errors.New("invalid solution")).Once()

		// Create a request and response recorder
		reqBody := `{"nonce":"incorrect-nonce"}`
		req := httptest.NewRequest(http.MethodPost, "/verify", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Remark", remark.String())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Execute the handler
		if assert.NoError(t, server.VerifySolution(c, pow.VerifySolutionParams{XRemark: remark})) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			// Parse the response body
			var response map[string]string
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Invalid PoW solution", response["error"])
		}

		mockSupervisor.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		// Create a request with an invalid JSON body
		reqBody := `{"nonce":}` // Malformed JSON
		req := httptest.NewRequest(http.MethodPost, "/verify", bytes.NewBuffer([]byte(reqBody)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Remark", remark.String())
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Execute the handler
		if assert.NoError(t, server.VerifySolution(c, pow.VerifySolutionParams{XRemark: remark})) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			// Parse the response body
			var response map[string]string
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "Invalid request", response["error"])
		}
	})
}
