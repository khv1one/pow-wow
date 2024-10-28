package gateway

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetChallenge: Generates a random challenge for the PoW
func (s *Server) GetChallenge(c echo.Context) error {
	ctx := c.Request().Context()

	remark, challenge, err := s.challenger.MakeChallenge(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Challenge emit failed"})
	}

	// Set the remark as a custom header
	c.Response().Header().Set("X-Remark", remark.String())

	// TODO: make hanler error and handler return entity
	response := map[string]string{
		"challenge":  challenge.Task,
		"difficulty": strconv.Itoa(challenge.Difficulty),
	}

	return c.JSON(http.StatusOK, response)
}
