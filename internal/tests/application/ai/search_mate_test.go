package ai_test

import (
	"testing"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestBestMoveChoosesImmediateCheckmateWhenAvailable(t *testing.T) {
	position := immediateMatePosition(t)

	if !hasImmediateMate(position) {
		t.Fatal("test fixture must contain immediate mate")
	}

	result := bestMoveForTest(t, position)
	if !result.HasMove {
		t.Fatal("want best move")
	}

	next := mustApplyMove(t, position, result.Move)
	if next.Status() != chess.Checkmate {
		t.Fatalf("want checkmate after %s, got %s", result.Move, next.Status())
	}
}

func TestBestMoveFindsForcedMateWithinTimeBudget(t *testing.T) {
	position := forcedMateDepthThreePosition(t)
	const forcedMateDepth = 3

	if hasImmediateMate(position) {
		t.Fatal("test fixture must require more than one move")
	}

	if !canForceMate(position, forcedMateDepth, chess.White) {
		t.Fatal("test fixture must contain forced mate within three plies")
	}

	result := bestMoveForTest(t, position)
	if !result.HasMove {
		t.Fatal("want best move")
	}

	next := mustApplyMove(t, position, result.Move)
	if !canForceMate(next, forcedMateDepth-1, chess.White) {
		t.Fatalf("want move that keeps forced mate, got %s", result.Move)
	}
}
