package ai_test

import (
	"testing"

	"github.com/KaiserSin/go-chess-ai/internal/application/ai"
)

func TestBestMoveAvoidsPoisonedCaptureAtFixedDepth(t *testing.T) {
	position := poisonedCapturePosition(t)
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
