package randomrepo

import (
	"aws-poc/internal/protocol"
	"regexp"
	"testing"
)

func TestGet(t *testing.T) {
	repo := randomRepository{}
	panRegex := regexp.MustCompile(`^\d{16}$`)

	if card, _ := repo.Get(&protocol.Dispute{}); card != nil && !panRegex.Match([]byte(card.Number)) {
		t.Error("card.Number is not a valid PAN")
	}
}
