package chess

import (
	"errors"
	"testing"
)

func TestPrivateTypes(t *testing.T) {
	if !NoPieceType.isValid() {
		t.Fatal("want good type")
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

func TestValidatePromotion(t *testing.T) {
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

func mustTestPosition(t *testing.T, builder *PositionBuilder) Position {
	t.Helper()

	position, err := builder.Build()
	if err != nil {
		t.Fatalf("want good position, got %v", err)
	}

	return position
}
