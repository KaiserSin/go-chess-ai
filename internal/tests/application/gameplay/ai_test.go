package gameplay_test

import (
	"errors"
	"testing"

	"github.com/KaiserSin/go-chess-ai/internal/application/ai"
	"github.com/KaiserSin/go-chess-ai/internal/application/gameplay"
	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestApplyAIMoveMakesMoveAndChangesTurn(t *testing.T) {
	service := gameplay.NewService()
	before := service.Snapshot()
	expected := ai.BestMove(chess.NewInitialPosition())
	if !expected.HasMove {
		t.Fatal("want best move for initial position")
	}

	if err := service.ApplyAIMove(); err != nil {
		t.Fatalf("want ai move, got %v", err)
	}

	snapshot := service.Snapshot()
	if snapshot.SideToMove != "black" {
		t.Fatalf("want black to move, got %q", snapshot.SideToMove)
	}

	movedPiece := squareByAlgebraic(t, before, expected.Move.From.String()).PieceKey
	if square := squareByAlgebraic(t, snapshot, expected.Move.From.String()); square.Occupied {
		t.Fatalf("did not expect piece on %s after ai move", expected.Move.From)
	}

	if square := squareByAlgebraic(t, snapshot, expected.Move.To.String()); square.PieceKey != movedPiece {
		t.Fatalf("want %s on %s, got %q", movedPiece, expected.Move.To, square.PieceKey)
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
