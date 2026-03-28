package input_test

import (
	"testing"

	boardinput "github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/input"
	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/theme"
)

func TestSquareAt(t *testing.T) {
	translator := boardinput.NewTranslator(theme.NewTheme())

	a1, ok := translator.SquareAt(80, 688)
	if !ok {
		t.Fatal("want a1 hit")
	}

	if a1.Algebraic() != "a1" {
		t.Fatalf("want a1, got %s", a1.Algebraic())
	}

	h8, ok := translator.SquareAt(640, 128)
	if !ok {
		t.Fatal("want h8 hit")
	}

	if h8.Algebraic() != "h8" {
		t.Fatalf("want h8, got %s", h8.Algebraic())
	}
}

func TestSquareAtOutsideBoard(t *testing.T) {
	translator := boardinput.NewTranslator(theme.NewTheme())

	if _, ok := translator.SquareAt(10, 10); ok {
		t.Fatal("did not expect board hit outside board")
	}
}

func TestPromotionChoiceAt(t *testing.T) {
	uiTheme := theme.NewTheme()
	translator := boardinput.NewTranslator(uiTheme)
	rects := boardinput.PromotionOptionRects(uiTheme, 4)
	pieceTypes := []string{"queen", "rook", "bishop", "knight"}

	queenX := rects[0].X + rects[0].Width/2
	queenY := rects[0].Y + rects[0].Height/2

	choice, ok := translator.PromotionChoiceAt(queenX, queenY, pieceTypes)
	if !ok {
		t.Fatal("want promotion hit")
	}

	if choice != "queen" {
		t.Fatalf("want queen, got %s", choice)
	}

	knightX := rects[3].X + rects[3].Width/2
	knightY := rects[3].Y + rects[3].Height/2

	choice, ok = translator.PromotionChoiceAt(knightX, knightY, pieceTypes)
	if !ok {
		t.Fatal("want promotion hit")
	}

	if choice != "knight" {
		t.Fatalf("want knight, got %s", choice)
	}
}
