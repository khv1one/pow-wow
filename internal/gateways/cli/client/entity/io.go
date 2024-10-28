package entity

// ChallengeResponse represents the structure of the JSON response from /challenge
type ChallengeResponse struct {
	Challenge  string `json:"challenge"`
	Difficulty string `json:"difficulty"`
}

type VerifyRequest struct {
	Challenge string `json:"challenge"`
	Nonce     string `json:"nonce"`
}

type VerifyResponse struct {
	Quote string `json:"quote"`
}
