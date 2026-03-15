package chess

import "fmt"

const (
	boardFiles = 8
	boardRanks = 8
)

// Square identifies a board coordinate. File and rank are zero-based.
type Square uint8

// NewSquare creates a square from zero-based file and rank.
func NewSquare(file, rank int) (Square, error) {
	if file < 0 || file >= boardFiles || rank < 0 || rank >= boardRanks {
		return 0, ErrInvalidSquare
	}

	return Square(rank*boardFiles + file), nil
}

func mustSquare(file, rank int) Square {
	square, err := NewSquare(file, rank)
	if err != nil {
		panic(err)
	}

	return square
}

// ParseSquare parses algebraic notation such as "e4".
func ParseSquare(raw string) (Square, error) {
	if len(raw) != 2 {
		return 0, ErrInvalidSquare
	}

	file := int(raw[0] - 'a')
	rank := int(raw[1] - '1')

	return NewSquare(file, rank)
}

// File returns the zero-based file.
func (s Square) File() int {
	return int(s) % boardFiles
}

// Rank returns the zero-based rank.
func (s Square) Rank() int {
	return int(s) / boardFiles
}

func (s Square) isValid() bool {
	value := int(s)
	return value >= 0 && value < boardFiles*boardRanks
}

func (s Square) bitboard() uint64 {
	return uint64(1) << s
}

func (s Square) offset(fileDelta, rankDelta int) (Square, bool) {
	if !s.isValid() {
		return 0, false
	}

	file := s.File() + fileDelta
	rank := s.Rank() + rankDelta
	if file < 0 || file >= boardFiles || rank < 0 || rank >= boardRanks {
		return 0, false
	}

	return mustSquare(file, rank), true
}

// String formats the square as algebraic notation.
func (s Square) String() string {
	if !s.isValid() {
		return "<invalid>"
	}

	return fmt.Sprintf("%c%d", 'a'+s.File(), s.Rank()+1)
}

type optionalSquare struct {
	value Square
	ok    bool
}

func noSquare() optionalSquare {
	return optionalSquare{}
}

func someSquare(square Square) optionalSquare {
	return optionalSquare{
		value: square,
		ok:    true,
	}
}
