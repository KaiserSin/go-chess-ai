package gameplay_test

import (
	"errors"
	"testing"

	"github.com/KaiserSin/go-chess-ai/internal/application/gameplay"
	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestApplyAIMoveMakesMoveAndChangesTurn(t *testing.T) {
	service := gameplay.NewService()

	if err := service.ApplyAIMove(); err != nil {
		t.Fatalf("want ai move, got %v", err)
	}

	snapshot := service.Snapshot()
	if snapshot.SideToMove != "black" {
		t.Fatalf("want black to move, got %q", snapshot.SideToMove)
	}

	if square := squareByAlgebraic(t, snapshot, "a3"); square.PieceKey != "white-pawn" {
		t.Fatalf("want white-pawn on a3, got %q", square.PieceKey)
	}

	if square := squareByAlgebraic(t, snapshot, "a2"); square.Occupied {
		t.Fatal("did not expect piece on a2")
	}
}

func TestApplyAIMoveReturnsGameFinishedWhenGameIsOver(t *testing.T) {
	service := gameplay.NewService()
	mustPlayMove(t, service, 5, 1, 5, 2)
	mustPlayMove(t, service, 4, 6, 4, 4)
	mustPlayMove(t, service, 6, 1, 6, 3)
	mustPlayMove(t, service, 3, 7, 7, 3)

	if err := service.ApplyAIMove(); !errors.Is(err, chess.ErrGameFinished) {
		t.Fatalf("want game finished error, got %v", err)
	}
}

func TestApplyAIMoveReturnsInvalidMoveWhenPromotionIsPending(t *testing.T) {
	service := gameplay.NewService()
	mustPlayMove(t, service, 7, 1, 7, 3)
	mustPlayMove(t, service, 6, 6, 6, 4)
	mustPlayMove(t, service, 7, 3, 6, 4)
	mustPlayMove(t, service, 6, 7, 5, 5)
	mustPlayMove(t, service, 6, 4, 6, 5)
	mustPlayMove(t, service, 0, 6, 0, 5)
	mustPlayMove(t, service, 6, 5, 6, 6)
	mustPlayMove(t, service, 0, 5, 0, 4)
	mustPlayMove(t, service, 6, 6, 6, 7)

	snapshot := service.Snapshot()
	if snapshot.Promotion == nil || !snapshot.Promotion.Visible {
		t.Fatal("want pending promotion")
	}

	if err := service.ApplyAIMove(); !errors.Is(err, chess.ErrInvalidMove) {
		t.Fatalf("want invalid move error, got %v", err)
	}
}

func mustPlayMove(t *testing.T, service *gameplay.Service, fromFile, fromRank, toFile, toRank int) {
	t.Helper()

	service.SelectSquareAt(fromFile, fromRank)
	if err := service.TryMoveAt(toFile, toRank); err != nil {
		t.Fatalf("want move %d,%d -> %d,%d, got %v", fromFile, fromRank, toFile, toRank, err)
	}
}
