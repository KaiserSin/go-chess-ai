package chess

import "fmt"

// Side identifies a chess side.
type Side uint8

const (
	White Side = iota
	Black
)

// Opponent returns the opposite side.
func (s Side) Opponent() Side {
	if s == Black {
		return White
	}

	return Black
}

// String returns a readable side name.
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

// PieceType identifies a chess piece type.
type PieceType uint8

const (
	NoPieceType PieceType = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

// String returns a readable piece type name.
func (pt PieceType) String() string {
	switch pt {
	case NoPieceType:
		return "none"
	case Pawn:
		return "pawn"
	case Knight:
		return "knight"
	case Bishop:
		return "bishop"
	case Rook:
		return "rook"
	case Queen:
		return "queen"
	case King:
		return "king"
	default:
		return fmt.Sprintf("piece_type(%d)", pt)
	}
}

func (pt PieceType) isValid() bool {
	return pt >= NoPieceType && pt <= King
}

func (pt PieceType) bitboardIndex() int {
	return int(pt - 1)
}

func (pt PieceType) isPromotionChoice() bool {
	switch pt {
	case Queen, Rook, Bishop, Knight:
		return true
	default:
		return false
	}
}

func (pt PieceType) symbol() string {
	switch pt {
	case Pawn:
		return "p"
	case Knight:
		return "n"
	case Bishop:
		return "b"
	case Rook:
		return "r"
	case Queen:
		return "q"
	case King:
		return "k"
	default:
		return "?"
	}
}

var promotionChoices = [...]PieceType{Queen, Rook, Bishop, Knight}

// Piece is a chess piece.
type Piece struct {
	side Side
	kind PieceType
}

func newPiece(side Side, kind PieceType) Piece {
	return Piece{
		side: side,
		kind: kind,
	}
}

// Side returns the owner of the piece.
func (p Piece) Side() Side {
	return p.side
}

// Type returns the kind of the piece.
func (p Piece) Type() PieceType {
	return p.kind
}

// String returns a readable piece description.
func (p Piece) String() string {
	if !p.side.isValid() || !p.kind.isValid() || p.kind == NoPieceType {
		return "empty"
	}

	return fmt.Sprintf("%s %s", p.side, p.kind)
}

// CastlingRights tracks which castling options remain available.
type CastlingRights uint8

const (
	WhiteKingSide CastlingRights = 1 << iota
	WhiteQueenSide
	BlackKingSide
	BlackQueenSide
)

func initialCastlingRights() CastlingRights {
	return WhiteKingSide | WhiteQueenSide | BlackKingSide | BlackQueenSide
}

// CanCastleKingside reports whether a side may castle kingside.
func (c CastlingRights) CanCastleKingside(side Side) bool {
	if side == White {
		return c&WhiteKingSide != 0
	}

	return c&BlackKingSide != 0
}

// CanCastleQueenside reports whether a side may castle queenside.
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

// Status represents the current state of play.
type Status uint8

const (
	Ongoing Status = iota
	Check
	Checkmate
	Stalemate
)

// String returns a readable status name.
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
