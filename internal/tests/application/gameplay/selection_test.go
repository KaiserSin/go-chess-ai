package gameplay_test

import (
	"testing"

	"github.com/KaiserSin/go-chess-ai/internal/application/gameplay"
)

func TestSelectSquareMarksSelectedPiece(t *testing.T) {
	service := gameplay.NewService()

	service.SelectSquareAt(4, 1)

	snapshot := service.Snapshot()
	if square := squareByAlgebraic(t, snapshot, "e2"); !square.Selected {
		t.Fatal("want e2 selected")
	}
}

func TestSelectSquareMarksLegalTargets(t *testing.T) {
	service := gameplay.NewService()

	service.SelectSquareAt(4, 1)

	snapshot := service.Snapshot()
	if square := squareByAlgebraic(t, snapshot, "e3"); !square.LegalTarget {
		t.Fatal("want e3 as legal target")
	}

	if square := squareByAlgebraic(t, snapshot, "e4"); !square.LegalTarget {
		t.Fatal("want e4 as legal target")
	}
}

func TestSelectingSameSquareClearsSelection(t *testing.T) {
	service := gameplay.NewService()

	service.SelectSquareAt(4, 1)
	service.SelectSquareAt(4, 1)

	snapshot := service.Snapshot()
	if square := squareByAlgebraic(t, snapshot, "e2"); square.Selected {
		t.Fatal("did not expect e2 selected")
	}
}
