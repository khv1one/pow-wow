package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/littlebugger/pow-wow/internal/gateways/cli/client/entity"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

const serverURL = "http://localhost:8080"

func GetChallenge() (string, string, int, error) {
	resp, err := http.Get(serverURL + "/challenge")
	if err != nil {
		return "", "", 0, err
	}
	defer resp.Body.Close()

	// Read the X-Remark header
	remark := resp.Header.Get("X-Remark")
	if remark == "" {
		return "", "", 0, fmt.Errorf("missing X-Remark header in response")
	}

	// Parse the JSON response
	var challengeResp entity.ChallengeResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", 0, err
	}
	if err := json.Unmarshal(body, &challengeResp); err != nil {
		return "", "", 0, err
	}

	diff, err := strconv.Atoi(challengeResp.Difficulty)
	if err != nil {
		return "", "", 0, err
	}

	return remark, challengeResp.Challenge, diff, nil
}

func VerifyChallenge(challenge, nonce, remark string) (string, error) {
	// Prepare the JSON request body
	requestBody, err := json.Marshal(entity.VerifyRequest{
		Challenge: challenge,
		Nonce:     nonce,
	})
	if err != nil {
		return "", err
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", serverURL+"/verify", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Remark", remark)

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check if the response indicates an error
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("verification failed: %s", string(body))
	}

	// Parse the JSON response
	var verifyResp entity.VerifyResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(body, &verifyResp); err != nil {
		return "", err
	}

	return verifyResp.Quote, nil
}
