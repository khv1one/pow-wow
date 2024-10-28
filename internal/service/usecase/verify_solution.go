package usecase

import "fmt"

type Verifier interface {
	VerifySolution(challenge, nonce string) (bool, error)
}

type WisdomBearer interface {
	ExpandWisdom() (string, error)
}

type Overseer struct {
	verifier Verifier
	sapient  WisdomBearer
}

func NewOverseer(v Verifier, s WisdomBearer) *Overseer {
	return &Overseer{verifier: v, sapient: s}
}

func (sv *Overseer) Oversee(challenge, nonce string) (string, error) {
	ok, err := sv.verifier.VerifySolution(challenge, nonce)
	if err != nil {
		return "", fmt.Errorf("failed to verify solution: %w", err)
	}
	if ok {
		return sv.sapient.ExpandWisdom()
	}

	return "", fmt.Errorf("icorrect solution presented")
}
