package proof_of_work

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"crypto/sha256"
	"encoding/hex"

	"github.com/littlebugger/pow-wow/internal/service/entity"
)

// TODO: make difficulty an argument
const difficulty = 6

type Hashcash struct{}

func NewHashcash() *Hashcash {
	return &Hashcash{}
}

// Dummy function to simulate PoW validation
func (hcv *Hashcash) VerifySolution(challenge entity.Challenge, solution string) (bool, error) {
	return verifyPoW(challenge, solution), nil
}

func (hcv *Hashcash) GenerateChallenge() (entity.Challenge, error) {
	rand.Seed(time.Now().UnixNano())
	return entity.Challenge{
		Task:       fmt.Sprintf("%x", rand.Intn(100000)),
		Difficulty: difficulty,
	}, nil
}

// Simple PoW: find a nonce such that hash(challenge + nonce) has leading zeros
func verifyPoW(challenge entity.Challenge, nonce string) bool {
	data := challenge.Task + nonce
	hash := sha256.Sum256([]byte(data))
	hashStr := hex.EncodeToString(hash[:])
	// Check if the hash has the required number of leading zeros
	return strings.HasPrefix(hashStr, strings.Repeat("0", challenge.Difficulty))
}

// SolveChallenge attempts to find a nonce that results in a hash with the required number of leading zeros
func SolveChallenge(challenge entity.Challenge) string {
	leadingZeros := strings.Repeat("0", challenge.Difficulty)
	var nonce int64

	for {
		// Generate a potential solution by concatenating the challenge and nonce
		input := fmt.Sprintf("%s%d", challenge.Task, nonce)
		hash := sha256.Sum256([]byte(input))

		// Convert the hash to a hex string and check if it meets the difficulty requirement
		hashStr := hex.EncodeToString(hash[:])
		if strings.HasPrefix(hashStr, leadingZeros) {
			return fmt.Sprintf("%d", nonce)
		}
		nonce++
	}
}

// SolveChallengeWithNonce same as SolveChallenge but with different interface for better scalability.
func SolveChallengeWithNonce(challenge entity.Challenge, nonce int64) error {
	leadingZeros := strings.Repeat("0", challenge.Difficulty)

	// Generate a potential solution by concatenating the challenge and nonce
	input := fmt.Sprintf("%s%d", challenge.Task, nonce)
	hash := sha256.Sum256([]byte(input))

	hashStr := hex.EncodeToString(hash[:])
	if strings.HasPrefix(hashStr, leadingZeros) {
		return nil
	}

	return fmt.Errorf("no solution found")
}
