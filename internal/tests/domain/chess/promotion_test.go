package chess_test

import (
	"errors"
	"testing"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestPromotionRequired(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "g7"), chess.White, chess.Pawn),
	)

	_, err := position.ApplyMove(mustMove(t, "g7", "g8"))
	if !errors.Is(err, chess.ErrPromotionRequired) {
		t.Fatalf("want promotion error, got %v", err)
	}
}

func TestBadPromotion(t *testing.T) {
	t.Run("bad piece", func(t *testing.T) {
		position := mustBuildPosition(t,
			chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
				Place(mustParseSquare(t, "g7"), chess.White, chess.Pawn),
		)

		_, err := position.ApplyMove(mustMove(t, "g7", "g8", chess.King))
		if !errors.Is(err, chess.ErrInvalidPromotion) {
			t.Fatalf("want bad promotion error, got %v", err)
		}
	})

	t.Run("bad type", func(t *testing.T) {
		position := mustBuildPosition(t,
			chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
				Place(mustParseSquare(t, "g6"), chess.White, chess.Pawn),
		)

		_, err := position.ApplyMove(mustMove(t, "g6", "g7", chess.PieceType(9)))
		if !errors.Is(err, chess.ErrInvalidPromotion) {
			t.Fatalf("want bad promotion error, got %v", err)
		}
	})

	t.Run("too early", func(t *testing.T) {
		position := mustBuildPosition(t,
			chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
				Place(mustParseSquare(t, "g6"), chess.White, chess.Pawn),
		)

		_, err := position.ApplyMove(mustMove(t, "g6", "g7", chess.Queen))
		if !errors.Is(err, chess.ErrInvalidPromotion) {
			t.Fatalf("want bad promotion error, got %v", err)
		}
	})

	t.Run("rook move", func(t *testing.T) {
		position := mustBuildPosition(t,
			chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
				Place(mustParseSquare(t, "a1"), chess.White, chess.Rook),
		)

		_, err := position.ApplyMove(mustMove(t, "a1", "a2", chess.Queen))
		if !errors.Is(err, chess.ErrInvalidPromotion) {
			t.Fatalf("want bad promotion error, got %v", err)
		}
	})
}

func TestPromotionWorks(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "g7"), chess.White, chess.Pawn),
	)

	next, err := position.ApplyMove(mustMove(t, "g7", "g8", chess.Queen))
	if err != nil {
		t.Fatalf("want good promotion, got %v", err)
	}

	assertPieceAt(t, next, "g8", chess.White, chess.Queen)
}
