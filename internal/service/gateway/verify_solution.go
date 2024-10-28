package gateway

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// VerifySolution: Verifies the PoW solution (nonce) provided by the client
func (s *Server) VerifySolution(c echo.Context) error {
	var req struct {
		Challenge string `json:"challenge"`
		Nonce     string `json:"nonce"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	quote, err := s.supervisor.Oversee(req.Challenge, req.Nonce)
	if err == nil {
		response := map[string]string{
			"quote": quote,
		}

		return c.JSON(http.StatusOK, response)
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid PoW solution"})
	}
}
