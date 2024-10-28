
# Word of Wisdom TCP Server with Proof of Work (PoW)

This project implements a **Word of Wisdom TCP server** with a Hashcash Proof of Work (PoW) mechanism to protect against DDoS attacks. Clients must solve a PoW challenge to receive a random wisdom quote from the server.

## Table of Contents
- [Introduction](#introduction)
- [Features](#features)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [API Specification](#api-specification)
- [PoW Solver Example](#pow-solver-example)
- [Running Tests](#running-tests)

## Introduction
The **Word of Wisdom TCP Server** provides inspirational quotes to clients who successfully solve a Hashcash-based PoW challenge. This project leverages **Go**, **Redis**, and **PostgreSQL** to store challenges, manage sessions, and serve quotes efficiently.

## Features
- **TCP Server with HTTP Interface**: Provides a RESTful API to receive and verify PoW challenges.
- **Hashcash-based PoW**: Uses Hashcash PoW to prevent abuse and protect against DDoS attacks.
- **Database Integration**: Stores quotes in a PostgreSQL database and manages PoW sessions in Redis.
- **Dockerized Setup**: Provides a `docker-compose` setup to manage all services.

## Architecture
The project consists of:
- **Server**: A TCP server exposed through a RESTful API using the Echo framework.
- **Database**: PostgreSQL database to store inspirational quotes.
- **Cache**: Redis to handle challenge sessions and temporary storage.
- **Client**: A command-line tool that connects to the server, retrieves a challenge, solves it, and verifies the solution.

## Prerequisites
- [Go 1.22.6+](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/littlebugger/pow-wow.git
   cd pow-wow
   ```

2. **Build and run services using Docker Compose**:
   ```bash
   ./scripts/docker-compose.sh
   ```

3. **Run the database migrations if nessesary**:
   ```bash
   goose -dir ./db/migrations postgres "host=localhost user=your_user password=your_password dbname=your_db port=5432 sslmode=disable" up
   ```

4. **Build and run storages using Docker Compose**:
   ```bash
   ./scripts/docker-compose.sh --storages-only
   ```

5. **Delete all containers and images**:
   ```bash
   ./scripts/docker-compose.sh down
   ```

## Usage

1. **Start the server** (if not using Docker Compose):
   ```bash
   go run cmd/app/main.go
   ```

2. **Run the client** to interact with the server:
   ```bash
   go run cmd/cli/main.go
   ```
   
3. **Regenerate openapi**:
   ```bash
   ./scripts/openapi-gen.sh
   ```

4. **Generate mocks**:
   ```bash
   ./scripts/mockery.sh
   ```
   
## API Specification
The server provides two main endpoints:

### 1. Get Challenge (`/challenge`)
- **Method**: `GET`
- **Description**: Fetches a PoW challenge with a difficulty level.
- **Response**:
    - `200 OK`
        - **Headers**:
            - `X-Remark`: A unique identifier for the challenge.
        - **Body**:
          ```json
          {
            "difficulty": 5,
            "challenge": "abcdef"
          }
          ```

### 2. Verify Solution (`/verify`)
- **Method**: `POST`
- **Description**: Verifies the PoW solution provided by the client.
- **Headers**:
    - `X-Remark`: A unique identifier for the challenge.
- **Request Body**:
  ```json
  {
    "nonce": "654321"
  }
  ```
- **Response**:
    - `200 OK`
        - **Body**:
          ```json
          {
            "quote": "Success is not final, failure is not fatal: it is the courage to continue that counts."
          }
          ```

## PoW Solver Example

The client fetches a challenge, solves it using a simple brute-force Hashcash algorithm, and posts the solution to the server for verification.

Client Logic (cmd/client/main.go)

The client logic involves:

	1.	Fetching a challenge from the server.
	2.	Solving the challenge using a hashcash package.
	3.	Posting the solution back to the server to receive a wisdom quote.

## Running Test

You can run tests in project:
   ```bash
   go test -json ./... 
   ```