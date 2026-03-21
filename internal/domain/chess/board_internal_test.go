package chess

import "testing"

func TestBoardPieceAt(t *testing.T) {
	board := newEmptyBoard()

	if _, ok := board.pieceAt(Square(99)); ok {
		t.Fatal("want no piece")
	}

	board.placePiece(newPiece(White, Rook), mustSquare(1, 1))
	piece, ok := board.pieceAt(mustSquare(1, 1))
	if !ok {
		t.Fatal("want piece")
	}

	if piece.kind != Rook || piece.side != White {
		t.Fatalf("want white rook, got %v", piece)
	}
}

func TestBoardPlacePiece(t *testing.T) {
	board := newEmptyBoard()

	board.placePiece(Piece{}, mustSquare(1, 1))
	if _, ok := board.pieceAt(mustSquare(1, 1)); ok {
		t.Fatal("want empty square")
	}

	board.placePiece(newPiece(White, Rook), Square(99))
	if board.all != 0 {
		t.Fatal("want empty board")
	}
}

func TestBoardRemovePiece(t *testing.T) {
	board := newEmptyBoard()

	if _, ok := board.removePiece(Square(99)); ok {
		t.Fatal("want no piece")
	}

	if _, ok := board.removePiece(mustSquare(0, 0)); ok {
		t.Fatal("want no piece")
	}

	board.placePiece(newPiece(White, Rook), mustSquare(0, 0))
	piece, ok := board.removePiece(mustSquare(0, 0))
	if !ok {
		t.Fatal("want piece")
	}

	if piece.kind != Rook || piece.side != White {
		t.Fatalf("want white rook, got %v", piece)
	}
}

func TestBoardMovePiece(t *testing.T) {
	board := newEmptyBoard()
	board.placePiece(newPiece(White, King), mustSquare(4, 0))
	board.placePiece(newPiece(Black, Knight), mustSquare(4, 7))

	moved, taken, ok := board.movePiece(mustSquare(4, 0), mustSquare(4, 7))
	if !ok {
		t.Fatal("want move")
	}

	if moved.kind != King || moved.side != White {
		t.Fatalf("want white king, got %v", moved)
	}

	if taken.kind != Knight || taken.side != Black {
		t.Fatalf("want black knight, got %v", taken)
	}

	if _, _, ok := board.movePiece(mustSquare(4, 0), mustSquare(4, 1)); ok {
		t.Fatal("want bad move")
	}
}

func TestBoardKingSquare(t *testing.T) {
	board := newEmptyBoard()
	if _, ok := board.kingSquare(Black); ok {
		t.Fatal("want no king")
	}

	board.placePiece(newPiece(White, King), mustSquare(4, 7))
	square, ok := board.kingSquare(White)
	if !ok || square != mustSquare(4, 7) {
		t.Fatalf("want e8, got %v, %t", square, ok)
	}
}
