package chess_test

import (
	"errors"
	"testing"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestStalemateDraw(t *testing.T) {
	game := mustBuildGame(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.Black).
			Place(mustParseSquare(t, "c6"), chess.White, chess.King).
			Place(mustParseSquare(t, "b6"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.King),
	)

	if !game.IsFinished() {
		t.Fatal("want finished game")
	}

	if game.Status() != chess.Stalemate {
		t.Fatalf("expected stalemate status, got %s", game.Status())
	}

	if game.Outcome().Reason() != chess.OutcomeByStalemate {
		t.Fatalf("expected stalemate outcome, got %s", game.Outcome().Reason())
	}

	if !game.Outcome().IsDraw() {
		t.Fatal("want draw")
	}

	if _, ok := game.Outcome().Winner(); ok {
		t.Fatal("want no winner")
	}

	if len(game.LegalMoves()) != 0 {
		t.Fatalf("expected no legal moves, got %d", len(game.LegalMoves()))
	}
}

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
		t.Fatal("want no draw yet")
	}

	applyMoves(t, game, cycle[:2]...)
	if game.IsFinished() {
		t.Fatal("want no draw yet")
	}

	applyMoves(t, game, cycle[2:]...)
	if !game.IsFinished() {
		t.Fatal("want draw by repetition")
	}

	if game.Outcome().Reason() != chess.OutcomeByThreefoldRepetition {
		t.Fatalf("expected threefold repetition outcome, got %s", game.Outcome().Reason())
	}

	if len(game.LegalMoves()) != 0 {
		t.Fatalf("want 0 legal moves, got %d", len(game.LegalMoves()))
	}

	if err := game.ApplyMove(mustMove(t, "e2", "e4")); !errors.Is(err, chess.ErrGameFinished) {
		t.Fatalf("expected ErrGameFinished after repetition draw, got %v", err)
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
		t.Fatal("want no draw yet")
	}

	applyMoves(t, game, cycle...)
	if game.Outcome().Reason() != chess.OutcomeByThreefoldRepetition {
		t.Fatalf("expected repetition to ignore clocks, got %s", game.Outcome().Reason())
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
		t.Fatal("want no draw yet")
	}

	applyMoves(t, game, cycle...)
	if game.Outcome().Reason() != chess.OutcomeByThreefoldRepetition {
		t.Fatalf("expected non capturable en passant to be ignored, got %s", game.Outcome().Reason())
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
		t.Fatal("want no draw here")
	}

	if game.Outcome().Reason() != chess.NoOutcomeReason {
		t.Fatalf("expected no outcome yet, got %s", game.Outcome().Reason())
	}
}

func TestFiftyMoveDraw(t *testing.T) {
	game := mustBuildGame(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			WithHalfmoveClock(99).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "h8"), chess.Black, chess.Rook),
	)

	if err := game.ApplyMove(mustMove(t, "a1", "a2")); err != nil {
		t.Fatalf("ApplyMove error: %v", err)
	}

	if !game.IsFinished() {
		t.Fatal("want draw by 50-move rule")
	}

	if game.Position().HalfmoveClock() != 100 {
		t.Fatalf("expected halfmove clock 100, got %d", game.Position().HalfmoveClock())
	}

	if game.Outcome().Reason() != chess.OutcomeByFiftyMoveRule {
		t.Fatalf("expected fifty-move outcome, got %s", game.Outcome().Reason())
	}
}

func TestClockResetOnPawnMove(t *testing.T) {
	game := mustBuildGame(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			WithHalfmoveClock(99).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "a2"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "h8"), chess.Black, chess.Rook),
	)

	if err := game.ApplyMove(mustMove(t, "a2", "a3")); err != nil {
		t.Fatalf("ApplyMove error: %v", err)
	}

	if game.Position().HalfmoveClock() != 0 {
		t.Fatalf("expected halfmove clock reset to 0, got %d", game.Position().HalfmoveClock())
	}

	if game.IsFinished() {
		t.Fatal("want no draw")
	}
}

func TestClockResetOnCapture(t *testing.T) {
	game := mustBuildGame(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			WithHalfmoveClock(99).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook),
	)

	if err := game.ApplyMove(mustMove(t, "a1", "a8")); err != nil {
		t.Fatalf("ApplyMove error: %v", err)
	}

	if game.Position().HalfmoveClock() != 0 {
		t.Fatalf("expected halfmove clock reset to 0, got %d", game.Position().HalfmoveClock())
	}

	if game.IsFinished() {
		t.Fatal("want no draw")
	}
}

func TestLowMaterialDraw(t *testing.T) {
	testCases := []struct {
		name    string
		builder *chess.PositionBuilder
	}{
		{
			name: "king versus king",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
		},
		{
			name: "king bishop versus king",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "c1"), chess.White, chess.Bishop).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
		},
		{
			name: "king knight versus king",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "b1"), chess.White, chess.Knight).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
		},
		{
			name: "same color bishops",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "c1"), chess.White, chess.Bishop).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
				Place(mustParseSquare(t, "f8"), chess.Black, chess.Bishop),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			game := mustBuildGame(t, tc.builder)

			if !game.IsFinished() {
				t.Fatal("want draw")
			}

			if game.Outcome().Reason() != chess.OutcomeByInsufficientMaterial {
				t.Fatalf("expected insufficient-material outcome, got %s", game.Outcome().Reason())
			}

			if _, ok := game.Outcome().Winner(); ok {
				t.Fatal("want no winner")
			}
		})
	}
}

func TestLowMaterialNotDraw(t *testing.T) {
	testCases := []struct {
		name    string
		builder *chess.PositionBuilder
	}{
		{
			name: "opposite color bishops",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "c1"), chess.White, chess.Bishop).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
				Place(mustParseSquare(t, "c8"), chess.Black, chess.Bishop),
		},
		{
			name: "rook versus king",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
		},
		{
			name: "two knights versus king",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "b1"), chess.White, chess.Knight).
				Place(mustParseSquare(t, "g1"), chess.White, chess.Knight).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			game := mustBuildGame(t, tc.builder)

			if game.IsFinished() {
				t.Fatalf("expected game to continue, got %s", game.Outcome().Reason())
			}

			if game.Outcome().Reason() != chess.NoOutcomeReason {
				t.Fatalf("expected no outcome, got %s", game.Outcome().Reason())
			}
		})
	}
}
