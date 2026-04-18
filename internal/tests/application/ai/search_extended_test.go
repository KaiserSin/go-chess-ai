//go:build extended
// +build extended

package ai_test

import (
	"testing"

	"github.com/KaiserSin/go-chess-ai/internal/application/ai"
	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

type extendedPositionCase struct {
	name     string
	position chess.Position
}

func TestBestMoveExtendedCorpusNonTerminalPositions(t *testing.T) {
	for _, tc := range extendedNonTerminalCases(t) {
		t.Run(tc.name, func(t *testing.T) {
			result := deterministicBestMove(t, tc.position)
			if !result.HasMove {
				t.Fatal("want best move")
			}

			if !containsMove(tc.position.LegalMoves(), result.Move) {
				t.Fatalf("want legal move, got %s", result.Move)
			}

			mover := tc.position.SideToMove()
			next := mustApplyMove(t, tc.position, result.Move)
			if tc.position.Status() == chess.Check && next.IsInCheck(mover) {
				t.Fatalf("did not expect %s to remain in check after %s", mover, result.Move)
			}
		})
	}
}

func TestBestMoveExtendedCorpusTerminalPositions(t *testing.T) {
	for _, tc := range extendedTerminalCases(t) {
		t.Run(tc.name, func(t *testing.T) {
			result := ai.BestMove(tc.position)
			if result.HasMove {
				t.Fatalf("did not expect move, got %s", result.Move)
			}
		})
	}
}

func deterministicBestMove(t *testing.T, position chess.Position) ai.SearchResult {
	t.Helper()

	first := ai.BestMove(position)
	for attempt := 0; attempt < 4; attempt++ {
		next := ai.BestMove(position)
		if next != first {
			t.Fatalf("want deterministic result, first=%+v next=%+v", first, next)
		}
	}

	return first
}

func extendedNonTerminalCases(t *testing.T) []extendedPositionCase {
	t.Helper()

	return []extendedPositionCase{
		{
			name: "e2e4_e7e5_g1f3_b8c6",
			position: mustPositionFromMoves(t,
				mustMove(t, "e2", "e4"),
				mustMove(t, "e7", "e5"),
				mustMove(t, "g1", "f3"),
				mustMove(t, "b8", "c6"),
			),
		},
		{
			name: "d2d4_d7d5_c2c4",
			position: mustPositionFromMoves(t,
				mustMove(t, "d2", "d4"),
				mustMove(t, "d7", "d5"),
				mustMove(t, "c2", "c4"),
			),
		},
		{
			name: "sicilian_open",
			position: mustPositionFromMoves(t,
				mustMove(t, "e2", "e4"),
				mustMove(t, "c7", "c5"),
				mustMove(t, "g1", "f3"),
				mustMove(t, "d7", "d6"),
				mustMove(t, "d2", "d4"),
				mustMove(t, "c5", "d4"),
				mustMove(t, "f3", "d4"),
			),
		},
		{
			name: "english_four_moves",
			position: mustPositionFromMoves(t,
				mustMove(t, "c2", "c4"),
				mustMove(t, "e7", "e5"),
				mustMove(t, "b1", "c3"),
				mustMove(t, "g8", "f6"),
			),
		},
		{name: "poisoned_capture", position: poisonedCapturePosition(t)},
		{name: "immediate_mate", position: immediateMatePosition(t)},
		{name: "forced_mate_depth_three", position: forcedMateDepthThreePosition(t)},
		{name: "only_legal_move_under_check", position: onlyLegalMoveUnderCheckPosition(t)},
		{name: "promotion_ready", position: promotionReadyPosition(t)},
	}
}

func extendedTerminalCases(t *testing.T) []extendedPositionCase {
	t.Helper()

	return []extendedPositionCase{
		{name: "checkmate", position: checkmateLossPosition(t)},
		{name: "stalemate", position: stalematePosition(t)},
	}
}
