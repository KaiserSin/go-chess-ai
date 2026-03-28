package viewmodel_test

import (
	"testing"

	"github.com/KaiserSin/go-chess-ai/internal/application/gameplay"
	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/theme"
	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/viewmodel"
)

func TestMapperInitialBoardOrientation(t *testing.T) {
	mapper := viewmodel.NewMapper(theme.NewTheme())
	snapshot := gameplay.NewService().Snapshot()

	board := mapper.Map(snapshot)

	a1 := squareByAlgebraic(t, board, "a1")
	if a1.X != 0 || a1.Y != 560 {
		t.Fatalf("want a1 at (0, 560), got (%d, %d)", a1.X, a1.Y)
	}

	h8 := squareByAlgebraic(t, board, "h8")
	if h8.X != 560 || h8.Y != 0 {
		t.Fatalf("want h8 at (560, 0), got (%d, %d)", h8.X, h8.Y)
	}
}

func TestMapperSquareCoordinatesUseSquareGrid(t *testing.T) {
	mapper := viewmodel.NewMapper(theme.NewTheme())
	snapshot := gameplay.NewService().Snapshot()

	board := mapper.Map(snapshot)

	for _, square := range board.Squares {
		if square.X%80 != 0 {
			t.Fatalf("want x coordinate divisible by 80, got %d for %s", square.X, square.Algebraic)
		}

		if square.Y%80 != 0 {
			t.Fatalf("want y coordinate divisible by 80, got %d for %s", square.Y, square.Algebraic)
		}
	}
}

func TestMapperSelectionAndLegalTargets(t *testing.T) {
	mapper := viewmodel.NewMapper(theme.NewTheme())
	service := gameplay.NewService()

	service.SelectSquareAt(4, 1)
	board := mapper.Map(service.Snapshot())

	if square := squareByAlgebraic(t, board, "e2"); !square.Selected {
		t.Fatal("want e2 selected")
	}

	if square := squareByAlgebraic(t, board, "e3"); !square.LegalTarget {
		t.Fatal("want e3 highlighted as legal target")
	}
}
