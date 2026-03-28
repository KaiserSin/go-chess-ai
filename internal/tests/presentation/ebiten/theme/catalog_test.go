package theme_test

import (
	"testing"

	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/theme"
)

func TestPieceCatalogLookup(t *testing.T) {
	uiTheme := theme.NewTheme()

	if got := uiTheme.PieceCatalog.Lookup("white-king").Label; got != "K" {
		t.Fatalf("want K, got %q", got)
	}

	if got := uiTheme.PieceCatalog.Lookup("black-queen").Label; got != "q" {
		t.Fatalf("want q, got %q", got)
	}
}

func TestPieceCatalogLookupUnknownDoesNotPanic(t *testing.T) {
	uiTheme := theme.NewTheme()

	visual := uiTheme.PieceCatalog.Lookup("mystery-piece")
	if visual.Label == "" {
		t.Fatal("want fallback label for unknown piece")
	}
}
