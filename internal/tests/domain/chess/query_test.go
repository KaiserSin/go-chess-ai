package chess_test

import (
	"testing"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestHasInsufficientMaterial(t *testing.T) {
	lowMaterial := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	if !chess.HasInsufficientMaterial(lowMaterial) {
		t.Fatal("want insufficient material")
	}

	normalMaterial := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	if chess.HasInsufficientMaterial(normalMaterial) {
		t.Fatal("did not expect insufficient material")
	}
}

func TestIsFiftyMoveDraw(t *testing.T) {
	drawPosition := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			WithHalfmoveClock(100).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	if !chess.IsFiftyMoveDraw(drawPosition) {
		t.Fatal("want fifty-move draw")
	}

	normalPosition := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			WithHalfmoveClock(99).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	if chess.IsFiftyMoveDraw(normalPosition) {
		t.Fatal("did not expect fifty-move draw")
	}
}

func TestNewRepetitionKeyIgnoresClocks(t *testing.T) {
	first := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			WithHalfmoveClock(4).
			WithFullmoveNumber(7).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "g1"), chess.White, chess.Knight).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "g8"), chess.Black, chess.Knight),
	)
	second := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			WithHalfmoveClock(80).
			WithFullmoveNumber(22).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "g1"), chess.White, chess.Knight).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "g8"), chess.Black, chess.Knight),
	)

	firstKey := chess.NewRepetitionKey(first)
	secondKey := chess.NewRepetitionKey(second)
	if firstKey != secondKey {
		t.Fatal("want repetition key to ignore clocks")
	}
}
