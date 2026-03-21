package position

import (
	"errors"
	"testing"

	chessmodel "github.com/KaiserSin/go-chess-ai/internal/domain/chess/model"
)

func mustTestPosition(t *testing.T, builder *PositionBuilder) Position {
	t.Helper()

	position, err := builder.Build()
	if err != nil {
		t.Fatalf("want good position, got %v", err)
	}

	return position
}

func TestPositionValidate(t *testing.T) {
	position := Position{
		board:      newEmptyBoard(),
		sideToMove: chessmodel.Side(9),
	}

	if !errors.Is(position.validate(), chessmodel.ErrInvalidPosition) {
		t.Fatal("want bad position")
	}

	position = Position{
		board:      newEmptyBoard(),
		sideToMove: chessmodel.White,
	}
	position.board.placePiece(chessmodel.NewPiece(chessmodel.Black, chessmodel.King), squareMust(4, 7))
	if !errors.Is(position.validate(), chessmodel.ErrInvalidPosition) {
		t.Fatal("want bad position")
	}

	position = Position{
		board:      newEmptyBoard(),
		sideToMove: chessmodel.White,
	}
	position.board.placePiece(chessmodel.NewPiece(chessmodel.White, chessmodel.King), squareMust(4, 0))
	if !errors.Is(position.validate(), chessmodel.ErrInvalidPosition) {
		t.Fatal("want bad position")
	}
}

func TestMoveCastlingRook(t *testing.T) {
	position := Position{board: newEmptyBoard()}
	position.moveCastlingRook(chessmodel.White, chessmodel.Move{From: squareMust(4, 0), To: squareMust(4, 1)})
	if position.board.all != 0 {
		t.Fatal("want empty board")
	}

	position = Position{board: newEmptyBoard()}
	position.moveCastlingRook(chessmodel.White, chessmodel.Move{From: squareMust(4, 0), To: squareMust(6, 0)})
	if position.board.all != 0 {
		t.Fatal("want empty board")
	}
}

func TestValidatePromotion(t *testing.T) {
	testCases := []struct {
		name  string
		piece chessmodel.Piece
		move  chessmodel.Move
		want  error
	}{
		{
			name:  "bad type",
			piece: chessmodel.NewPiece(chessmodel.White, chessmodel.Pawn),
			move: chessmodel.Move{
				From:      squareMust(0, 6),
				To:        squareMust(0, 7),
				Promotion: chessmodel.PieceType(9),
			},
			want: chessmodel.ErrInvalidPromotion,
		},
		{
			name:  "rook move",
			piece: chessmodel.NewPiece(chessmodel.White, chessmodel.Rook),
			move: chessmodel.Move{
				From: squareMust(0, 0),
				To:   squareMust(0, 1),
			},
			want: nil,
		},
		{
			name:  "rook with promo",
			piece: chessmodel.NewPiece(chessmodel.White, chessmodel.Rook),
			move: chessmodel.Move{
				From:      squareMust(0, 0),
				To:        squareMust(0, 1),
				Promotion: chessmodel.Queen,
			},
			want: chessmodel.ErrInvalidPromotion,
		},
		{
			name:  "need promo",
			piece: chessmodel.NewPiece(chessmodel.White, chessmodel.Pawn),
			move: chessmodel.Move{
				From: squareMust(0, 6),
				To:   squareMust(0, 7),
			},
			want: chessmodel.ErrPromotionRequired,
		},
		{
			name:  "bad promo piece",
			piece: chessmodel.NewPiece(chessmodel.White, chessmodel.Pawn),
			move: chessmodel.Move{
				From:      squareMust(0, 6),
				To:        squareMust(0, 7),
				Promotion: chessmodel.King,
			},
			want: chessmodel.ErrInvalidPromotion,
		},
		{
			name:  "good promo",
			piece: chessmodel.NewPiece(chessmodel.White, chessmodel.Pawn),
			move: chessmodel.Move{
				From:      squareMust(0, 6),
				To:        squareMust(0, 7),
				Promotion: chessmodel.Queen,
			},
			want: nil,
		},
		{
			name:  "early promo",
			piece: chessmodel.NewPiece(chessmodel.White, chessmodel.Pawn),
			move: chessmodel.Move{
				From:      squareMust(0, 5),
				To:        squareMust(0, 6),
				Promotion: chessmodel.Queen,
			},
			want: chessmodel.ErrInvalidPromotion,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePromotion(tc.piece, tc.move)
			if !errors.Is(err, tc.want) {
				t.Fatalf("want %v, got %v", tc.want, err)
			}
		})
	}
}

