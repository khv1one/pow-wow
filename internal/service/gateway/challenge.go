package gateway

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// GetChallenge: Generates a random challenge for the PoW
func (s *Server) GetChallenge(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*1)
	defer cancel()

	remark, challenge, err := s.Challenger.MakeChallenge(ctx)
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
