package model

import "fmt"

const (
	boardFiles = 8
	boardRanks = 8
)

type Square uint8

func NewSquare(file, rank int) (Square, error) {
	if file < 0 || file >= boardFiles || rank < 0 || rank >= boardRanks {
		return 0, ErrInvalidSquare
	}

	return Square(rank*boardFiles + file), nil
}

func ParseSquare(raw string) (Square, error) {
	if len(raw) != 2 {
		return 0, ErrInvalidSquare
	}

	file := int(raw[0] - 'a')
	rank := int(raw[1] - '1')

	return NewSquare(file, rank)
}

func (s Square) File() int {
	return int(s) % boardFiles
}

func (s Square) Rank() int {
	return int(s) / boardFiles
}

func (s Square) String() string {
	if !s.isValid() {
		return "<invalid>"
	}

	return fmt.Sprintf("%c%d", 'a'+s.File(), s.Rank()+1)
}

func (s Square) isValid() bool {
	return int(s) < boardFiles*boardRanks
}
