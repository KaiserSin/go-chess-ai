package position

import (
	"github.com/KaiserSin/go-chess-ai/internal/domain/chess/internal/geom"
	chessmodel "github.com/KaiserSin/go-chess-ai/internal/domain/chess/model"
)

func initialCastlingRights() chessmodel.CastlingRights {
	return chessmodel.WhiteKingSide |
		chessmodel.WhiteQueenSide |
		chessmodel.BlackKingSide |
		chessmodel.BlackQueenSide
}

func backRank(side chessmodel.Side) int {
	return geom.BackRank(side == chessmodel.White)
}

func pawnDirection(side chessmodel.Side) int {
	return geom.PawnDirection(side == chessmodel.White)
}

func pawnStartRank(side chessmodel.Side) int {
	return geom.PawnStartRank(side == chessmodel.White)
}

func pawnPromotionRank(side chessmodel.Side) int {
	return geom.PawnPromotionRank(side == chessmodel.White)
}

func kingStartSquare(side chessmodel.Side) chessmodel.Square {
	return squareMust(4, backRank(side))
}

func rookStartSquare(side chessmodel.Side, castle castleSide) chessmodel.Square {
	file := 0
	if castle == kingSide {
		file = 7
	}

	return squareMust(file, backRank(side))
}

func castleKingSquare(side chessmodel.Side, castle castleSide) chessmodel.Square {
	file := 2
	if castle == kingSide {
		file = 6
	}

	return squareMust(file, backRank(side))
}

func castleRookSquare(side chessmodel.Side, castle castleSide) chessmodel.Square {
	file := 3
	if castle == kingSide {
		file = 5
	}

	return squareMust(file, backRank(side))
}

func castleTravelSquares(side chessmodel.Side, castle castleSide) []chessmodel.Square {
	rank := backRank(side)
	if castle == kingSide {
		return []chessmodel.Square{squareMust(5, rank), squareMust(6, rank)}
	}

	return []chessmodel.Square{
		squareMust(1, rank),
		squareMust(2, rank),
		squareMust(3, rank),
	}
}

func castleSafeSquares(side chessmodel.Side, castle castleSide) []chessmodel.Square {
	rank := backRank(side)
	if castle == kingSide {
		return []chessmodel.Square{squareMust(4, rank), squareMust(5, rank), squareMust(6, rank)}
	}

	return []chessmodel.Square{squareMust(4, rank), squareMust(3, rank), squareMust(2, rank)}
}

func rookCastlingRef(square chessmodel.Square) (chessmodel.Side, castleSide, bool) {
	switch square {
	case rookStartSquare(chessmodel.White, queenSide):
		return chessmodel.White, queenSide, true
	case rookStartSquare(chessmodel.White, kingSide):
		return chessmodel.White, kingSide, true
	case rookStartSquare(chessmodel.Black, queenSide):
		return chessmodel.Black, queenSide, true
	case rookStartSquare(chessmodel.Black, kingSide):
		return chessmodel.Black, kingSide, true
	default:
		return chessmodel.White, kingSide, false
	}
}

func castleSideFromKingTarget(side chessmodel.Side, square chessmodel.Square) (castleSide, bool) {
	switch square {
	case castleKingSquare(side, kingSide):
		return kingSide, true
	case castleKingSquare(side, queenSide):
		return queenSide, true
	default:
		return kingSide, false
	}
}
