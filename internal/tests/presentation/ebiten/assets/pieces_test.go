package assets_test

import (
	"testing"

	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/assets"
)

func TestLoadPieceSprites(t *testing.T) {
	sprites := assets.LoadPieceSprites()

	if got := sprites.Count(); got != 12 {
		t.Fatalf("want 12 sprites, got %d", got)
	}

	if _, ok := sprites.Lookup("white-king"); !ok {
		t.Fatal("want white-king sprite")
	}

	if _, ok := sprites.Lookup("black-queen"); !ok {
		t.Fatal("want black-queen sprite")
	}

	if _, ok := sprites.Lookup("unknown-piece"); ok {
		t.Fatal("did not expect unknown-piece sprite")
	}
}
