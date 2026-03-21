package geom

var (
	PawnCaptureFiles = [...]int{-1, 1}
	KnightOffsets    = [][2]int{
		{1, 2}, {2, 1}, {2, -1}, {1, -2},
		{-1, -2}, {-2, -1}, {-2, 1}, {-1, 2},
	}
	KingOffsets = [][2]int{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0}, {1, 0},
		{-1, 1}, {0, 1}, {1, 1},
	}
	BishopDirections = [][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	RookDirections   = [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	QueenDirections  = [][2]int{
		{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
		{1, 0}, {-1, 0}, {0, 1}, {0, -1},
	}
)

func BackRank(isWhite bool) int {
	if isWhite {
		return 0
	}

	return 7
}

func PawnDirection(isWhite bool) int {
	if isWhite {
		return 1
	}

	return -1
}

func PawnStartRank(isWhite bool) int {
	if isWhite {
		return 1
	}

	return 6
}

func PawnPromotionRank(isWhite bool) int {
	if isWhite {
		return 7
	}

	return 0
}
