package gameplay

import (
	"testing"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestApplyAIMoveUsesConfiguredDepth(t *testing.T) {
	depthOne := newServiceWithGame(mustGameFromPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustSquareAt(t, 4, 0), chess.White, chess.King).
			Place(mustSquareAt(t, 3, 0), chess.White, chess.Queen).
			Place(mustSquareAt(t, 4, 7), chess.Black, chess.King).
			Place(mustSquareAt(t, 3, 7), chess.Black, chess.Rook),
	))
	depthOne.SetAISearchDepth(1)

	if err := depthOne.ApplyAIMove(); err != nil {
		t.Fatalf("depth 1 apply ai move: %v", err)
	}

	depthOneSnapshot := depthOne.Snapshot()
	if square := squareByAlgebraic(t, depthOneSnapshot, "d8"); square.PieceKey != "white-queen" {
		t.Fatalf("want white-queen on d8 after depth 1, got %q", square.PieceKey)
	}

	depthTwo := newServiceWithGame(mustGameFromPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustSquareAt(t, 4, 0), chess.White, chess.King).
			Place(mustSquareAt(t, 3, 0), chess.White, chess.Queen).
			Place(mustSquareAt(t, 4, 7), chess.Black, chess.King).
			Place(mustSquareAt(t, 3, 7), chess.Black, chess.Rook),
	))
	depthTwo.SetAISearchDepth(2)

	if err := depthTwo.ApplyAIMove(); err != nil {
		t.Fatalf("depth 2 apply ai move: %v", err)
	}

	depthTwoSnapshot := depthTwo.Snapshot()
	if square := squareByAlgebraic(t, depthTwoSnapshot, "d8"); square.PieceKey != "black-rook" {
		t.Fatalf("want black-rook to stay on d8 after depth 2, got %q", square.PieceKey)
	}
}
