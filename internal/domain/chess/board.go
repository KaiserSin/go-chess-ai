package chess

import "math/bits"

// Board stores piece placement using hidden bitboards.
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

	for file := 0; file < boardFiles; file++ {
		board.placePiece(newPiece(White, Pawn), mustSquare(file, 1))
		board.placePiece(newPiece(Black, Pawn), mustSquare(file, 6))
	}

	backRank := [...]PieceType{Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook}
	for file, pieceType := range backRank {
		board.placePiece(newPiece(White, pieceType), mustSquare(file, 0))
		board.placePiece(newPiece(Black, pieceType), mustSquare(file, 7))
	}

	return board
}

// PieceAt returns the piece located on the given square.
func (b Board) PieceAt(square Square) (Piece, bool) {
	return b.pieceAt(square)
}

func (b Board) pieceAt(square Square) (Piece, bool) {
	if !square.isValid() {
		return Piece{}, false
	}

	mask := square.bitboard()
	for _, side := range []Side{White, Black} {
		for pieceType := Pawn; pieceType <= King; pieceType++ {
			if b.pieces[side.index()][pieceType.bitboardIndex()]&mask != 0 {
				return newPiece(side, pieceType), true
			}
		}
	}

	return Piece{}, false
}

func (b *Board) placePiece(piece Piece, square Square) {
	if !square.isValid() || !piece.side.isValid() || piece.kind < Pawn || piece.kind > King {
		return
	}

	b.removePiece(square)

	mask := square.bitboard()
	sideIndex := piece.side.index()
	pieceIndex := piece.kind.bitboardIndex()

	b.pieces[sideIndex][pieceIndex] |= mask
	b.occupancy[sideIndex] |= mask
	b.all |= mask
}

func (b *Board) removePiece(square Square) (Piece, bool) {
	if !square.isValid() {
		return Piece{}, false
	}

	mask := square.bitboard()
	for _, side := range []Side{White, Black} {
		sideIndex := side.index()
		for pieceType := Pawn; pieceType <= King; pieceType++ {
			pieceIndex := pieceType.bitboardIndex()
			if b.pieces[sideIndex][pieceIndex]&mask == 0 {
				continue
			}

			b.pieces[sideIndex][pieceIndex] &^= mask
			b.occupancy[sideIndex] &^= mask
			b.all &^= mask

			return newPiece(side, pieceType), true
		}
	}

	return Piece{}, false
}

func (b *Board) movePiece(from, to Square) (Piece, Piece, bool) {
	piece, ok := b.removePiece(from)
	if !ok {
		return Piece{}, Piece{}, false
	}

	captured, _ := b.removePiece(to)
	b.placePiece(piece, to)

	return piece, captured, true
}

func (b Board) occupied(square Square) bool {
	return b.all&square.bitboard() != 0
}

func (b Board) occupiedBy(side Side, square Square) bool {
	return b.occupancy[side.index()]&square.bitboard() != 0
}

func (b Board) kingSquare(side Side) (Square, bool) {
	bitboard := b.pieces[side.index()][King.bitboardIndex()]
	if bitboard == 0 {
		return 0, false
	}

	return Square(bits.TrailingZeros64(bitboard)), true
}

func (b Board) bitboard(side Side, pieceType PieceType) uint64 {
	return b.pieces[side.index()][pieceType.bitboardIndex()]
}

func forEachSquare(bitboard uint64, visit func(Square)) {
	for bitboard != 0 {
		index := bits.TrailingZeros64(bitboard)
		visit(Square(index))
		bitboard &= bitboard - 1
	}
}
