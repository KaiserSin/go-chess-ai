package chess_test

import "testing"

import chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"

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
		t.Fatal("want draw")
	}

	if game.Position().HalfmoveClock() != 100 {
		t.Fatalf("want 100, got %d", game.Position().HalfmoveClock())
	}

	if game.Outcome().Reason() != chess.OutcomeByFiftyMoveRule {
		t.Fatalf("want fifty rule, got %s", game.Outcome().Reason())
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
		t.Fatalf("want 0, got %d", game.Position().HalfmoveClock())
	}

	if game.IsFinished() {
		t.Fatal("want game on")
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
		t.Fatalf("want 0, got %d", game.Position().HalfmoveClock())
	}

	if game.IsFinished() {
		t.Fatal("want game on")
	}
}

func TestMoveCounters(t *testing.T) {
	game := chess.NewGame()
	applyMoves(t, game,
		mustMove(t, "h2", "h4"),
		mustMove(t, "a7", "a6"),
		mustMove(t, "h1", "h3"),
	)

	position := game.Position()
	if position.FullmoveNumber() != 2 {
		t.Fatalf("want 2, got %d", position.FullmoveNumber())
	}

	if position.HalfmoveClock() != 1 {
		t.Fatalf("want 1, got %d", position.HalfmoveClock())
	}

	if position.SideToMove() != chess.Black {
		t.Fatalf("want black, got %s", position.SideToMove())
	}
}

func TestLowMaterialDraw(t *testing.T) {
	testCases := []struct {
		name    string
		builder *chess.PositionBuilder
	}{
		{
			name: "king vs king",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
		},
		{
			name: "bishop vs king",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "c1"), chess.White, chess.Bishop).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
		},
		{
			name: "knight vs king",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "b1"), chess.White, chess.Knight).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
		},
		{
			name: "same bishops",
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
				t.Fatalf("want low material, got %s", game.Outcome().Reason())
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
			name: "other bishops",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "c1"), chess.White, chess.Bishop).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
				Place(mustParseSquare(t, "c8"), chess.Black, chess.Bishop),
		},
		{
			name: "rook vs king",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
		},
		{
			name: "two knights",
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
				t.Fatalf("want game on, got %s", game.Outcome().Reason())
			}

			if game.Outcome().Reason() != chess.NoOutcomeReason {
				t.Fatalf("want no outcome, got %s", game.Outcome().Reason())
			}
		})
	}
}

func TestLowMaterialWithPawn(t *testing.T) {
	game := mustBuildGame(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "a2"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	if game.IsFinished() {
		t.Fatal("want game on")
	}
}
