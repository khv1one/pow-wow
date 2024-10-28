package main

import (
	"github.com/littlebugger/pow-wow/internal/service/repository/redis"
	"log"

	api "github.com/littlebugger/pow-wow/deps/api"
	"github.com/littlebugger/pow-wow/internal/pkg/proof_of_work"
	"github.com/littlebugger/pow-wow/internal/service"
	"github.com/littlebugger/pow-wow/internal/service/gateway"
	"github.com/littlebugger/pow-wow/internal/service/repository/potgresql"
	"github.com/littlebugger/pow-wow/internal/service/usecase"
)

func main() {
	app := service.NewApp()
	server := createServer(app)
	// Register the routes from the generated handler
	api.RegisterHandlers(app.Echo, server)

	// Start the api
	log.Println("Starting api on :8080")
	log.Fatal(app.Echo.Start(":8080"))
}

// createServer encapsulates all business end of application init.
func createServer(app *service.App) *gateway.Server {
	// Method for create and verify POW challenges.
	hashcash := proof_of_work.NewHashcash()

	// Storage for wisdom quotes.
	repo := potgresql.NewPGWisdomRepository(app.Gorm)
	rdb := redis.NewJournal(app.RDB)
	journal := usecase.NewJournal(rdb)

	// Use cases for emit challenges and check solutions.
	challenger := usecase.NewChallenger(hashcash, journal)
	overseer := usecase.NewOverseer(hashcash, repo, journal)

	// Create the api implementation
	return gateway.NewServer(challenger, overseer)
}
