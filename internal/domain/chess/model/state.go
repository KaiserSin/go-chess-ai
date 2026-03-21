package model

import "fmt"

type CastlingRights uint8

const (
	WhiteKingSide CastlingRights = 1 << iota
	WhiteQueenSide
	BlackKingSide
	BlackQueenSide
)

func (c CastlingRights) CanCastleKingside(side Side) bool {
	if side == White {
		return c&WhiteKingSide != 0
	}

	return c&BlackKingSide != 0
}

func (c CastlingRights) CanCastleQueenside(side Side) bool {
	if side == White {
		return c&WhiteQueenSide != 0
	}

	return c&BlackQueenSide != 0
}

func (c CastlingRights) WithoutKingside(side Side) CastlingRights {
	if side == White {
		return c &^ WhiteKingSide
	}

	return c &^ BlackKingSide
}

func (c CastlingRights) WithoutQueenside(side Side) CastlingRights {
	if side == White {
		return c &^ WhiteQueenSide
	}

	return c &^ BlackQueenSide
}

func (c CastlingRights) WithoutSide(side Side) CastlingRights {
	return c.WithoutKingside(side).WithoutQueenside(side)
}

type Status uint8

const (
	Ongoing Status = iota
	Check
	Checkmate
	Stalemate
)

func (s Status) String() string {
	switch s {
	case Ongoing:
		return "ongoing"
	case Check:
		return "check"
	case Checkmate:
		return "checkmate"
	case Stalemate:
		return "stalemate"
	default:
		return fmt.Sprintf("status(%d)", s)
	}
}
