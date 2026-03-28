package viewmodel_test

import (
	"testing"

	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/viewmodel"
)

func squareByAlgebraic(t *testing.T, board viewmodel.BoardViewModel, algebraic string) viewmodel.SquareViewModel {
	t.Helper()

	for _, square := range board.Squares {
		if square.Algebraic == algebraic {
			return square
		}
	}

	t.Fatalf("square %q not found", algebraic)
	return viewmodel.SquareViewModel{}
}
