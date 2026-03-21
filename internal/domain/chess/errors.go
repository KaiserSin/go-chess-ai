package chess

import chessmodel "github.com/KaiserSin/go-chess-ai/internal/domain/chess/model"

var (
	ErrInvalidSquare     = chessmodel.ErrInvalidSquare
	ErrInvalidPosition   = chessmodel.ErrInvalidPosition
	ErrNoPiece           = chessmodel.ErrNoPiece
	ErrWrongSide         = chessmodel.ErrWrongSide
	ErrInvalidMove       = chessmodel.ErrInvalidMove
	ErrGameFinished      = chessmodel.ErrGameFinished
	ErrPromotionRequired = chessmodel.ErrPromotionRequired
	ErrInvalidPromotion  = chessmodel.ErrInvalidPromotion
)
