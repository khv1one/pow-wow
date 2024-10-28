package potgresql

import (
	"github.com/littlebugger/pow-wow/internal/service/entity"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type PGWisdomRepository struct {
	repo *gorm.DB
}

func NewPGWisdomRepository(db *gorm.DB) *PGWisdomRepository {
	return &PGWisdomRepository{repo: db}
}

func (r PGWisdomRepository) ExpandWisdom() (string, error) {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Get the total number of quotes in the table
	var count int
	err := r.repo.Raw("SELECT COUNT(*) FROM quotes").Scan(&count).Error
	if err != nil {
		return "", err
	}

	// If there are no quotes, return an error or fallback quote
	if count == 0 {
		return "", gorm.ErrRecordNotFound
	}

	// Generate a random offset
	randomOffset := rand.Intn(count)

	// Retrieve a random quote using a random offset in a raw SQL query
	var selectedQuote entity.Quote
	err = r.repo.Raw("SELECT id, quote, author FROM quotes OFFSET ? LIMIT 1", randomOffset).Scan(&selectedQuote).Error
	if err != nil {
		return "", err
	}

	// Format the quote with the author if available
	if selectedQuote.Author != nil {
		return selectedQuote.Quote + " — " + *selectedQuote.Author, nil
	}
	return selectedQuote.Quote, nil
}
