package chess

import "testing"

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
			pos: mustTestPosition(t,
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
			pos: mustTestPosition(t,
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
			pos: mustTestPosition(t,
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
			pos: mustTestPosition(t,
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
			pos: mustTestPosition(t,
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
			pos: mustTestPosition(t,
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
			pos: mustTestPosition(t,
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
			pos: mustTestPosition(t,
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

func TestCastleHelpers(t *testing.T) {
	position := mustTestPosition(t,
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

	position = mustTestPosition(t,
		NewPositionBuilder().
			Place(mustSquare(4, 0), White, King).
			Place(mustSquare(7, 0), White, Rook).
			Place(mustSquare(4, 7), Black, King),
	)

	if position.canCastleKingside(White) {
		t.Fatal("want bad castle")
	}
}
