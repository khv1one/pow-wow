package proof_of_work

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Hashcash struct{}

func NewHashcash() *Hashcash {
	return &Hashcash{}
}

// Dummy function to simulate PoW validation
func (hcv *Hashcash) VerifySolution(challenge, nonce string) (bool, error) {
	// In real-world scenario, you'd implement actual PoW verification logic here
	// For now, just return true for simplicity
	return true, nil
}

func (hcv *Hashcash) GenerateChallenge() (string, error) {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%x", rand.Intn(100000)), nil
}

// Simple PoW: find a nonce such that hash(challenge + nonce) has leading zeros
func verifyPoW(challenge, nonce string, difficulty int) bool {
	data := challenge + nonce
	hash := sha256.Sum256([]byte(data))
	hashStr := hex.EncodeToString(hash[:])
	// Check if the hash has the required number of leading zeros
	return strings.HasPrefix(hashStr, strings.Repeat("0", difficulty))
}
