package chess_test

import (
	"errors"
	"testing"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestBadBuilder(t *testing.T) {
	testCases := []struct {
		name    string
		builder *chess.PositionBuilder
		want    error
	}{
		{
			name:    "bad side",
			builder: chess.NewPositionBuilder().WithSideToMove(chess.Side(9)),
			want:    chess.ErrInvalidPosition,
		},
		{
			name:    "bad en passant",
			builder: chess.NewPositionBuilder().WithEnPassantSquare(chess.Square(99)),
			want:    chess.ErrInvalidSquare,
		},
		{
			name:    "bad halfmove",
			builder: chess.NewPositionBuilder().WithHalfmoveClock(-1),
			want:    chess.ErrInvalidPosition,
		},
		{
			name:    "bad fullmove",
			builder: chess.NewPositionBuilder().WithFullmoveNumber(0),
			want:    chess.ErrInvalidPosition,
		},
		{
			name: "bad piece",
			builder: chess.NewPositionBuilder().
				Place(mustParseSquare(t, "e4"), chess.White, chess.PieceType(9)),
			want: chess.ErrInvalidPosition,
		},
		{
			name:    "bad place square",
			builder: chess.NewPositionBuilder().Place(chess.Square(99), chess.White, chess.Pawn),
			want:    chess.ErrInvalidSquare,
		},
		{
			name: "missing king",
			builder: chess.NewPositionBuilder().
				Place(mustParseSquare(t, "e1"), chess.White, chess.King),
			want: chess.ErrInvalidPosition,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.builder.Build()
			if !errors.Is(err, tc.want) {
				t.Fatalf("want %v, got %v", tc.want, err)
			}
		})
	}
}

func TestBuilderStops(t *testing.T) {
	builder := chess.NewPositionBuilder().
		WithSideToMove(chess.Side(9)).
		WithSideToMove(chess.Black).
		WithCastlingRights(chess.WhiteKingSide).
		WithEnPassantSquare(mustParseSquare(t, "e3")).
		WithHalfmoveClock(4).
		WithFullmoveNumber(2).
		Place(mustParseSquare(t, "e1"), chess.White, chess.King)

	_, err := builder.Build()
	if !errors.Is(err, chess.ErrInvalidPosition) {
		t.Fatalf("want bad position error, got %v", err)
	}
}

func TestBadGamePosition(t *testing.T) {
	if _, err := chess.NewGameFromPosition(chess.Position{}); !errors.Is(err, chess.ErrInvalidPosition) {
		t.Fatalf("want bad position error, got %v", err)
	}

	if (chess.Position{}).IsInCheck(chess.White) {
		t.Fatal("want no check")
	}
}
