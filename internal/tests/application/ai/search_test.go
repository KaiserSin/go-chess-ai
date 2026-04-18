package ai_test

import (
	"testing"

	"github.com/KaiserSin/go-chess-ai/internal/application/ai"
	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestBestMoveReturnsLegalMoveFromInitialPosition(t *testing.T) {
	position := chess.NewInitialPosition()

	result := ai.BestMove(position)
	if !result.HasMove {
		t.Fatal("want best move")
	}

	if !containsMove(position.LegalMoves(), result.Move) {
		t.Fatalf("want legal move, got %s", result.Move)
	}
}

func TestBestMoveTerminalPositionReturnsNoMove(t *testing.T) {
	game := chess.NewGame()
	applyMoves(t, game,
		mustMove(t, "f2", "f3"),
		mustMove(t, "e7", "e5"),
		mustMove(t, "g2", "g4"),
		mustMove(t, "d8", "h4"),
	)

	position := game.Position()
	result := ai.BestMove(position)

	if result.HasMove {
		t.Fatal("did not expect move")
	}

	if result.Score >= 0 {
		t.Fatalf("want negative score for lost terminal position, got %d", result.Score)
	}
}

func TestBestMoveChoosesImmediateCheckmateWhenAvailable(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "d3"), chess.White, chess.King).
			Place(mustParseSquare(t, "f2"), chess.White, chess.Knight).
			Place(mustParseSquare(t, "g1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "c1"), chess.Black, chess.King),
	)

	if !hasImmediateMate(position) {
		t.Fatal("test fixture must contain immediate mate")
	}

	result := ai.BestMove(position)
	if !result.HasMove {
		t.Fatal("want best move")
	}

	next := mustApplyMove(t, position, result.Move)
	if next.Status() != chess.Checkmate {
		t.Fatalf("want checkmate after %s, got %s", result.Move, next.Status())
	}
}

func TestBestMoveFindsForcedMateWithinFixedDepth(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e2"), chess.White, chess.King).
			Place(mustParseSquare(t, "f2"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "b5"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "h4"), chess.Black, chess.King),
	)

	if hasImmediateMate(position) {
		t.Fatal("test fixture must require more than one move")
	}

	if !canForceMate(position, ai.FixedSearchDepth, chess.White) {
		t.Fatal("test fixture must contain forced mate within fixed depth")
	}

	result := ai.BestMove(position)
	if !result.HasMove {
		t.Fatal("want best move")
	}

	next := mustApplyMove(t, position, result.Move)
	if !canForceMate(next, ai.FixedSearchDepth-1, chess.White) {
		t.Fatalf("want move that keeps forced mate, got %s", result.Move)
	}
}

func TestBestMoveAvoidsPoisonedCaptureAtFixedDepth(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "g1"), chess.White, chess.King).
			Place(mustParseSquare(t, "d1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "g8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "d8"), chess.Black, chess.Queen).
			Place(mustParseSquare(t, "d7"), chess.Black, chess.Rook),
	)

	poisonedCapture := mustMove(t, "d1", "d7")
	if !containsMove(position.LegalMoves(), poisonedCapture) {
		t.Fatal("test fixture must offer poisoned capture")
	}

	result := ai.BestMove(position)

	if !result.HasMove {
		t.Fatal("want best move")
	}

	if result.Move == poisonedCapture {
		t.Fatalf("did not expect poisoned capture %s", poisonedCapture)
	}
}