func TestAttackKinds(t *testing.T) {
	testCases := []struct {
		name   string
		pos    Position
		target chessmodel.Square
		side   chessmodel.Side
		want   bool
	}{
		{
			name: "white pawn",
			pos: mustTestPosition(t,
				NewPositionBuilder().
					Place(squareMust(0, 0), chessmodel.White, chessmodel.King).
					Place(squareMust(7, 7), chessmodel.Black, chessmodel.King).
					Place(squareMust(3, 3), chessmodel.White, chessmodel.Pawn),
			),
			target: squareMust(4, 4),
			side:   chessmodel.White,
			want:   true,
		},
		{
			name: "black pawn",
			pos: mustTestPosition(t,
				NewPositionBuilder().
					WithSideToMove(chessmodel.Black).
					Place(squareMust(0, 0), chessmodel.White, chessmodel.King).
					Place(squareMust(7, 7), chessmodel.Black, chessmodel.King).
					Place(squareMust(3, 4), chessmodel.Black, chessmodel.Pawn),
			),
			target: squareMust(4, 3),
			side:   chessmodel.Black,
			want:   true,
		},
		{
			name: "knight",
			pos: mustTestPosition(t,
				NewPositionBuilder().
					Place(squareMust(0, 0), chessmodel.White, chessmodel.King).
					Place(squareMust(7, 7), chessmodel.Black, chessmodel.King).
					Place(squareMust(5, 3), chessmodel.White, chessmodel.Knight),
			),
			target: squareMust(4, 5),
			side:   chessmodel.White,
			want:   true,
		},
		{
			name: "rook",
			pos: mustTestPosition(t,
				NewPositionBuilder().
					Place(squareMust(0, 0), chessmodel.White, chessmodel.King).
					Place(squareMust(7, 7), chessmodel.Black, chessmodel.King).
					Place(squareMust(4, 0), chessmodel.White, chessmodel.Rook),
			),
			target: squareMust(4, 4),
			side:   chessmodel.White,
			want:   true,
		},
		{
			name: "bishop",
			pos: mustTestPosition(t,
				NewPositionBuilder().
					Place(squareMust(0, 0), chessmodel.White, chessmodel.King).
					Place(squareMust(7, 7), chessmodel.Black, chessmodel.King).
					Place(squareMust(1, 0), chessmodel.White, chessmodel.Bishop),
			),
			target: squareMust(4, 3),
			side:   chessmodel.White,
			want:   true,
		},
		{
			name: "queen",
			pos: mustTestPosition(t,
				NewPositionBuilder().
					Place(squareMust(0, 0), chessmodel.White, chessmodel.King).
					Place(squareMust(7, 7), chessmodel.Black, chessmodel.King).
					Place(squareMust(4, 0), chessmodel.White, chessmodel.Queen),
			),
			target: squareMust(4, 4),
			side:   chessmodel.White,
			want:   true,
		},
		{
			name: "king",
			pos: mustTestPosition(t,
				NewPositionBuilder().
					Place(squareMust(4, 3), chessmodel.White, chessmodel.King).
					Place(squareMust(7, 7), chessmodel.Black, chessmodel.King),
			),
			target: squareMust(4, 4),
			side:   chessmodel.White,
			want:   true,
		},
		{
			name: "blocked",
			pos: mustTestPosition(t,
				NewPositionBuilder().
					Place(squareMust(0, 0), chessmodel.White, chessmodel.King).
					Place(squareMust(7, 7), chessmodel.Black, chessmodel.King).
					Place(squareMust(4, 0), chessmodel.White, chessmodel.Rook).
					Place(squareMust(4, 2), chessmodel.Black, chessmodel.Pawn),
			),
			target: squareMust(4, 4),
			side:   chessmodel.White,
			want:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.pos.isSquareAttacked(tc.target, tc.side); got != tc.want {
				t.Fatalf("want %t, got %t", tc.want, got)
			}
		})
	}
}

