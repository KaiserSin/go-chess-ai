package model

import "fmt"

type Side uint8

const (
	White Side = iota
	Black
)

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
