package ai_test

import (
	"testing"

	"github.com/KaiserSin/go-chess-ai/internal/application/ai"
	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestEvaluateInitialPosition(t *testing.T) {
	position := chess.NewInitialPosition()

	if got := ai.Evaluate(position, chess.White); got != 0 {
		t.Fatalf("want 0 for white, got %d", got)
	}

	if got := ai.Evaluate(position, chess.Black); got != 0 {
		t.Fatalf("want 0 for black, got %d", got)
	}
}

func TestEvaluateMaterialAdvantageForWhite(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "d1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	if got := ai.Evaluate(position, chess.White); got != 900 {
		t.Fatalf("want 900 for white, got %d", got)
	}

	if got := ai.Evaluate(position, chess.Black); got != -900 {
		t.Fatalf("want -900 for black, got %d", got)
	}
}

func TestEvaluateMaterialAdvantageForBlack(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook),
	)

	if got := ai.Evaluate(position, chess.White); got != -500 {
		t.Fatalf("want -500 for white, got %d", got)
	}

	if got := ai.Evaluate(position, chess.Black); got != 500 {
		t.Fatalf("want 500 for black, got %d", got)
	}
}

func TestEvaluateCheckmate(t *testing.T) {
	game := chess.NewGame()
	applyMoves(t, game,
		mustMove(t, "f2", "f3"),
		mustMove(t, "e7", "e5"),
		mustMove(t, "g2", "g4"),
		mustMove(t, "d8", "h4"),
	)

	position := game.Position()
	if got := ai.Evaluate(position, chess.Black); got != 100000 {
		t.Fatalf("want 100000 for black, got %d", got)
	}

	if got := ai.Evaluate(position, chess.White); got != -100000 {
		t.Fatalf("want -100000 for white, got %d", got)
	}
}

func TestEvaluateStalemateReturnsZero(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.Black).
			Place(mustParseSquare(t, "c6"), chess.White, chess.King).
			Place(mustParseSquare(t, "b6"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.King),
	)

	if got := ai.Evaluate(position, chess.White); got != 0 {
		t.Fatalf("want 0 for white, got %d", got)
	}

	if got := ai.Evaluate(position, chess.Black); got != 0 {
		t.Fatalf("want 0 for black, got %d", got)
	}
}

func TestEvaluateInsufficientMaterialReturnsZero(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	if got := ai.Evaluate(position, chess.White); got != 0 {
		t.Fatalf("want 0 for white, got %d", got)
	}

	if got := ai.Evaluate(position, chess.Black); got != 0 {
		t.Fatalf("want 0 for black, got %d", got)
	}
}

func TestEvaluateFiftyMoveDrawReturnsZero(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			WithHalfmoveClock(100).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	if got := ai.Evaluate(position, chess.White); got != 0 {
		t.Fatalf("want 0 for white, got %d", got)
	}

	if got := ai.Evaluate(position, chess.Black); got != 0 {
		t.Fatalf("want 0 for black, got %d", got)
	}
}

func TestEvaluateInvalidPerspectiveReturnsZero(t *testing.T) {
	position := chess.NewInitialPosition()

	if got := ai.Evaluate(position, chess.Side(9)); got != 0 {
		t.Fatalf("want 0 for invalid side, got %d", got)
	}
}

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
