package scene

import "testing"

func TestSpritePlacementForSquare(t *testing.T) {
	placement := spritePlacementForRect(40, 88, 80, 128, 128)

	if placement.X != 40 {
		t.Fatalf("want x 40, got %v", placement.X)
	}

	if placement.Y != 88 {
		t.Fatalf("want y 88, got %v", placement.Y)
	}

	if placement.ScaleX != 0.625 {
		t.Fatalf("want scaleX 0.625, got %v", placement.ScaleX)
	}

	if placement.ScaleY != 0.625 {
		t.Fatalf("want scaleY 0.625, got %v", placement.ScaleY)
	}
}
