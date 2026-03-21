package chess_test

import (
	"errors"
	"testing"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestThreefoldDraw(t *testing.T) {
	game := chess.NewGame()
	cycle := []chess.Move{
		mustMove(t, "g1", "f3"),
		mustMove(t, "g8", "f6"),
		mustMove(t, "f3", "g1"),
		mustMove(t, "f6", "g8"),
	}

	applyMoves(t, game, cycle...)
	if game.IsFinished() {
		t.Fatal("want no draw")
	}

	applyMoves(t, game, cycle[:2]...)
	if game.IsFinished() {
		t.Fatal("want no draw")
	}

	applyMoves(t, game, cycle[2:]...)
	if !game.IsFinished() {
		t.Fatal("want draw")
	}

	if game.Outcome().Reason() != chess.OutcomeByThreefoldRepetition {
		t.Fatalf("want threefold, got %s", game.Outcome().Reason())
	}

	if moves := game.LegalMoves(); moves != nil {
		t.Fatal("want no moves")
	}

	if err := game.ApplyMove(mustMove(t, "e2", "e4")); !errors.Is(err, chess.ErrGameFinished) {
		t.Fatalf("want game over error, got %v", err)
	}
}

func TestRepetitionIgnoresClocks(t *testing.T) {
	game := mustBuildGame(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			WithHalfmoveClock(20).
			WithFullmoveNumber(7).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "g1"), chess.White, chess.Knight).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "g8"), chess.Black, chess.Knight),
	)

	cycle := []chess.Move{
		mustMove(t, "g1", "f3"),
		mustMove(t, "g8", "f6"),
		mustMove(t, "f3", "g1"),
		mustMove(t, "f6", "g8"),
	}

	applyMoves(t, game, cycle...)
	if game.IsFinished() {
		t.Fatal("want no draw")
	}

	applyMoves(t, game, cycle...)
	if game.Outcome().Reason() != chess.OutcomeByThreefoldRepetition {
		t.Fatalf("want threefold, got %s", game.Outcome().Reason())
	}
}

func TestRepetitionIgnoresBadEnPassant(t *testing.T) {
	game := mustBuildGame(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.Black).
			WithEnPassantSquare(mustParseSquare(t, "e3")).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "g1"), chess.White, chess.Knight).
			Place(mustParseSquare(t, "e4"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "g8"), chess.Black, chess.Knight),
	)

	cycle := []chess.Move{
		mustMove(t, "g8", "f6"),
		mustMove(t, "g1", "f3"),
		mustMove(t, "f6", "g8"),
		mustMove(t, "f3", "g1"),
	}

	applyMoves(t, game, cycle...)
	if game.IsFinished() {
		t.Fatal("want no draw")
	}

	applyMoves(t, game, cycle...)
	if game.Outcome().Reason() != chess.OutcomeByThreefoldRepetition {
		t.Fatalf("want threefold, got %s", game.Outcome().Reason())
	}
}

func TestRepetitionKeepsGoodEnPassant(t *testing.T) {
	game := mustBuildGame(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			WithEnPassantSquare(mustParseSquare(t, "d6")).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "g1"), chess.White, chess.Knight).
			Place(mustParseSquare(t, "e5"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "g8"), chess.Black, chess.Knight).
			Place(mustParseSquare(t, "d5"), chess.Black, chess.Pawn),
	)

	cycle := []chess.Move{
		mustMove(t, "g1", "f3"),
		mustMove(t, "g8", "f6"),
		mustMove(t, "f3", "g1"),
		mustMove(t, "f6", "g8"),
	}

	applyMoves(t, game, cycle...)
	applyMoves(t, game, cycle...)

	if game.IsFinished() {
		t.Fatal("want game on")
	}

	if game.Outcome().Reason() != chess.NoOutcomeReason {
		t.Fatalf("want no outcome, got %s", game.Outcome().Reason())
	}
}
