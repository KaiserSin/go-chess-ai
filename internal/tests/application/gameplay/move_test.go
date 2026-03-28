package gameplay_test

import (
	"testing"

	"github.com/KaiserSin/go-chess-ai/internal/application/gameplay"
)

func TestTryMoveUpdatesBoardAndSide(t *testing.T) {
	service := gameplay.NewService()

	service.SelectSquareAt(4, 1)
	if err := service.TryMoveAt(4, 3); err != nil {
		t.Fatalf("want legal move, got %v", err)
	}

	snapshot := service.Snapshot()
	if snapshot.SideToMove != "black" {
		t.Fatalf("want black to move, got %q", snapshot.SideToMove)
	}

	if square := squareByAlgebraic(t, snapshot, "e4"); square.PieceKey != "white-pawn" {
		t.Fatalf("want white-pawn on e4, got %q", square.PieceKey)
	}

	if square := squareByAlgebraic(t, snapshot, "e2"); square.Occupied {
		t.Fatal("did not expect piece on e2")
	}
}

func TestMoveClearsSelection(t *testing.T) {
	service := gameplay.NewService()

	service.SelectSquareAt(4, 1)
	if err := service.TryMoveAt(4, 3); err != nil {
		t.Fatalf("want legal move, got %v", err)
	}

	snapshot := service.Snapshot()
	if square := squareByAlgebraic(t, snapshot, "e2"); square.Selected {
		t.Fatal("did not expect e2 selected after move")
	}

	if square := squareByAlgebraic(t, snapshot, "e4"); square.Selected {
		t.Fatal("did not expect e4 selected after move")
	}
}
