package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/littlebugger/pow-wow/internal/service/entity"
)

type Trackable interface {
	Save(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (*string, error)
	Delete(ctx context.Context, key string) error
}

type Journal struct {
	storage Trackable
}

func NewJournal(s Trackable) *Journal {
	return &Journal{storage: s}
}

func (j *Journal) Score(ctx context.Context, key uuid.UUID, challenge entity.Challenge) error {
	value, err := json.Marshal(challenge)
	if err != nil {
		return fmt.Errorf("failed to transmute challenge: %w", err)
	}

	return j.storage.Save(ctx, key.String(), string(value))
}

func (j *Journal) Match(ctx context.Context, key uuid.UUID) (entity.Challenge, error) {
	challenge := entity.Challenge{}
	stored, err := j.storage.Get(ctx, key.String())
	if err != nil {
		return challenge, fmt.Errorf("something wrong with challenge journal")
	}
	if stored == nil {
		return challenge, fmt.Errorf("no such challenge")
	}

	err = json.Unmarshal([]byte(*stored), &challenge)
	if err != nil {
		return challenge, fmt.Errorf("challenge was corrupted")
	}

	return challenge, nil
}

func (j *Journal) MarkSolved(ctx context.Context, key uuid.UUID) error {
	err := j.storage.Delete(ctx, key.String())
	if err != nil {
		return fmt.Errorf("failed to mark challenge solved: %w", err)
	}

	return nil
}
