package chess

type Game struct {
	position Position
	outcome  Outcome
	history  repetitionHistory
}

func NewGame() *Game {
	game, err := NewGameFromPosition(NewInitialPosition())
	if err != nil {
		panic(err)
	}

	return game
}

func NewGameFromPosition(position Position) (*Game, error) {
	if err := position.validate(); err != nil {
		return nil, err
	}

	game := &Game{
		position: position,
		history:  newRepetitionHistory(),
	}

	game.history.record(position)
	game.outcome = game.resolveOutcome()

	return game, nil
}

func (g *Game) Position() Position {
	return g.position
}

func (g *Game) LegalMoves() []Move {
	if g.IsFinished() {
		return nil
	}

	return g.position.LegalMoves()
}

func (g *Game) ApplyMove(move Move) error {
	if g.IsFinished() {
		return ErrGameFinished
	}

	next, err := g.position.ApplyMove(move)
	if err != nil {
		return err
	}

	g.position = next
	g.history.record(next)
	g.outcome = g.resolveOutcome()

	return nil
}

func (g *Game) Status() Status {
	return g.position.Status()
}

func (g *Game) Outcome() Outcome {
	return g.outcome
}

func (g *Game) IsFinished() bool {
	return g.outcome.IsFinished()
}

func (g *Game) resolveOutcome() Outcome {
	switch g.position.Status() {
	case Checkmate:
		return decisiveOutcome(g.position.SideToMove().Opponent(), OutcomeByCheckmate)
	case Stalemate:
		return drawOutcome(OutcomeByStalemate)
	}

	if g.position.hasInsufficientMaterial() {
		return drawOutcome(OutcomeByInsufficientMaterial)
	}

	if g.position.isFiftyMoveDraw() {
		return drawOutcome(OutcomeByFiftyMoveRule)
	}

	if g.history.count(g.position) >= 3 {
		return drawOutcome(OutcomeByThreefoldRepetition)
	}

	return noOutcome()
}
