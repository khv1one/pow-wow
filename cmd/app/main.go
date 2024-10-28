package main

import (
	"github.com/labstack/echo/v4"
	echo_mid "github.com/labstack/echo/v4/middleware"
	"github.com/littlebugger/pow-wow/deps/api"
	"github.com/littlebugger/pow-wow/internal/pkg/proof_of_work"
	"github.com/littlebugger/pow-wow/internal/service/gateway"
	"github.com/littlebugger/pow-wow/internal/service/repository/potgresql"
	"github.com/littlebugger/pow-wow/internal/service/usecase"
	api_mid "github.com/oapi-codegen/echo-middleware"
	"log"
)

func main() {
	// Load the OpenAPI specification
	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatalf("Error loading OpenAPI spec: %v", err)
	}

	// Create Echo instance
	e := echo.New()

	// Add Echo middleware (optional for logging, recovery, etc.)
	e.Use(echo_mid.Logger())
	e.Use(echo_mid.Recover())

	// Add OpenAPI request validation middleware using oapi-codegen's middleware
	e.Use(api_mid.OapiRequestValidator(swagger))

	server := createServer()
	// Register the routes from the generated handler
	api.RegisterHandlers(e, server)

	// Start the api
	log.Println("Starting api on :8080")
	log.Fatal(e.Start(":8080"))
}

// createServer encapsulates all business end of application init.
func createServer() *gateway.Server {
	// Method for create and verify POW challenges.
	hashcash := proof_of_work.NewHashcash()

	// Storage for wisdom quotes.
	repo := potgresql.NewPGWisdomRepository()

	// Use cases for emit challenges and check solutions.
	challenger := usecase.NewChallenger(hashcash)
	overseer := usecase.NewOverseer(hashcash, repo)

	// Create the api implementation
	return gateway.NewServer(challenger, overseer)
}
