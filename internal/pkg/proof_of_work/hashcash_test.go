package proof_of_work

import (
	"testing"

	"github.com/littlebugger/pow-wow/internal/service/entity"
	"github.com/stretchr/testify/assert"
)

func TestHashcash_GenerateChallenge(t *testing.T) {
	hcv := &Hashcash{}

	t.Run("Generate a valid challenge", func(t *testing.T) {
		challenge, err := hcv.GenerateChallenge()

		// Assertions
		assert.NoError(t, err)
		assert.NotEmpty(t, challenge.Task, "Task should not be empty")
		assert.GreaterOrEqual(t, challenge.Difficulty, 1, "Difficulty should be at least 1")
	})
}

func Test_verifyPoW(t *testing.T) {
	tests := []struct {
		name      string
		challenge entity.Challenge
		nonce     string
		want      bool
	}{
		{
			name: "Correct nonce with low difficulty",
			challenge: entity.Challenge{
				Task:       "test-challenge",
				Difficulty: 1,
			},
			nonce: "00001",
			want:  true,
		},
		{
			name: "Incorrect nonce",
			challenge: entity.Challenge{
				Task:       "test-challenge",
				Difficulty: 3,
			},
			nonce: "123456",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := verifyPoW(tt.challenge, tt.nonce)
			assert.Equal(t, tt.want, got, "verifyPoW() result should match expected value")
		})
	}
}

func TestSolveChallenge(t *testing.T) {
	tests := []struct {
		name      string
		challenge entity.Challenge
		nonce     string
		wantErr   bool
	}{
		{
			name: "Successfully solve challenge with low difficulty",
			challenge: entity.Challenge{
				Task:       "test-challenge",
				Difficulty: 3,
			},
			nonce:   "388",
			wantErr: false,
		},
		{
			name: "Solve challenge with higher difficulty",
			challenge: entity.Challenge{
				Task:       "difficult-task",
				Difficulty: 5,
			},
			nonce:   "747542",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nonce, err := SolveChallenge(tt.challenge)

			if (err != nil) != tt.wantErr {
				t.Errorf("SolveChallenge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify that the nonce successfully passes the verifyPoW check
			if !tt.wantErr {
				assert.NotEmpty(t, nonce, "Nonce should not be empty")
				assert.Equal(t, tt.nonce, nonce, "Nonce should match")
				assert.True(t, verifyPoW(tt.challenge, nonce), "Solved nonce should pass the verifyPoW() check")
			}
		})
	}
}
