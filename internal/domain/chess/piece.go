package chess

import "fmt"

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

var promotionChoices = [...]PieceType{Queen, Rook, Bishop, Knight}

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

func (p Piece) Side() Side {
	return p.side
}

func (p Piece) Type() PieceType {
	return p.kind
}

func (p Piece) String() string {
	if !p.side.isValid() || !p.kind.isValid() || p.kind == NoPieceType {
		return "empty"
	}

	return fmt.Sprintf("%s %s", p.side, p.kind)
}
