package chess

import chessposition "github.com/KaiserSin/go-chess-ai/internal/domain/chess/position"

func HasInsufficientMaterial(position Position) bool {
	return chessposition.HasInsufficientMaterial(position)
}

func IsFiftyMoveDraw(position Position) bool {
	return chessposition.IsFiftyMoveDraw(position)
}
