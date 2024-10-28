package usecase

import (
	"fmt"
	"log"

	"github.com/littlebugger/pow-wow/internal/gateways/cli/client"
	"github.com/littlebugger/pow-wow/internal/pkg/proof_of_work"
	"github.com/littlebugger/pow-wow/internal/service/entity"
)

func Handle() {
	// Step 1: Get the challenge
	remark, task, difficulty, err := client.GetChallenge()
	if err != nil {
		log.Fatalf("Failed to get challenge: %v", err)
	}

	challenge := entity.Challenge{
		Task:       task,
		Difficulty: difficulty,
	}

	// Step 2: Solve the PoW challenge
	fmt.Printf("Solving challenge: %s with difficulty: %d\n", challenge.Task, challenge.Difficulty)
	nonce, err := proof_of_work.SolveChallenge(challenge)
	if err != nil {
		log.Fatalf("Failed to solve challenge: %v", err)
	}
	fmt.Printf("Solution found! Nonce: %s\n", nonce)

	// Step 3: Post the solution to /verify
	quote, err := client.VerifyChallenge(challenge.Task, nonce, remark)
	if err != nil {
		log.Fatalf("Failed to verify challenge: %v", err)
	}
	fmt.Printf("Received quote: %s\n", quote)
}
