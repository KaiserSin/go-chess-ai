package chess

// Game is the aggregate root for chess play.
type Game struct {
	position Position
}

// NewGame creates a new game in the starting position.
func NewGame() *Game {
	return &Game{
		position: NewInitialPosition(),
	}
}

// Position returns the current position.
func (g *Game) Position() Position {
	return g.position
}

// LegalMoves returns all legal moves from the current position.
func (g *Game) LegalMoves() []Move {
	return g.position.LegalMoves()
}

// ApplyMove applies a move to the game.
func (g *Game) ApplyMove(move Move) error {
	next, err := g.position.ApplyMove(move)
	if err != nil {
		return err
	}

	g.position = next
	return nil
}

// Status returns the current game status.
func (g *Game) Status() Status {
	return g.position.Status()
}
