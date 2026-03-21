package position

import (
	"github.com/KaiserSin/go-chess-ai/internal/domain/chess/internal/bitboard"
	chessmodel "github.com/KaiserSin/go-chess-ai/internal/domain/chess/model"
)

type Board struct {
	pieces    [2][6]uint64
	occupancy [2]uint64
	all       uint64
}

func newEmptyBoard() Board {
	return Board{}
}

func newInitialBoard() Board {
	board := newEmptyBoard()

	for file := 0; file < 8; file++ {
		board.placePiece(chessmodel.NewPiece(chessmodel.White, chessmodel.Pawn), squareMust(file, 1))
		board.placePiece(chessmodel.NewPiece(chessmodel.Black, chessmodel.Pawn), squareMust(file, 6))
	}

	backRank := [...]chessmodel.PieceType{
		chessmodel.Rook,
		chessmodel.Knight,
		chessmodel.Bishop,
		chessmodel.Queen,
		chessmodel.King,
		chessmodel.Bishop,
		chessmodel.Knight,
		chessmodel.Rook,
	}

	for file, pieceType := range backRank {
		board.placePiece(chessmodel.NewPiece(chessmodel.White, pieceType), squareMust(file, 0))
		board.placePiece(chessmodel.NewPiece(chessmodel.Black, pieceType), squareMust(file, 7))
	}

	return board
}

func (b Board) PieceAt(square chessmodel.Square) (chessmodel.Piece, bool) {
	return b.pieceAt(square)
}

func (b Board) pieceAt(square chessmodel.Square) (chessmodel.Piece, bool) {
	if !squareIsValid(square) {
		return chessmodel.Piece{}, false
	}

	mask := squareBitboard(square)
	for _, side := range allSides {
		for pieceType := chessmodel.Pawn; pieceType <= chessmodel.King; pieceType++ {
			if b.pieces[sideIndex(side)][bitboardIndex(pieceType)]&mask != 0 {
				return chessmodel.NewPiece(side, pieceType), true
			}
		}
	}

	return chessmodel.Piece{}, false
}

func (b *Board) placePiece(piece chessmodel.Piece, square chessmodel.Square) {
	if !squareIsValid(square) || !isValidSide(piece.Side()) {
		return
	}

	pieceType := piece.Type()
	if pieceType < chessmodel.Pawn || pieceType > chessmodel.King {
		return
	}

	b.removePiece(square)

	mask := squareBitboard(square)
	side := sideIndex(piece.Side())
	kind := bitboardIndex(pieceType)

	b.pieces[side][kind] |= mask
	b.occupancy[side] |= mask
	b.all |= mask
}

func (b *Board) removePiece(square chessmodel.Square) (chessmodel.Piece, bool) {
	if !squareIsValid(square) {
		return chessmodel.Piece{}, false
	}

	mask := squareBitboard(square)
	for _, side := range allSides {
		side := sideIndex(side)
		for pieceType := chessmodel.Pawn; pieceType <= chessmodel.King; pieceType++ {
			kind := bitboardIndex(pieceType)
			if b.pieces[side][kind]&mask == 0 {
				continue
			}

			b.pieces[side][kind] &^= mask
			b.occupancy[side] &^= mask
			b.all &^= mask

			return chessmodel.NewPiece(allSides[side], pieceType), true
		}
	}

	return chessmodel.Piece{}, false
}

func (b *Board) movePiece(from, to chessmodel.Square) (chessmodel.Piece, chessmodel.Piece, bool) {
	piece, ok := b.removePiece(from)
	if !ok {
		return chessmodel.Piece{}, chessmodel.Piece{}, false
	}

	captured, _ := b.removePiece(to)
	b.placePiece(piece, to)

	return piece, captured, true
}

func (b Board) occupied(square chessmodel.Square) bool {
	return b.all&squareBitboard(square) != 0
}

func (b Board) occupiedBy(side chessmodel.Side, square chessmodel.Square) bool {
	return b.occupancy[sideIndex(side)]&squareBitboard(square) != 0
}

func (b Board) kingSquare(side chessmodel.Side) (chessmodel.Square, bool) {
	index, ok := bitboard.First(b.pieces[sideIndex(side)][bitboardIndex(chessmodel.King)])
	if !ok {
		return 0, false
	}

	return chessmodel.Square(index), true
}

func (b Board) bitboard(side chessmodel.Side, pieceType chessmodel.PieceType) uint64 {
	return b.pieces[sideIndex(side)][bitboardIndex(pieceType)]
}

func forEachSquare(value uint64, visit func(chessmodel.Square)) {
	bitboard.ForEach(value, func(index int) {
		visit(chessmodel.Square(index))
	})
}
