package gameplay

import (
	"testing"

	"github.com/KaiserSin/go-chess-ai/internal/application/ai"
	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestApplyAIMoveUsesFixedDepth(t *testing.T) {
	game := mustGameFromPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustSquareAt(t, 4, 0), chess.White, chess.King).
			Place(mustSquareAt(t, 3, 0), chess.White, chess.Queen).
			Place(mustSquareAt(t, 0, 0), chess.White, chess.Rook).
			Place(mustSquareAt(t, 4, 7), chess.Black, chess.King).
			Place(mustSquareAt(t, 3, 7), chess.Black, chess.Queen).
			Place(mustSquareAt(t, 0, 7), chess.Black, chess.Bishop),
	)
	service := newServiceWithGame(game)
	before := service.Snapshot()
	expected := ai.BestMove(game.Position())
	if !expected.HasMove {
		t.Fatal("want best move")
	}

	if err := service.ApplyAIMove(); err != nil {
		t.Fatalf("apply ai move: %v", err)
	}

	movedPiece := squareByAlgebraic(t, before, expected.Move.From.String()).PieceKey
	after := service.Snapshot()

	if square := squareByAlgebraic(t, after, expected.Move.From.String()); square.Occupied {
		t.Fatalf("did not expect piece on %s after ai move", expected.Move.From)
	}

	if square := squareByAlgebraic(t, after, expected.Move.To.String()); square.PieceKey != movedPiece {
		t.Fatalf("want %s on %s, got %q", movedPiece, expected.Move.To, square.PieceKey)
	}
}
