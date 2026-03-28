package gameplay

import (
	"testing"

	"github.com/KaiserSin/go-chess-ai/internal/application/dto"
)

func TestSnapshotInitialPosition(t *testing.T) {
	service := NewService()

	snapshot := service.Snapshot()

	if got := len(snapshot.Squares); got != 64 {
		t.Fatalf("want 64 squares, got %d", got)
	}

	if snapshot.SideToMove != "white" {
		t.Fatalf("want white to move, got %q", snapshot.SideToMove)
	}

	if square := squareByAlgebraic(t, snapshot, "e1"); square.PieceKey != "white-king" {
		t.Fatalf("want white-king on e1, got %q", square.PieceKey)
	}

	if square := squareByAlgebraic(t, snapshot, "d1"); square.PieceKey != "white-queen" {
		t.Fatalf("want white-queen on d1, got %q", square.PieceKey)
	}

	if square := squareByAlgebraic(t, snapshot, "e8"); square.PieceKey != "black-king" {
		t.Fatalf("want black-king on e8, got %q", square.PieceKey)
	}

	if square := squareByAlgebraic(t, snapshot, "a2"); square.PieceKey != "white-pawn" {
		t.Fatalf("want white-pawn on a2, got %q", square.PieceKey)
	}

	if square := squareByAlgebraic(t, snapshot, "h7"); square.PieceKey != "black-pawn" {
		t.Fatalf("want black-pawn on h7, got %q", square.PieceKey)
	}
}

func squareByAlgebraic(t *testing.T, snapshot dto.GameSnapshot, algebraic string) dto.SquareSnapshot {
	t.Helper()

	for _, square := range snapshot.Squares {
		if square.Algebraic == algebraic {
			return square
		}
	}

	t.Fatalf("square %q not found", algebraic)
	return dto.SquareSnapshot{}
}
