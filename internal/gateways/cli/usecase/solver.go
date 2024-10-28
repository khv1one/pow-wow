package usecase

import (
	"context"
	"fmt"
	"github.com/littlebugger/pow-wow/internal/gateways/cli/client"
	"github.com/littlebugger/pow-wow/internal/pkg/proof_of_work"
	"github.com/littlebugger/pow-wow/internal/service/entity"
	"log"
	"sync"
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

	nonce, err := SolveChallengeConcurrently(context.Background(), challenge)
	if err != nil {
		log.Fatalf("Failed to solve challenge, terminated: %v", err)
	}

	fmt.Printf("Solution found! Nonce: %s\n", nonce)

	// Step 3: Post the solution to /verify
	quote, err := client.VerifyChallenge(challenge.Task, fmt.Sprintf("%s", nonce), remark)
	if err != nil {
		log.Fatalf("Failed to verify challenge: %v", err)
	}
	fmt.Printf("Received quote: %s\n", quote)
}

func SolveChallengeConcurrently(ctx context.Context, challenge entity.Challenge) (string, error) {
	var wg sync.WaitGroup
	resultChan := make(chan string, 1) // Channel to capture the first found solution
	quitChan := make(chan struct{})    // Channel to signal workers to quit

	var nonce int64
	// Launch workers
	go func() {
		for {
			select {
			case <-quitChan:
				break
			default:
				wg.Add(1)
				go func(solution int64) {
					defer wg.Done()
					err := proof_of_work.SolveChallengeWithNonce(challenge, solution)
					if err == nil {
						resultChan <- fmt.Sprintf("%d", solution)
					}
				}(nonce)

				nonce++
			}
		}
	}()

	var solution string

	select {
	case solution = <-resultChan:
		close(quitChan) // Signal all workers to stop
	case <-ctx.Done():
		wg.Wait() // Ensure all workers are finished before exiting
		return "", ctx.Err()
	}

	// Wait for all workers to finish
	wg.Wait()
	return solution, nil
}
