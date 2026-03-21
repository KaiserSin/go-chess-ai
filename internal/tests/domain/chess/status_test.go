package chess_test

import (
	"errors"
	"testing"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestCheck(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.Black).
			Place(mustParseSquare(t, "h1"), chess.White, chess.King).
			Place(mustParseSquare(t, "e1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	if position.Status() != chess.Check {
		t.Fatalf("want check, got %s", position.Status())
	}

	if !position.IsInCheck(chess.Black) {
		t.Fatal("want black in check")
	}

	if len(position.LegalMoves()) == 0 {
		t.Fatal("want move")
	}
}

func TestCheckmate(t *testing.T) {
	game := chess.NewGame()
	applyMoves(t, game,
		mustMove(t, "f2", "f3"),
		mustMove(t, "e7", "e5"),
		mustMove(t, "g2", "g4"),
		mustMove(t, "d8", "h4"),
	)

	if game.Status() != chess.Checkmate {
		t.Fatalf("want checkmate, got %s", game.Status())
	}

	if !game.IsFinished() {
		t.Fatal("want game over")
	}

	if !game.Outcome().IsDecisive() {
		t.Fatal("want win")
	}

	if game.Outcome().Reason() != chess.OutcomeByCheckmate {
		t.Fatalf("want checkmate, got %s", game.Outcome().Reason())
	}

	winner, ok := game.Outcome().Winner()
	if !ok || winner != chess.Black {
		t.Fatalf("want black win, got %v, %t", winner, ok)
	}

	if moves := game.LegalMoves(); moves != nil {
		t.Fatal("want no moves")
	}

	if err := game.ApplyMove(mustMove(t, "g2", "g3")); !errors.Is(err, chess.ErrGameFinished) {
		t.Fatalf("want game over error, got %v", err)
	}
}

func TestStalemate(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.Black).
			Place(mustParseSquare(t, "c6"), chess.White, chess.King).
			Place(mustParseSquare(t, "b6"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.King),
	)

	if position.Status() != chess.Stalemate {
		t.Fatalf("want stalemate, got %s", position.Status())
	}

	_, err := position.ApplyMove(mustMove(t, "a8", "a7"))
	if !errors.Is(err, chess.ErrGameFinished) {
		t.Fatalf("want game over error, got %v", err)
	}

	game := mustNewGameFromPosition(t, position)
	if !game.IsFinished() {
		t.Fatal("want game over")
	}

	if game.Status() != chess.Stalemate {
		t.Fatalf("want stalemate, got %s", game.Status())
	}

	if game.Outcome().Reason() != chess.OutcomeByStalemate {
		t.Fatalf("want stalemate, got %s", game.Outcome().Reason())
	}

	if !game.Outcome().IsDraw() {
		t.Fatal("want draw")
	}

	if game.Outcome().IsDecisive() {
		t.Fatal("want no win")
	}

	if _, ok := game.Outcome().Winner(); ok {
		t.Fatal("want no winner")
	}

	if moves := game.LegalMoves(); moves != nil {
		t.Fatal("want no moves")
	}

	if err := game.ApplyMove(mustMove(t, "a8", "a7")); !errors.Is(err, chess.ErrGameFinished) {
		t.Fatalf("want game over error, got %v", err)
	}
}
