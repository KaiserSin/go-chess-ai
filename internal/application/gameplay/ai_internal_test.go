package gameplay

import (
	"testing"

	"github.com/KaiserSin/go-chess-ai/internal/application/dto"
	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestApplyAIMoveUsesSearchResult(t *testing.T) {
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

	if err := service.ApplyAIMove(); err != nil {
		t.Fatalf("apply ai move: %v", err)
	}

	after := service.Snapshot()
	if after.SideToMove != "black" {
		t.Fatalf("want black to move, got %q", after.SideToMove)
	}

	if changed := changedSquareCount(before, after); changed == 0 {
		t.Fatal("want board to change after ai move")
	}
}

func changedSquareCount(before, after dto.GameSnapshot) int {
	changed := 0
	afterBySquare := make(map[string]string, len(after.Squares))
	for _, square := range after.Squares {
		afterBySquare[square.Algebraic] = square.PieceKey
	}

	for _, square := range before.Squares {
		if afterBySquare[square.Algebraic] != square.PieceKey {
			changed++
		}
	}

	return changed
}
