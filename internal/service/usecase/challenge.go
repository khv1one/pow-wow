package usecase

type Minter interface {
	GenerateChallenge() (string, error)
}

type Challenger struct {
	minter Minter
}

func NewChallenger(m Minter) *Challenger {
	return &Challenger{minter: m}
}

func (rc *Challenger) MakeChallenge() (string, error) {
	return rc.minter.GenerateChallenge()
}
