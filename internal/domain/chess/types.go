package chess

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

func (s Side) index() int {
	return int(s)
}

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

type OutcomeReason uint8

const (
	NoOutcomeReason OutcomeReason = iota
	OutcomeByCheckmate
	OutcomeByStalemate
	OutcomeByThreefoldRepetition
	OutcomeByFiftyMoveRule
	OutcomeByInsufficientMaterial
)

func (r OutcomeReason) String() string {
	switch r {
	case NoOutcomeReason:
		return "none"
	case OutcomeByCheckmate:
		return "checkmate"
	case OutcomeByStalemate:
		return "stalemate"
	case OutcomeByThreefoldRepetition:
		return "same position 3 times"
	case OutcomeByFiftyMoveRule:
		return "50-move rule"
	case OutcomeByInsufficientMaterial:
		return "not enough material"
	default:
		return fmt.Sprintf("outcome_reason(%d)", r)
	}
}

type Outcome struct {
	reason    OutcomeReason
	winner    Side
	hasWinner bool
}

func noOutcome() Outcome {
	return Outcome{}
}

func decisiveOutcome(winner Side, reason OutcomeReason) Outcome {
	return Outcome{
		reason:    reason,
		winner:    winner,
		hasWinner: true,
	}
}

func drawOutcome(reason OutcomeReason) Outcome {
	return Outcome{
		reason: reason,
	}
}

func (o Outcome) Reason() OutcomeReason {
	return o.reason
}

func (o Outcome) Winner() (Side, bool) {
	return o.winner, o.hasWinner
}

func (o Outcome) IsFinished() bool {
	return o.reason != NoOutcomeReason
}

func (o Outcome) IsDraw() bool {
	return o.IsFinished() && !o.hasWinner
}

func (o Outcome) IsDecisive() bool {
	return o.hasWinner
}
