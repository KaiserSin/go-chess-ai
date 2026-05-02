package ai

import (
	"testing"
	"time"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestExpiredDeadlineReturnsLegalFallbackMove(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "d1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook),
	)

	result := bestMoveWithOptions(position, searchOptions{
		deadline:      time.Now(),
		useAspiration: true,
	})
	if !result.HasMove {
		t.Fatal("want fallback move")
	}

	if !containsMove(position.LegalMoves(), result.Move) {
		t.Fatalf("want legal fallback move, got %s", result.Move)
	}
}

func TestAspirationSearchMatchesFullWindowSearch(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "c3"), chess.White, chess.King).
			Place(mustParseSquare(t, "d3"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a1"), chess.Black, chess.King),
	)

	previous, completed := searchAtDepth(position, 2, 0, searchOptions{
		useAspiration: false,
	})
	if !completed {
		t.Fatal("want completed previous depth")
	}

	withAspiration, completed := searchAtDepth(position, 3, previous.Score, searchOptions{
		useAspiration: true,
	})
	if !completed {
		t.Fatal("want completed aspiration search")
	}

	withoutAspiration, completed := searchAtDepth(position, 3, 0, searchOptions{
		useAspiration: false,
	})
	if !completed {
		t.Fatal("want completed full-window search")
	}

	if withAspiration != withoutAspiration {
		t.Fatalf("want matching search results, with aspiration=%+v without aspiration=%+v", withAspiration, withoutAspiration)
	}
}

func TestSearchAtDepthFindsForcedMateAtThreePlies(t *testing.T) {
	position := forcedMateDepthThreeSearchPosition(t)

	if !canForceMateForSearchTest(position, 3, chess.White) {
		t.Fatal("test fixture must contain forced mate within three plies")
	}

	result, completed := searchAtDepth(position, 3, 0, searchOptions{
		useAspiration: false,
	})
	if !completed {
		t.Fatal("want completed controlled-depth search")
	}

	next := mustApplyMoveForSearchTest(t, position, result.Move)
	if !canForceMateForSearchTest(next, 2, chess.White) {
		t.Fatalf("want move that keeps forced mate, got %s", result.Move)
	}
}

func TestQuiescenceSearchesEvasionsWhenSideToMoveIsInCheck(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "h7"), chess.Black, chess.King).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.Rook),
	)

	if position.Status() != chess.Check {
		t.Fatalf("want check, got %s", position.Status())
	}

	result := quiescence(position, -searchInfinity, searchInfinity, chess.White, searchOptions{})
	if !result.completed {
		t.Fatal("want completed quiescence search")
	}

	standPat := evaluateStatic(position, chess.White)
	if result.score == standPat {
		t.Fatalf("did not expect stand-pat score while in check, got %d", result.score)
	}

	expected := bestQuiescenceEvasionScore(t, position, chess.White)
	if result.score != expected {
		t.Fatalf("want best evasion score %d, got %d", expected, result.score)
	}
}

func mustBuildPosition(t *testing.T, builder *chess.PositionBuilder) chess.Position {
	t.Helper()

	position, err := builder.Build()
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	return position
}

func bestQuiescenceEvasionScore(t *testing.T, position chess.Position, rootPerspective chess.Side) int {
	t.Helper()

	moves := orderMoves(position, position.LegalMoves())
	if len(moves) == 0 {
		t.Fatal("test fixture must contain at least one legal evasion")
	}

	bestScore := 0
	for index, move := range moves {
		next := mustApplyMoveForSearchTest(t, position, move)
		result := quiescence(next, -searchInfinity, searchInfinity, rootPerspective, searchOptions{})
		if !result.completed {
			t.Fatal("want completed child quiescence search")
		}

		if index == 0 || result.score > bestScore {
			bestScore = result.score
		}
	}

	return bestScore
}

func mustParseSquare(t *testing.T, raw string) chess.Square {
	t.Helper()

	square, err := chess.ParseSquare(raw)
	if err != nil {
		t.Fatalf("ParseSquare(%q) error: %v", raw, err)
	}

	return square
}

func mustApplyMoveForSearchTest(t *testing.T, position chess.Position, move chess.Move) chess.Position {
	t.Helper()

	next, err := position.ApplyMove(move)
	if err != nil {
		t.Fatalf("ApplyMove(%s) error: %v", move, err)
	}

	return next
}

func forcedMateDepthThreeSearchPosition(t *testing.T) chess.Position {
	t.Helper()

	return mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "c3"), chess.White, chess.King).
			Place(mustParseSquare(t, "d3"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a1"), chess.Black, chess.King),
	)
}

func canForceMateForSearchTest(position chess.Position, plies int, attacker chess.Side) bool {
	if position.Status() == chess.Checkmate {
		return position.SideToMove().Opponent() == attacker
	}

	if position.Status() == chess.Stalemate || chess.HasInsufficientMaterial(position) || chess.IsFiftyMoveDraw(position) || plies == 0 {
		return false
	}

	moves := position.LegalMoves()
	if len(moves) == 0 {
		return false
	}

	if position.SideToMove() == attacker {
		for _, move := range moves {
			next, err := position.ApplyMove(move)
			if err != nil {
				panic(err)
			}

			if canForceMateForSearchTest(next, plies-1, attacker) {
				return true
			}
		}

		return false
	}

	for _, move := range moves {
		next, err := position.ApplyMove(move)
		if err != nil {
			panic(err)
		}

		if !canForceMateForSearchTest(next, plies-1, attacker) {
			return false
		}
	}

	return true
}

func containsMove(moves []chess.Move, target chess.Move) bool {
	for _, move := range moves {
		if move == target {
			return true
		}
	}

	return false
}
