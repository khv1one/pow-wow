package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/littlebugger/pow-wow/internal/service/entity"
)

type Verifier interface {
	VerifySolution(challenge entity.Challenge, solution string) (bool, error)
}

type Registry interface {
	Match(ctx context.Context, key uuid.UUID) (entity.Challenge, error)
	MarkSolved(ctx context.Context, key uuid.UUID) error
}

type WisdomBearer interface {
	ExpandWisdom(ctx context.Context) (string, error)
}

type Overseer struct {
	verifier  Verifier
	sapient   WisdomBearer
	validator Registry
}

func NewOverseer(v Verifier, s WisdomBearer, val Registry) *Overseer {
	return &Overseer{verifier: v, sapient: s, validator: val}
}

func (sv *Overseer) Oversee(ctx context.Context, remark uuid.UUID, solution string) (string, error) {
	challenge, err := sv.validator.Match(ctx, remark)
	if err != nil {
		return "", fmt.Errorf("failed to verify solution: %w", err)
	}

	ok, err := sv.verifier.VerifySolution(challenge, solution)
	if err != nil {
		return "", fmt.Errorf("failed to verify solution: %w", err)
	}
	if ok {
		err = sv.validator.MarkSolved(ctx, remark)
		if err != nil {
			return "", fmt.Errorf("failed to mark solved solution: %w", err)
		}

		return sv.sapient.ExpandWisdom(ctx)
	}

	return "", fmt.Errorf("icorrect solution presented")
}
