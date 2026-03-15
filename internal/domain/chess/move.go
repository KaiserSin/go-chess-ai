package chess

import "fmt"

// Move is a chess move from one square to another.
type Move struct {
	From      Square
	To        Square
	Promotion PieceType
}

// String returns a readable move representation.
func (m Move) String() string {
	if m.Promotion == NoPieceType {
		return fmt.Sprintf("%s%s", m.From, m.To)
	}

	return fmt.Sprintf("%s%s=%s", m.From, m.To, m.Promotion.symbol())
}

func (m Move) validateSquares() error {
	if !m.From.isValid() || !m.To.isValid() {
		return ErrInvalidSquare
	}

	return nil
}
