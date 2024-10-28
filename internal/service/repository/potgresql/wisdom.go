package potgresql

import "math/rand"

var quotes = []string{
	"The only limit to our realization of tomorrow is our doubts of today.",
	"Do not wait to strike till the iron is hot; but make it hot by striking.",
	"The greatest glory in living lies not in never falling, but in rising every time we fall.",
}

type PGWisdomRepository struct{}

func NewPGWisdomRepository() *PGWisdomRepository {
	return &PGWisdomRepository{}
}

func (r PGWisdomRepository) ExpandWisdom() (string, error) {
	return quotes[rand.Intn(len(quotes))], nil
}
