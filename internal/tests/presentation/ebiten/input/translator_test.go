package input_test

import (
	"testing"

	boardinput "github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/input"
	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/theme"
)

func TestSquareAt(t *testing.T) {
	translator := boardinput.NewTranslator(theme.NewTheme())

	a1, ok := translator.SquareAt(80, 688, false)
	if !ok {
		t.Fatal("want a1 hit")
	}

	if a1.Algebraic() != "a1" {
		t.Fatalf("want a1, got %s", a1.Algebraic())
	}

	h8, ok := translator.SquareAt(640, 128, false)
	if !ok {
		t.Fatal("want h8 hit")
	}

	if h8.Algebraic() != "h8" {
		t.Fatalf("want h8, got %s", h8.Algebraic())
	}
}

func TestSquareAtOutsideBoard(t *testing.T) {
	translator := boardinput.NewTranslator(theme.NewTheme())

	if _, ok := translator.SquareAt(10, 10, false); ok {
		t.Fatal("did not expect board hit outside board")
	}
}

func TestSquareAtBlackPerspective(t *testing.T) {
	translator := boardinput.NewTranslator(theme.NewTheme())

	h8, ok := translator.SquareAt(80, 688, true)
	if !ok {
		t.Fatal("want h8 hit")
	}

	if h8.Algebraic() != "h8" {
		t.Fatalf("want h8, got %s", h8.Algebraic())
	}

	a1, ok := translator.SquareAt(640, 128, true)
	if !ok {
		t.Fatal("want a1 hit")
	}

	if a1.Algebraic() != "a1" {
		t.Fatalf("want a1, got %s", a1.Algebraic())
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

func TestSideChoiceAt(t *testing.T) {
	uiTheme := theme.NewTheme()
	translator := boardinput.NewTranslator(uiTheme)
	choices := boardinput.SideChoiceRects(uiTheme.WindowWidth)

	whiteRect := choices[0].Rect
	whiteX := whiteRect.X + whiteRect.Width/2
	whiteY := whiteRect.Y + whiteRect.Height/2

	choice, ok := translator.SideChoiceAt(whiteX, whiteY)
	if !ok {
		t.Fatal("want white button hit")
	}

	if choice != "white" {
		t.Fatalf("want white, got %s", choice)
	}

	blackRect := choices[1].Rect
	blackX := blackRect.X + blackRect.Width/2
	blackY := blackRect.Y + blackRect.Height/2

	choice, ok = translator.SideChoiceAt(blackX, blackY)
	if !ok {
		t.Fatal("want black button hit")
	}

	if choice != "black" {
		t.Fatalf("want black, got %s", choice)
	}
}

func TestDepthInputAt(t *testing.T) {
	uiTheme := theme.NewTheme()
	translator := boardinput.NewTranslator(uiTheme)

	rect := boardinput.DepthInputRect(uiTheme.WindowWidth)
	centerX := rect.X + rect.Width/2
	centerY := rect.Y + rect.Height/2

	if !translator.DepthInputAt(centerX, centerY) {
		t.Fatal("want depth input hit")
	}

	if translator.DepthInputAt(10, 10) {
		t.Fatal("did not expect depth input hit outside field")
	}
}
