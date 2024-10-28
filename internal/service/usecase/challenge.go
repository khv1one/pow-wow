package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/littlebugger/pow-wow/internal/service/entity"
)

type Minter interface {
	GenerateChallenge() (entity.Challenge, error)
}

type Recordable interface {
	Score(ctx context.Context, key uuid.UUID, challenge entity.Challenge) error
}

type Challenger struct {
	minter     Minter
	recordable Recordable
}

func NewChallenger(m Minter, r Recordable) *Challenger {
	return &Challenger{minter: m, recordable: r}
}

func (rc *Challenger) MakeChallenge(ctx context.Context) (uuid.UUID, entity.Challenge, error) {
	challenge, err := rc.minter.GenerateChallenge()
	if err != nil {
		return uuid.Nil, entity.Challenge{}, fmt.Errorf("failed to generate challenge: %w", err)
	}

	remark := uuid.New()

	err = rc.recordable.Score(ctx, remark, challenge)
	if err != nil {
		return uuid.Nil, entity.Challenge{}, fmt.Errorf("failed to score challenge: %w", err)
	}

	return remark, challenge, nil
}
