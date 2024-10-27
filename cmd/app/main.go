package main

import (
	"fmt"
	"github.com/littlebugger/pow-wow/deps/api"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	echo_mid "github.com/labstack/echo/v4/middleware"
	api_mid "github.com/oapi-codegen/echo-middleware"
)

type ServerImpl struct{}

var quotes = []string{
	"The only limit to our realization of tomorrow is our doubts of today.",
	"Do not wait to strike till the iron is hot; but make it hot by striking.",
	"The greatest glory in living lies not in never falling, but in rising every time we fall.",
}

// GetChallenge: Generates a random challenge for the PoW
func (s *ServerImpl) GetChallenge(c echo.Context) error {
	rand.Seed(time.Now().UnixNano())
	challenge := fmt.Sprintf("%x", rand.Intn(100000))
	response := map[string]string{
		"challenge": challenge,
	}
	return c.JSON(http.StatusOK, response)
}

// VerifySolution: Verifies the PoW solution (nonce) provided by the client
func (s *ServerImpl) VerifySolution(c echo.Context) error {
	var req struct {
		Challenge string `json:"challenge"`
		Nonce     string `json:"nonce"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Here, add your logic to verify the PoW solution (nonce + challenge)
	if isValidPoW(req.Challenge, req.Nonce) {
		quote := quotes[rand.Intn(len(quotes))]
		response := map[string]string{
			"quote": quote,
		}
		return c.JSON(http.StatusOK, response)
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid PoW solution"})
	}
}

// Dummy function to simulate PoW validation
func isValidPoW(challenge, nonce string) bool {
	// In real-world scenario, you'd implement actual PoW verification logic here
	// For now, just return true for simplicity
	return true
}

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

	// Create the api implementation
	server := &ServerImpl{}

	// Register the routes from the generated handler
	api.RegisterHandlers(e, server)

	// Start the api
	log.Println("Starting api on :8080")
	log.Fatal(e.Start(":8080"))
}
