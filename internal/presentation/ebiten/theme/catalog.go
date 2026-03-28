package theme

import "image/color"

type PieceVisual struct {
	Key      string
	Label    string
	Color    color.RGBA
	AssetKey string
}

type PieceCatalog struct {
	visuals map[string]PieceVisual
	unknown PieceVisual
}

func NewPieceCatalog(whiteColor, blackColor color.RGBA) PieceCatalog {
	return PieceCatalog{
		visuals: map[string]PieceVisual{
			"white-king":   newPieceVisual("white-king", "K", whiteColor),
			"white-queen":  newPieceVisual("white-queen", "Q", whiteColor),
			"white-rook":   newPieceVisual("white-rook", "R", whiteColor),
			"white-bishop": newPieceVisual("white-bishop", "B", whiteColor),
			"white-knight": newPieceVisual("white-knight", "N", whiteColor),
			"white-pawn":   newPieceVisual("white-pawn", "P", whiteColor),
			"black-king":   newPieceVisual("black-king", "k", blackColor),
			"black-queen":  newPieceVisual("black-queen", "q", blackColor),
			"black-rook":   newPieceVisual("black-rook", "r", blackColor),
			"black-bishop": newPieceVisual("black-bishop", "b", blackColor),
			"black-knight": newPieceVisual("black-knight", "n", blackColor),
			"black-pawn":   newPieceVisual("black-pawn", "p", blackColor),
		},
		unknown: PieceVisual{
			Label:    "?",
			Color:    color.RGBA{R: 196, G: 196, B: 196, A: 255},
			AssetKey: "unknown-piece",
		},
	}
}

func (c PieceCatalog) Lookup(pieceKey string) PieceVisual {
	if pieceKey == "" {
		return PieceVisual{}
	}

	if visual, ok := c.visuals[pieceKey]; ok {
		return visual
	}

	unknown := c.unknown
	unknown.Key = pieceKey
	return unknown
}

func newPieceVisual(pieceKey, label string, color color.RGBA) PieceVisual {
	return PieceVisual{
		Key:      pieceKey,
		Label:    label,
		Color:    color,
		AssetKey: pieceKey,
	}
}
