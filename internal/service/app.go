package service

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	echo_mid "github.com/labstack/echo/v4/middleware"
	api_mid "github.com/oapi-codegen/echo-middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	api "github.com/littlebugger/pow-wow/deps/api"
)

type App struct {
	Echo *echo.Echo
	RDB  *redis.Client
	Gorm *gorm.DB
}

func NewApp() *App {
	a := App{}
	a.init()
	return &a
}

func (a *App) init() {
	a.redisConn()
	a.postgresConn()
	a.serverSetup()
}

// redisConn init for rdb.
func (a *App) redisConn() {
	// Redis setup
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Test Redis connection
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")

	a.RDB = rdb
}

func (a *App) postgresConn() {
	// PostgreSQL setup
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASSWORD")
	postgresDB := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		postgresHost, postgresPort, postgresUser, postgresPassword, postgresDB)

	psql, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	a.Gorm = psql
}

func (a *App) serverSetup() {
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

	log.Println("Server init done")

	a.Echo = e
}
