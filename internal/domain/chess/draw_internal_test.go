package chess

import "testing"

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

func TestDrawHelpers(t *testing.T) {
	position := mustTestPosition(t,
		NewPositionBuilder().
			Place(mustSquare(0, 0), White, King).
			Place(mustSquare(7, 7), Black, King),
	)

	if position.bishopsShareColorComplex() {
		t.Fatal("want false")
	}

	position = mustTestPosition(t,
		NewPositionBuilder().
			Place(mustSquare(0, 0), White, King).
			Place(mustSquare(3, 3), White, Queen).
			Place(mustSquare(7, 7), Black, King),
	)

	if position.hasInsufficientMaterial() {
		t.Fatal("want false")
	}

	position = mustTestPosition(t,
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
