package ai

import (
	"testing"
	"time"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestExpiredDeadlineReturnsLegalFallbackMove(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "d1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook),
	)

	result := bestMoveWithOptions(position, searchOptions{
		deadline:      time.Now(),
		useAspiration: true,
	})
	if !result.HasMove {
		t.Fatal("want fallback move")
	}

	if !containsMove(position.LegalMoves(), result.Move) {
		t.Fatalf("want legal fallback move, got %s", result.Move)
	}
}

func TestAspirationSearchMatchesFullWindowSearch(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e2"), chess.White, chess.King).
			Place(mustParseSquare(t, "f2"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "b5"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "h4"), chess.Black, chess.King),
	)

	deadline := time.Now().Add(10 * time.Second)
	withAspiration := bestMoveWithOptions(position, searchOptions{
		deadline:      deadline,
		useAspiration: true,
	})
	withoutAspiration := bestMoveWithOptions(position, searchOptions{
		deadline:      deadline,
		useAspiration: false,
	})

	if withAspiration != withoutAspiration {
		t.Fatalf("want matching search results, with aspiration=%+v without aspiration=%+v", withAspiration, withoutAspiration)
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

func mustApplyMove(t *testing.T, position chess.Position, move chess.Move) chess.Position {
	t.Helper()

	next, err := position.ApplyMove(move)
	if err != nil {
		t.Fatalf("ApplyMove(%s) error: %v", move, err)
	}

	return next
}

func containsMove(moves []chess.Move, target chess.Move) bool {
	for _, move := range moves {
		if move == target {
			return true
		}
	}

	return false
}
