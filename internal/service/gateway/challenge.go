package gateway

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetChallenge: Generates a random challenge for the PoW
func (s *Server) GetChallenge(c echo.Context) error {
	challenge, err := s.challenger.MakeChallenge()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Challenge emit failed"})
	}

	response := map[string]string{
		"challenge": challenge,
	}

	return c.JSON(http.StatusOK, response)
}
