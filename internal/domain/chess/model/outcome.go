package model

import "fmt"

type OutcomeReason uint8

const (
	NoOutcomeReason OutcomeReason = iota
	OutcomeByCheckmate
	OutcomeByStalemate
	OutcomeByThreefoldRepetition
	OutcomeByFiftyMoveRule
	OutcomeByInsufficientMaterial
)

type Outcome struct {
	reason    OutcomeReason
	winner    Side
	hasWinner bool
}

func NoOutcome() Outcome {
	return Outcome{}
}

func NewDecisiveOutcome(winner Side, reason OutcomeReason) Outcome {
	return Outcome{
		reason:    reason,
		winner:    winner,
		hasWinner: true,
	}
}

func NewDrawOutcome(reason OutcomeReason) Outcome {
	return Outcome{
		reason: reason,
	}
}

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
