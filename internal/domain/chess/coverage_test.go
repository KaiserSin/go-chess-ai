package chess

import (
	"errors"
	"testing"
)

func TestBoardHelpers(t *testing.T) {
	board := newEmptyBoard()

	if _, ok := board.pieceAt(Square(99)); ok {
		t.Fatal("want no piece")
	}

	if _, ok := board.removePiece(Square(99)); ok {
		t.Fatal("want no piece")
	}

	if _, ok := board.removePiece(mustSquare(0, 0)); ok {
		t.Fatal("want no piece")
	}

	board.placePiece(Piece{}, mustSquare(1, 1))
	if _, ok := board.pieceAt(mustSquare(1, 1)); ok {
		t.Fatal("want empty square")
	}

	board.placePiece(newPiece(White, Rook), Square(99))
	if board.all != 0 {
		t.Fatal("want empty board")
	}

	board.placePiece(newPiece(White, King), mustSquare(4, 0))
	board.placePiece(newPiece(Black, Knight), mustSquare(4, 7))

	moved, taken, ok := board.movePiece(mustSquare(4, 0), mustSquare(4, 7))
	if !ok {
		t.Fatal("want move")
	}

	if moved.kind != King || moved.side != White {
		t.Fatalf("bad moved piece: %v", moved)
	}

	if taken.kind != Knight || taken.side != Black {
		t.Fatalf("bad taken piece: %v", taken)
	}

	if _, _, ok := board.movePiece(mustSquare(4, 0), mustSquare(4, 1)); ok {
		t.Fatal("want bad move")
	}

	kingSquare, ok := board.kingSquare(White)
	if !ok || kingSquare != mustSquare(4, 7) {
		t.Fatalf("want king on e8, got %v, %t", kingSquare, ok)
	}

	emptyBoard := newEmptyBoard()
	if _, ok := emptyBoard.kingSquare(Black); ok {
		t.Fatal("want no king")
	}
}

func TestPrivateTypes(t *testing.T) {
	if !NoPieceType.isValid() {
		t.Fatal("want valid none type")
	}

	if PieceType(9).isValid() {
		t.Fatal("want bad type")
	}

	testCases := []struct {
		name string
		kind PieceType
		want string
	}{
		{name: "pawn", kind: Pawn, want: "p"},
		{name: "knight", kind: Knight, want: "n"},
		{name: "bishop", kind: Bishop, want: "b"},
		{name: "rook", kind: Rook, want: "r"},
		{name: "queen", kind: Queen, want: "q"},
		{name: "king", kind: King, want: "k"},
		{name: "bad", kind: NoPieceType, want: "?"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.kind.symbol(); got != tc.want {
				t.Fatalf("want %q, got %q", tc.want, got)
			}
		})
	}
}

func TestSquareHelpers(t *testing.T) {
	if _, ok := Square(99).offset(1, 1); ok {
		t.Fatal("want bad offset")
	}

	if _, ok := mustSquare(7, 7).offset(1, 0); ok {
		t.Fatal("want bad offset")
	}

	defer func() {
		if recover() == nil {
			t.Fatal("want panic")
		}
	}()

	_ = mustSquare(9, 9)
}

func TestFirstSquare(t *testing.T) {
	if _, ok := firstSquare(0); ok {
		t.Fatal("want no square")
	}

	square, ok := firstSquare(mustSquare(3, 3).bitboard())
	if !ok || square != mustSquare(3, 3) {
		t.Fatalf("want d4, got %v, %t", square, ok)
	}

	square, ok = firstSquare(mustSquare(3, 3).bitboard() | mustSquare(4, 4).bitboard())
	if !ok || square != mustSquare(3, 3) {
		t.Fatalf("want d4, got %v, %t", square, ok)
	}
}

