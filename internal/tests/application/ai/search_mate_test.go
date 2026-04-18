package ai_test

import (
	"testing"

	"github.com/KaiserSin/go-chess-ai/internal/application/ai"
	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestBestMoveChoosesImmediateCheckmateWhenAvailable(t *testing.T) {
	position := immediateMatePosition(t)

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
	position := forcedMateDepthThreePosition(t)

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