func TestCastleHelpers(t *testing.T) {
	position := mustTestPosition(t,
		NewPositionBuilder().
			WithCastlingRights(chessmodel.WhiteKingSide).
			Place(squareMust(4, 0), chessmodel.White, chessmodel.King).
			Place(squareMust(7, 0), chessmodel.White, chessmodel.Rook).
			Place(squareMust(0, 7), chessmodel.Black, chessmodel.King).
			Place(squareMust(5, 7), chessmodel.Black, chessmodel.Rook),
	)

	if position.canCastleKingside(chessmodel.White) {
		t.Fatal("want bad castle")
	}

	position = mustTestPosition(t,
		NewPositionBuilder().
			Place(squareMust(4, 0), chessmodel.White, chessmodel.King).
			Place(squareMust(7, 0), chessmodel.White, chessmodel.Rook).
			Place(squareMust(4, 7), chessmodel.Black, chessmodel.King),
	)

	if position.canCastleKingside(chessmodel.White) {
		t.Fatal("want bad castle")
	}
}

func TestBoardPieceAt(t *testing.T) {
	board := newEmptyBoard()

	if _, ok := board.pieceAt(chessmodel.Square(99)); ok {
		t.Fatal("want no piece")
	}

	board.placePiece(chessmodel.NewPiece(chessmodel.White, chessmodel.Rook), squareMust(1, 1))
	piece, ok := board.pieceAt(squareMust(1, 1))
	if !ok {
		t.Fatal("want piece")
	}

	if piece.Type() != chessmodel.Rook || piece.Side() != chessmodel.White {
		t.Fatalf("want white rook, got %v", piece)
	}
}

func TestBoardPlacePiece(t *testing.T) {
	board := newEmptyBoard()

	board.placePiece(chessmodel.Piece{}, squareMust(1, 1))
	if _, ok := board.pieceAt(squareMust(1, 1)); ok {
		t.Fatal("want empty square")
	}

	board.placePiece(chessmodel.NewPiece(chessmodel.White, chessmodel.Rook), chessmodel.Square(99))
	if board.all != 0 {
		t.Fatal("want empty board")
	}
}

func TestBoardRemovePiece(t *testing.T) {
	board := newEmptyBoard()

	if _, ok := board.removePiece(chessmodel.Square(99)); ok {
		t.Fatal("want no piece")
	}

	if _, ok := board.removePiece(squareMust(0, 0)); ok {
		t.Fatal("want no piece")
	}

	board.placePiece(chessmodel.NewPiece(chessmodel.White, chessmodel.Rook), squareMust(0, 0))
	piece, ok := board.removePiece(squareMust(0, 0))
	if !ok {
		t.Fatal("want piece")
	}

	if piece.Type() != chessmodel.Rook || piece.Side() != chessmodel.White {
		t.Fatalf("want white rook, got %v", piece)
	}
}

func TestBoardMovePiece(t *testing.T) {
	board := newEmptyBoard()
	board.placePiece(chessmodel.NewPiece(chessmodel.White, chessmodel.King), squareMust(4, 0))
	board.placePiece(chessmodel.NewPiece(chessmodel.Black, chessmodel.Knight), squareMust(4, 7))

	moved, taken, ok := board.movePiece(squareMust(4, 0), squareMust(4, 7))
	if !ok {
		t.Fatal("want move")
	}

	if moved.Type() != chessmodel.King || moved.Side() != chessmodel.White {
		t.Fatalf("want white king, got %v", moved)
	}

	if taken.Type() != chessmodel.Knight || taken.Side() != chessmodel.Black {
		t.Fatalf("want black knight, got %v", taken)
	}

	if _, _, ok := board.movePiece(squareMust(4, 0), squareMust(4, 1)); ok {
		t.Fatal("want bad move")
	}
}

func TestBoardKingSquare(t *testing.T) {
	board := newEmptyBoard()
	if _, ok := board.kingSquare(chessmodel.Black); ok {
		t.Fatal("want no king")
	}

	board.placePiece(chessmodel.NewPiece(chessmodel.White, chessmodel.King), squareMust(4, 7))
	square, ok := board.kingSquare(chessmodel.White)
	if !ok || square != squareMust(4, 7) {
		t.Fatalf("want e8, got %v, %t", square, ok)
	}
}

func TestFirstSquare(t *testing.T) {
	if _, ok := firstSquare(0); ok {
		t.Fatal("want no square")
	}

	square, ok := firstSquare(squareBitboard(squareMust(3, 3)))
	if !ok || square != squareMust(3, 3) {
		t.Fatalf("want d4, got %v, %t", square, ok)
	}

	square, ok = firstSquare(squareBitboard(squareMust(3, 3)) | squareBitboard(squareMust(4, 4)))
	if !ok || square != squareMust(3, 3) {
		t.Fatalf("want d4, got %v, %t", square, ok)
	}
}

