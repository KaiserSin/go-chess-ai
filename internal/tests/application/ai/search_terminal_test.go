package ai_test

import (
	"testing"

	"github.com/KaiserSin/go-chess-ai/internal/application/ai"
	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestBestMoveTerminalCheckmateReturnsNoMove(t *testing.T) {
	position := checkmateLossPosition(t)

	if position.Status() != chess.Checkmate {
		t.Fatalf("want checkmate, got %s", position.Status())
	}

	result := ai.BestMove(position)
	if result.HasMove {
		t.Fatal("did not expect move")
	}

	if result.Score >= 0 {
		t.Fatalf("want negative score for lost terminal position, got %d", result.Score)
	}
}

func TestBestMoveTerminalStalemateReturnsNoMove(t *testing.T) {
	position := stalematePosition(t)

	if position.Status() != chess.Stalemate {
		t.Fatalf("want stalemate, got %s", position.Status())
	}

	result := ai.BestMove(position)
	if result.HasMove {
		t.Fatal("did not expect move")
	}
}
