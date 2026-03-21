package game

import (
	chessmodel "github.com/KaiserSin/go-chess-ai/internal/domain/chess/model"
	chessposition "github.com/KaiserSin/go-chess-ai/internal/domain/chess/position"
)

type Game struct {
	position chessposition.Position
	outcome  chessmodel.Outcome
	history  repetitionHistory
}

func NewGame() *Game {
	return newGameFromValidatedPosition(chessposition.NewInitialPosition())
}

func NewGameFromPosition(pos chessposition.Position) (*Game, error) {
	if err := chessposition.Validate(pos); err != nil {
		return nil, err
	}

	return newGameFromValidatedPosition(pos), nil
}

func newGameFromValidatedPosition(pos chessposition.Position) *Game {
	game := &Game{
		position: pos,
		history:  newRepetitionHistory(),
	}

	game.history.record(pos)
	game.outcome = game.resolveOutcome()

	return game
}

func (g *Game) Position() chessposition.Position {
	return g.position
}

func (g *Game) LegalMoves() []chessmodel.Move {
	if g.IsFinished() {
		return nil
	}

	return g.position.LegalMoves()
}

func (g *Game) ApplyMove(move chessmodel.Move) error {
	if g.IsFinished() {
		return chessmodel.ErrGameFinished
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

func (g *Game) Status() chessmodel.Status {
	return g.position.Status()
}

func (g *Game) Outcome() chessmodel.Outcome {
	return g.outcome
}

func (g *Game) IsFinished() bool {
	return g.outcome.IsFinished()
}

func (g *Game) resolveOutcome() chessmodel.Outcome {
	switch g.position.Status() {
	case chessmodel.Checkmate:
		return chessmodel.NewDecisiveOutcome(g.position.SideToMove().Opponent(), chessmodel.OutcomeByCheckmate)
	case chessmodel.Stalemate:
		return chessmodel.NewDrawOutcome(chessmodel.OutcomeByStalemate)
	}

	if chessposition.HasInsufficientMaterial(g.position) {
		return chessmodel.NewDrawOutcome(chessmodel.OutcomeByInsufficientMaterial)
	}

	if chessposition.IsFiftyMoveDraw(g.position) {
		return chessmodel.NewDrawOutcome(chessmodel.OutcomeByFiftyMoveRule)
	}

	if g.history.count(g.position) >= 3 {
		return chessmodel.NewDrawOutcome(chessmodel.OutcomeByThreefoldRepetition)
	}

	return chessmodel.NoOutcome()
}
