package model

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

func TestMoveValidateSquares(t *testing.T) {
	bad := Move{
		From: mustSquare(0, 1),
		To:   Square(99),
	}

	if err := bad.validateSquares(); !errors.Is(err, ErrInvalidSquare) {
		t.Fatal("want bad square")
	}

	move := Move{
		From: mustSquare(0, 1),
		To:   mustSquare(0, 2),
	}

	if err := move.validateSquares(); err != nil {
		t.Fatalf("want good move, got %v", err)
	}
}

func TestSquareHelpers(t *testing.T) {
	if got := mustSquare(1, 1); got != Square(9) {
		t.Fatalf("want 9, got %d", got)
	}

	if _, ok := Square(99).offset(1, 1); ok {
		t.Fatal("want bad offset")
	}

	if _, ok := mustSquare(7, 7).offset(1, 0); ok {
		t.Fatal("want bad offset")
	}

	if _, ok := mustSquare(0, 0).offset(0, -1); ok {
		t.Fatal("want bad rank")
	}

	square, ok := mustSquare(3, 3).offset(1, 1)
	if !ok || square != mustSquare(4, 4) {
		t.Fatalf("want e5, got %v, %t", square, ok)
	}

	if got := mustSquare(3, 3).bitboard(); got != uint64(1)<<27 {
		t.Fatalf("want bit 27, got %d", got)
	}

	if got := mustSquare(0, 0).color(); got != 0 {
		t.Fatalf("want 0, got %d", got)
	}

	if got := mustSquare(0, 1).color(); got != 1 {
		t.Fatalf("want 1, got %d", got)
	}
}

func TestMustSquarePanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("want panic")
		}
	}()

	_ = mustSquare(9, 9)
}
