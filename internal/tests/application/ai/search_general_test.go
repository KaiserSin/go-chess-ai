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

func TestBestMoveReturnsOnlyLegalMoveWhileInCheck(t *testing.T) {
	position := onlyLegalMoveUnderCheckPosition(t)
	onlyMove := mustMove(t, "a8", "b7")

	if position.Status() != chess.Check {
		t.Fatalf("want check, got %s", position.Status())
	}

	legalMoves := position.LegalMoves()
	if len(legalMoves) != 1 {
		t.Fatalf("want exactly one legal move, got %d", len(legalMoves))
	}

	if legalMoves[0] != onlyMove {
		t.Fatalf("want only move %s, got %s", onlyMove, legalMoves[0])
	}

	result := ai.BestMove(position)
	if !result.HasMove {
		t.Fatal("want best move")
	}

	if result.Move != onlyMove {
		t.Fatalf("want only legal move %s, got %s", onlyMove, result.Move)
	}
}

func TestBestMoveChoosesLegalPromotionMove(t *testing.T) {
	position := promotionReadyPosition(t)

	promotionMoves := 0
	for _, move := range position.LegalMoves() {
		if move.Promotion != chess.NoPieceType {
			promotionMoves++
		}
	}

	if promotionMoves == 0 {
		t.Fatal("test fixture must offer promotion moves")
	}

	result := ai.BestMove(position)
	if !result.HasMove {
		t.Fatal("want best move")
	}

	if !containsMove(position.LegalMoves(), result.Move) {
		t.Fatalf("want legal move, got %s", result.Move)
	}

	if result.Move.Promotion == chess.NoPieceType {
		t.Fatalf("want promotion move, got %s", result.Move)
	}

	next := mustApplyMove(t, position, result.Move)
	piece, ok := next.PieceAt(mustParseSquare(t, "g8"))
	if !ok {
		t.Fatal("want promoted piece on g8")
	}

	if piece.Side() != chess.White || piece.Type() != result.Move.Promotion {
		t.Fatalf("want promoted %s on g8, got %s", result.Move.Promotion, piece)
	}
}
