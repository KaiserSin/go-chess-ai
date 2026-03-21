package chess

var (
	pawnCaptureFiles = [...]int{-1, 1}
	knightOffsets    = [][2]int{
		{1, 2}, {2, 1}, {2, -1}, {1, -2},
		{-1, -2}, {-2, -1}, {-2, 1}, {-1, 2},
	}
	kingOffsets = [][2]int{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0}, {1, 0},
		{-1, 1}, {0, 1}, {1, 1},
	}
	bishopDirections = [][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	rookDirections   = [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	queenDirections  = [][2]int{
		{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
		{1, 0}, {-1, 0}, {0, 1}, {0, -1},
	}
)

func backRank(side Side) int {
	if side == White {
		return 0
	}

	return 7
}

func pawnDirection(side Side) int {
	if side == White {
		return 1
	}

	return -1
}

func pawnStartRank(side Side) int {
	if side == White {
		return 1
	}

	return 6
}

func pawnPromotionRank(side Side) int {
	if side == White {
		return 7
	}

	return 0
}

func kingStartSquare(side Side) Square {
	return mustSquare(4, backRank(side))
}

func rookStartSquare(side Side, castle castleSide) Square {
	file := 0
	if castle == kingSide {
		file = 7
	}

	return mustSquare(file, backRank(side))
}

func castleKingSquare(side Side, castle castleSide) Square {
	file := 2
	if castle == kingSide {
		file = 6
	}

	return mustSquare(file, backRank(side))
}

func castleRookSquare(side Side, castle castleSide) Square {
	file := 3
	if castle == kingSide {
		file = 5
	}

	return mustSquare(file, backRank(side))
}

func castleTravelSquares(side Side, castle castleSide) []Square {
	rank := backRank(side)
	if castle == kingSide {
		return []Square{mustSquare(5, rank), mustSquare(6, rank)}
	}

	return []Square{mustSquare(1, rank), mustSquare(2, rank), mustSquare(3, rank)}
}

func castleSafeSquares(side Side, castle castleSide) []Square {
	rank := backRank(side)
	if castle == kingSide {
		return []Square{mustSquare(4, rank), mustSquare(5, rank), mustSquare(6, rank)}
	}

	return []Square{mustSquare(4, rank), mustSquare(3, rank), mustSquare(2, rank)}
}

func rookCastlingRef(square Square) (Side, castleSide, bool) {
	switch square {
	case rookStartSquare(White, queenSide):
		return White, queenSide, true
	case rookStartSquare(White, kingSide):
		return White, kingSide, true
	case rookStartSquare(Black, queenSide):
		return Black, queenSide, true
	case rookStartSquare(Black, kingSide):
		return Black, kingSide, true
	default:
		return White, kingSide, false
	}
}

func castleSideFromKingTarget(side Side, square Square) (castleSide, bool) {
	switch square {
	case castleKingSquare(side, kingSide):
		return kingSide, true
	case castleKingSquare(side, queenSide):
		return queenSide, true
	default:
		return kingSide, false
	}
}
