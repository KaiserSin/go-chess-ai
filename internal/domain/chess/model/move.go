package model

import "fmt"

type Move struct {
	From      Square
	To        Square
	Promotion PieceType
}

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