func TestDrawHelpers(t *testing.T) {
	position := mustTestPosition(t,
		NewPositionBuilder().
			Place(squareMust(0, 0), chessmodel.White, chessmodel.King).
			Place(squareMust(7, 7), chessmodel.Black, chessmodel.King),
	)

	if position.bishopsShareColorComplex() {
		t.Fatal("want false")
	}

	position = mustTestPosition(t,
		NewPositionBuilder().
			Place(squareMust(0, 0), chessmodel.White, chessmodel.King).
			Place(squareMust(3, 3), chessmodel.White, chessmodel.Queen).
			Place(squareMust(7, 7), chessmodel.Black, chessmodel.King),
	)

	if HasInsufficientMaterial(position) {
		t.Fatal("want false")
	}

	position = mustTestPosition(t,
		NewPositionBuilder().
			Place(squareMust(0, 0), chessmodel.White, chessmodel.King).
			Place(squareMust(1, 0), chessmodel.White, chessmodel.Bishop).
			Place(squareMust(2, 0), chessmodel.White, chessmodel.Knight).
			Place(squareMust(3, 0), chessmodel.White, chessmodel.Knight).
			Place(squareMust(7, 7), chessmodel.Black, chessmodel.King),
	)

	if HasInsufficientMaterial(position) {
		t.Fatal("want false")
	}
}

func TestHelpers(t *testing.T) {
	empty := noSquare()
	if empty.ok {
		t.Fatal("want no square")
	}

	filled := someSquare(squareMust(3, 3))
	if !filled.ok || filled.value != squareMust(3, 3) {
		t.Fatal("want d4")
	}

	if squareIsValid(chessmodel.Square(99)) {
		t.Fatal("want bad square")
	}

	if !squareIsValid(squareMust(0, 0)) {
		t.Fatal("want good square")
	}

	if got := squareBitboard(squareMust(3, 3)); got != uint64(1)<<27 {
		t.Fatalf("want bit 27, got %d", got)
	}

	if _, ok := squareOffset(chessmodel.Square(99), 1, 1); ok {
		t.Fatal("want bad source")
	}

	if _, ok := squareOffset(squareMust(7, 7), 1, 0); ok {
		t.Fatal("want bad file")
	}

	if _, ok := squareOffset(squareMust(0, 0), 0, -1); ok {
		t.Fatal("want bad rank")
	}

	offset, ok := squareOffset(squareMust(3, 3), 1, 1)
	if !ok || offset != squareMust(4, 4) {
		t.Fatalf("want e5, got %v, %t", offset, ok)
	}

	if got := squareColor(squareMust(0, 0)); got != 0 {
		t.Fatalf("want 0, got %d", got)
	}

	if got := squareColor(squareMust(0, 1)); got != 1 {
		t.Fatalf("want 1, got %d", got)
	}

	if !isValidSide(chessmodel.White) {
		t.Fatal("want good side")
	}

	if isValidSide(chessmodel.Side(9)) {
		t.Fatal("want bad side")
	}

	if got := sideIndex(chessmodel.Black); got != 1 {
		t.Fatalf("want 1, got %d", got)
	}

	if got := sideIndex(chessmodel.White); got != 0 {
		t.Fatalf("want 0, got %d", got)
	}

	if got := bitboardIndex(chessmodel.Pawn); got != 0 {
		t.Fatalf("want 0, got %d", got)
	}

	if got := bitboardIndex(chessmodel.King); got != 5 {
		t.Fatalf("want 5, got %d", got)
	}

	if !isPromotionChoice(chessmodel.Queen) {
		t.Fatal("want good promo")
	}

	if isPromotionChoice(chessmodel.King) {
		t.Fatal("want bad promo")
	}

	goodMove := chessmodel.Move{
		From: squareMust(0, 1),
		To:   squareMust(0, 2),
	}

	if err := validateMoveSquares(goodMove); err != nil {
		t.Fatalf("want good move, got %v", err)
	}

	badMove := chessmodel.Move{
		From: squareMust(0, 1),
		To:   chessmodel.Square(99),
	}

	if err := validateMoveSquares(badMove); !errors.Is(err, chessmodel.ErrInvalidSquare) {
		t.Fatalf("want bad square, got %v", err)
	}

	if got := absInt(3); got != 3 {
		t.Fatalf("want 3, got %d", got)
	}

	if got := absInt(-3); got != 3 {
		t.Fatalf("want 3, got %d", got)
	}
}

func TestSquareMustPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("want panic")
		}
	}()

	_ = squareMust(9, 9)
}
