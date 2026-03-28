package theme

import "testing"

func TestPieceCatalogLookup(t *testing.T) {
	theme := NewTheme()

	if got := theme.PieceCatalog.Lookup("white-king").Label; got != "K" {
		t.Fatalf("want K, got %q", got)
	}

	if got := theme.PieceCatalog.Lookup("black-queen").Label; got != "q" {
		t.Fatalf("want q, got %q", got)
	}
}

func TestPieceCatalogLookupUnknownDoesNotPanic(t *testing.T) {
	theme := NewTheme()

	visual := theme.PieceCatalog.Lookup("mystery-piece")
	if visual.Label == "" {
		t.Fatal("want fallback label for unknown piece")
	}
}
