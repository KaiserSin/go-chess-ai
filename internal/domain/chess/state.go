package chess

import "fmt"

type CastlingRights uint8

const (
	WhiteKingSide CastlingRights = 1 << iota
	WhiteQueenSide
	BlackKingSide
	BlackQueenSide
)

type castleSide uint8

const (
	kingSide castleSide = iota
	queenSide
)

func initialCastlingRights() CastlingRights {
	return WhiteKingSide | WhiteQueenSide | BlackKingSide | BlackQueenSide
}

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

func (c *CastlingRights) removeKingside(side Side) {
	if side == White {
		*c &^= WhiteKingSide
		return
	}

	*c &^= BlackKingSide
}

func (c *CastlingRights) removeQueenside(side Side) {
	if side == White {
		*c &^= WhiteQueenSide
		return
	}

	*c &^= BlackQueenSide
}

func (c *CastlingRights) removeSide(side Side) {
	c.removeKingside(side)
	c.removeQueenside(side)
}

func (c *CastlingRights) remove(side Side, castle castleSide) {
	if castle == kingSide {
		c.removeKingside(side)
		return
	}

	c.removeQueenside(side)
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
