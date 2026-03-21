package chess

import "math/bits"

func (p Position) hasInsufficientMaterial() bool {
	if p.board.bitboard(White, Pawn)|p.board.bitboard(Black, Pawn) != 0 {
		return false
	}

	if p.board.bitboard(White, Rook)|p.board.bitboard(Black, Rook) != 0 {
		return false
	}

	if p.board.bitboard(White, Queen)|p.board.bitboard(Black, Queen) != 0 {
		return false
	}

	whiteBishops := bits.OnesCount64(p.board.bitboard(White, Bishop))
	blackBishops := bits.OnesCount64(p.board.bitboard(Black, Bishop))
	whiteKnights := bits.OnesCount64(p.board.bitboard(White, Knight))
	blackKnights := bits.OnesCount64(p.board.bitboard(Black, Knight))
	totalMinors := whiteBishops + blackBishops + whiteKnights + blackKnights

	switch totalMinors {
	case 0:
		return true
	case 1:
		return true
	case 2:
		return whiteBishops == 1 &&
			blackBishops == 1 &&
			whiteKnights == 0 &&
			blackKnights == 0 &&
			p.bishopsShareColorComplex()
	default:
		return false
	}
}

func (p Position) bishopsShareColorComplex() bool {
	whiteSquare, whiteFound := firstSquare(p.board.bitboard(White, Bishop))
	blackSquare, blackFound := firstSquare(p.board.bitboard(Black, Bishop))
	if !whiteFound || !blackFound {
		return false
	}

	return whiteSquare.color() == blackSquare.color()
}

func firstSquare(bitboard uint64) (Square, bool) {
	if bitboard == 0 {
		return 0, false
	}

	var square Square
	found := false
	forEachSquare(bitboard, func(next Square) {
		if found {
			return
		}

		square = next
		found = true
	})

	return square, found
}
