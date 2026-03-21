package position

import (
	"github.com/KaiserSin/go-chess-ai/internal/domain/chess/internal/bitboard"
	chessmodel "github.com/KaiserSin/go-chess-ai/internal/domain/chess/model"
)

func HasInsufficientMaterial(position Position) bool {
	if position.board.bitboard(chessmodel.White, chessmodel.Pawn)|position.board.bitboard(chessmodel.Black, chessmodel.Pawn) != 0 {
		return false
	}

	if position.board.bitboard(chessmodel.White, chessmodel.Rook)|position.board.bitboard(chessmodel.Black, chessmodel.Rook) != 0 {
		return false
	}

	if position.board.bitboard(chessmodel.White, chessmodel.Queen)|position.board.bitboard(chessmodel.Black, chessmodel.Queen) != 0 {
		return false
	}

	whiteBishops := bitboard.Count(position.board.bitboard(chessmodel.White, chessmodel.Bishop))
	blackBishops := bitboard.Count(position.board.bitboard(chessmodel.Black, chessmodel.Bishop))
	whiteKnights := bitboard.Count(position.board.bitboard(chessmodel.White, chessmodel.Knight))
	blackKnights := bitboard.Count(position.board.bitboard(chessmodel.Black, chessmodel.Knight))
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
			position.bishopsShareColorComplex()
	default:
		return false
	}
}

func (p Position) bishopsShareColorComplex() bool {
	whiteSquare, whiteFound := firstSquare(p.board.bitboard(chessmodel.White, chessmodel.Bishop))
	blackSquare, blackFound := firstSquare(p.board.bitboard(chessmodel.Black, chessmodel.Bishop))
	if !whiteFound || !blackFound {
		return false
	}

	return squareColor(whiteSquare) == squareColor(blackSquare)
}

func firstSquare(value uint64) (chessmodel.Square, bool) {
	index, ok := bitboard.First(value)
	if !ok {
		return 0, false
	}

	return chessmodel.Square(index), true
}
