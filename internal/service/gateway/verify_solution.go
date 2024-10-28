package gateway

import (
	"net/http"

	"github.com/labstack/echo/v4"

	pow "github.com/littlebugger/pow-wow/deps/api"
)

// VerifySolution: Verifies the PoW solution (nonce) provided by the client
func (s *Server) VerifySolution(c echo.Context, params pow.VerifySolutionParams) error {
	var req struct {
		Nonce string `json:"nonce"`
	}

	ctx := c.Request().Context()

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	quote, err := s.supervisor.Oversee(ctx, params.XRemark, req.Nonce)
	if err == nil {
		response := map[string]string{
			"quote": quote,
		}

		return c.JSON(http.StatusOK, response)
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid PoW solution"})
	}
}
