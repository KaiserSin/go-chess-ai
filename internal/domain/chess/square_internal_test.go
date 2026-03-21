package chess

import "testing"

func TestSquareOffset(t *testing.T) {
	if _, ok := Square(99).offset(1, 1); ok {
		t.Fatal("want bad offset")
	}

	if _, ok := mustSquare(7, 7).offset(1, 0); ok {
		t.Fatal("want bad offset")
	}

	square, ok := mustSquare(3, 3).offset(1, 1)
	if !ok || square != mustSquare(4, 4) {
		t.Fatalf("want e5, got %v, %t", square, ok)
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
