package chess

import "fmt"

type Side uint8

const (
	White Side = iota
	Black
)

var allSides = [...]Side{White, Black}

func (s Side) Opponent() Side {
	if s == Black {
		return White
	}

	return Black
}

func (s Side) String() string {
	switch s {
	case White:
		return "white"
	case Black:
		return "black"
	default:
		return fmt.Sprintf("side(%d)", s)
	}
}

func (s Side) isValid() bool {
	return s == White || s == Black
}

func (s Side) index() int {
	return int(s)
}
