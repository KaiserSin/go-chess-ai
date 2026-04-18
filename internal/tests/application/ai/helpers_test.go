package ai_test

import (
	"testing"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func mustBuildPosition(t *testing.T, builder *chess.PositionBuilder) chess.Position {
	t.Helper()

	position, err := builder.Build()
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	return position
}

func mustParseSquare(t *testing.T, raw string) chess.Square {
	t.Helper()

	square, err := chess.ParseSquare(raw)
	if err != nil {
		t.Fatalf("ParseSquare(%q) error: %v", raw, err)
	}

	return square
}

func mustMove(t *testing.T, from, to string, promotion ...chess.PieceType) chess.Move {
	t.Helper()

	move := chess.Move{
		From: mustParseSquare(t, from),
		To:   mustParseSquare(t, to),
	}

	if len(promotion) > 0 {
		move.Promotion = promotion[0]
	}

	return move
}

func applyMoves(t *testing.T, game *chess.Game, moves ...chess.Move) {
	t.Helper()

	for _, move := range moves {
		if err := game.ApplyMove(move); err != nil {
			t.Fatalf("ApplyMove(%s) error: %v", move, err)
		}
	}
}

func containsMove(moves []chess.Move, target chess.Move) bool {
	for _, move := range moves {
		if move == target {
			return true
		}
	}

	return false
}

func mustApplyMove(t *testing.T, position chess.Position, move chess.Move) chess.Position {
	t.Helper()

	next, err := position.ApplyMove(move)
	if err != nil {
		t.Fatalf("ApplyMove(%s) error: %v", move, err)
	}

	return next
}

func hasImmediateMate(position chess.Position) bool {
	attacker := position.SideToMove()

	for _, move := range position.LegalMoves() {
		next, err := position.ApplyMove(move)
		if err != nil {
			panic(err)
		}

		if next.Status() == chess.Checkmate && next.SideToMove().Opponent() == attacker {
			return true
		}
	}

	return false
}

func canForceMate(position chess.Position, plies int, attacker chess.Side) bool {
	if position.Status() == chess.Checkmate {
		return position.SideToMove().Opponent() == attacker
	}

	if isDrawnTerminal(position) || plies == 0 {
		return false
	}

	moves := position.LegalMoves()
	if len(moves) == 0 {
		return false
	}

	if position.SideToMove() == attacker {
		for _, move := range moves {
			next, err := position.ApplyMove(move)
			if err != nil {
				panic(err)
			}

			if canForceMate(next, plies-1, attacker) {
				return true
			}
		}

		return false
	}

	for _, move := range moves {
		next, err := position.ApplyMove(move)
		if err != nil {
			panic(err)
		}

		if !canForceMate(next, plies-1, attacker) {
			return false
		}
	}

	return true
}

func isDrawnTerminal(position chess.Position) bool {
	status := position.Status()
	if status == chess.Stalemate {
		return true
	}

	return chess.HasInsufficientMaterial(position) || chess.IsFiftyMoveDraw(position)
}
