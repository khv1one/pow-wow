package gateway

import (
	"context"
	"github.com/google/uuid"
	"github.com/littlebugger/pow-wow/internal/service/entity"
)

type Challenger interface {
	MakeChallenge(ctx context.Context) (uuid.UUID, entity.Challenge, error)
}

type Supervisor interface {
	Oversee(ctx context.Context, remark uuid.UUID, solution string) (string, error)
}

type Server struct {
	challenger Challenger
	supervisor Supervisor
}

func NewServer(c Challenger, s Supervisor) *Server {
	return &Server{challenger: c, supervisor: s}
}