func TestMoveCastlingRook(t *testing.T) {
	position := Position{board: newEmptyBoard()}
	position.moveCastlingRook(White, Move{From: mustSquare(4, 0), To: mustSquare(4, 1)})
	if position.board.all != 0 {
		t.Fatal("want empty board")
	}

	position = Position{board: newEmptyBoard()}
	position.moveCastlingRook(White, Move{From: mustSquare(4, 0), To: mustSquare(6, 0)})
	if position.board.all != 0 {
		t.Fatal("want empty board")
	}
}

func TestAttackKinds(t *testing.T) {
	testCases := []struct {
		name   string
		pos    Position
		target Square
		side   Side
		want   bool
	}{
		{
			name: "white pawn",
			pos: testPosition(t,
				NewPositionBuilder().
					Place(mustSquare(0, 0), White, King).
					Place(mustSquare(7, 7), Black, King).
					Place(mustSquare(3, 3), White, Pawn),
			),
			target: mustSquare(4, 4),
			side:   White,
			want:   true,
		},
		{
			name: "black pawn",
			pos: testPosition(t,
				NewPositionBuilder().
					WithSideToMove(Black).
					Place(mustSquare(0, 0), White, King).
					Place(mustSquare(7, 7), Black, King).
					Place(mustSquare(3, 4), Black, Pawn),
			),
			target: mustSquare(4, 3),
			side:   Black,
			want:   true,
		},
		{
			name: "knight",
			pos: testPosition(t,
				NewPositionBuilder().
					Place(mustSquare(0, 0), White, King).
					Place(mustSquare(7, 7), Black, King).
					Place(mustSquare(5, 3), White, Knight),
			),
			target: mustSquare(4, 5),
			side:   White,
			want:   true,
		},
		{
			name: "rook",
			pos: testPosition(t,
				NewPositionBuilder().
					Place(mustSquare(0, 0), White, King).
					Place(mustSquare(7, 7), Black, King).
					Place(mustSquare(4, 0), White, Rook),
			),
			target: mustSquare(4, 4),
			side:   White,
			want:   true,
		},
		{
			name: "bishop",
			pos: testPosition(t,
				NewPositionBuilder().
					Place(mustSquare(0, 0), White, King).
					Place(mustSquare(7, 7), Black, King).
					Place(mustSquare(1, 0), White, Bishop),
			),
			target: mustSquare(4, 3),
			side:   White,
			want:   true,
		},
		{
			name: "queen",
			pos: testPosition(t,
				NewPositionBuilder().
					Place(mustSquare(0, 0), White, King).
					Place(mustSquare(7, 7), Black, King).
					Place(mustSquare(4, 0), White, Queen),
			),
			target: mustSquare(4, 4),
			side:   White,
			want:   true,
		},
		{
			name: "king",
			pos: testPosition(t,
				NewPositionBuilder().
					Place(mustSquare(4, 3), White, King).
					Place(mustSquare(7, 7), Black, King),
			),
			target: mustSquare(4, 4),
			side:   White,
			want:   true,
		},
		{
			name: "blocked",
			pos: testPosition(t,
				NewPositionBuilder().
					Place(mustSquare(0, 0), White, King).
					Place(mustSquare(7, 7), Black, King).
					Place(mustSquare(4, 0), White, Rook).
					Place(mustSquare(4, 2), Black, Pawn),
			),
			target: mustSquare(4, 4),
			side:   White,
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

func TestDrawHelpers(t *testing.T) {
	position := testPosition(t,
		NewPositionBuilder().
			Place(mustSquare(0, 0), White, King).
			Place(mustSquare(7, 7), Black, King),
	)

	if position.bishopsShareColorComplex() {
		t.Fatal("want false")
	}

	position = testPosition(t,
		NewPositionBuilder().
			Place(mustSquare(0, 0), White, King).
			Place(mustSquare(3, 3), White, Queen).
			Place(mustSquare(7, 7), Black, King),
	)

	if position.hasInsufficientMaterial() {
		t.Fatal("want false")
	}

	position = testPosition(t,
		NewPositionBuilder().
			Place(mustSquare(0, 0), White, King).
			Place(mustSquare(1, 0), White, Bishop).
			Place(mustSquare(2, 0), White, Knight).
			Place(mustSquare(3, 0), White, Knight).
			Place(mustSquare(7, 7), Black, King),
	)

	if position.hasInsufficientMaterial() {
		t.Fatal("want false")
	}
}

func TestPositionValidate(t *testing.T) {
	position := Position{
		board:      newEmptyBoard(),
		sideToMove: Side(9),
	}

	if !errors.Is(position.validate(), ErrInvalidPosition) {
		t.Fatal("want bad position")
	}

	position = Position{
		board:      newEmptyBoard(),
		sideToMove: White,
	}
	position.board.placePiece(newPiece(Black, King), mustSquare(4, 7))
	if !errors.Is(position.validate(), ErrInvalidPosition) {
		t.Fatal("want bad position")
	}

	position = Position{
		board:      newEmptyBoard(),
		sideToMove: White,
	}
	position.board.placePiece(newPiece(White, King), mustSquare(4, 0))
	if !errors.Is(position.validate(), ErrInvalidPosition) {
		t.Fatal("want bad position")
	}
}

func TestCastleRuleHelpers(t *testing.T) {
	position := testPosition(t,
		NewPositionBuilder().
			WithCastlingRights(WhiteKingSide).
			Place(mustSquare(4, 0), White, King).
			Place(mustSquare(7, 0), White, Rook).
			Place(mustSquare(0, 7), Black, King).
			Place(mustSquare(5, 7), Black, Rook),
	)

	if position.canCastleKingside(White) {
		t.Fatal("want bad castle")
	}

	position = testPosition(t,
		NewPositionBuilder().
			Place(mustSquare(4, 0), White, King).
			Place(mustSquare(7, 0), White, Rook).
			Place(mustSquare(4, 7), Black, King),
	)

	if position.canCastleKingside(White) {
		t.Fatal("want bad castle")
	}
}

func TestValidatePromotionRules(t *testing.T) {
	testCases := []struct {
		name  string
		piece Piece
		move  Move
		want  error
	}{
		{
			name:  "bad type",
			piece: newPiece(White, Pawn),
			move: Move{
				From:      mustSquare(0, 6),
				To:        mustSquare(0, 7),
				Promotion: PieceType(9),
			},
			want: ErrInvalidPromotion,
		},
		{
			name:  "rook move",
			piece: newPiece(White, Rook),
			move: Move{
				From: mustSquare(0, 0),
				To:   mustSquare(0, 1),
			},
			want: nil,
		},
		{
			name:  "rook with promo",
			piece: newPiece(White, Rook),
			move: Move{
				From:      mustSquare(0, 0),
				To:        mustSquare(0, 1),
				Promotion: Queen,
			},
			want: ErrInvalidPromotion,
		},
		{
			name:  "need promo",
			piece: newPiece(White, Pawn),
			move: Move{
				From: mustSquare(0, 6),
				To:   mustSquare(0, 7),
			},
			want: ErrPromotionRequired,
		},
		{
			name:  "bad promo piece",
			piece: newPiece(White, Pawn),
			move: Move{
				From:      mustSquare(0, 6),
				To:        mustSquare(0, 7),
				Promotion: King,
			},
			want: ErrInvalidPromotion,
		},
		{
			name:  "good promo",
			piece: newPiece(White, Pawn),
			move: Move{
				From:      mustSquare(0, 6),
				To:        mustSquare(0, 7),
				Promotion: Queen,
			},
			want: nil,
		},
		{
			name:  "early promo",
			piece: newPiece(White, Pawn),
			move: Move{
				From:      mustSquare(0, 5),
				To:        mustSquare(0, 6),
				Promotion: Queen,
			},
			want: ErrInvalidPromotion,
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

func testPosition(t *testing.T, builder *PositionBuilder) Position {
	t.Helper()

	position, err := builder.Build()
	if err != nil {
		t.Fatalf("want good position, got %v", err)
	}

	return position
}
