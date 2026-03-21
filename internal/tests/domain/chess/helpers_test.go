package chess_test

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

func mustNewGameFromPosition(t *testing.T, position chess.Position) *chess.Game {
	t.Helper()

	game, err := chess.NewGameFromPosition(position)
	if err != nil {
		t.Fatalf("NewGameFromPosition error: %v", err)
	}

	return game
}

func mustBuildGame(t *testing.T, builder *chess.PositionBuilder) *chess.Game {
	t.Helper()

	return mustNewGameFromPosition(t, mustBuildPosition(t, builder))
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

func assertPieceAt(t *testing.T, position chess.Position, rawSquare string, side chess.Side, pieceType chess.PieceType) {
	t.Helper()

	piece, ok := position.PieceAt(mustParseSquare(t, rawSquare))
	if !ok {
		t.Fatalf("want piece on %s", rawSquare)
	}

	if piece.Side() != side || piece.Type() != pieceType {
		t.Fatalf("want %s %s on %s, got %s", side, pieceType, rawSquare, piece)
	}
}
