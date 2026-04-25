package ai

import (
	"testing"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestEvaluateRewardsMaterialAdvantage(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "d1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	score := Evaluate(position, chess.White)
	if score < 800 {
		t.Fatalf("want clear material advantage for White, got %d", score)
	}
}

func TestEvaluateUsesBlackPerspectiveSign(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "d1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	whiteScore := Evaluate(position, chess.White)
	blackScore := Evaluate(position, chess.Black)
	if whiteScore != -blackScore {
		t.Fatalf("want opposite perspective scores, white=%d black=%d", whiteScore, blackScore)
	}
}

func TestEvaluateRewardsCentralPiecePlacement(t *testing.T) {
	centerKnight := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "a2"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "d4"), chess.White, chess.Knight).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "h7"), chess.Black, chess.Queen),
	)
	edgeKnight := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "a2"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Knight).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "h7"), chess.Black, chess.Queen),
	)

	centerScore := Evaluate(centerKnight, chess.White)
	edgeScore := Evaluate(edgeKnight, chess.White)
	if centerScore <= edgeScore {
		t.Fatalf("want central knight score above edge knight score, center=%d edge=%d", centerScore, edgeScore)
	}
}

func TestEvaluatePenalizesWeakPawnStructure(t *testing.T) {
	connectedPawns := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "c4"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "d5"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)
	doubledIsolatedPawns := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "c4"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "c5"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	connectedScore := Evaluate(connectedPawns, chess.White)
	doubledScore := Evaluate(doubledIsolatedPawns, chess.White)
	if connectedScore <= doubledScore {
		t.Fatalf("want connected pawns above doubled isolated pawns, connected=%d doubled=%d", connectedScore, doubledScore)
	}
}

func TestEvaluateRewardsKingShieldInMiddlegame(t *testing.T) {
	shieldedKing := middlegameKingShieldPosition(t, "f2", "g2", "h2")
	exposedKing := middlegameKingShieldPosition(t, "f3", "g3", "h3")

	shieldedScore := Evaluate(shieldedKing, chess.White)
	exposedScore := Evaluate(exposedKing, chess.White)
	if shieldedScore <= exposedScore {
		t.Fatalf("want shielded king above exposed king, shielded=%d exposed=%d", shieldedScore, exposedScore)
	}
}

func TestEvaluateRewardsActiveKingInEndgame(t *testing.T) {
	activeKing := endgameKingPosition(t, "d4")
	cornerKing := endgameKingPosition(t, "a1")

	activeScore := Evaluate(activeKing, chess.White)
	cornerScore := Evaluate(cornerKing, chess.White)
	if activeScore <= cornerScore {
		t.Fatalf("want active king above corner king, active=%d corner=%d", activeScore, cornerScore)
	}
}

func middlegameKingShieldPosition(t *testing.T, firstPawn, secondPawn, thirdPawn string) chess.Position {
	t.Helper()

	return mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "g1"), chess.White, chess.King).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "b1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, firstPawn), chess.White, chess.Pawn).
			Place(mustParseSquare(t, secondPawn), chess.White, chess.Pawn).
			Place(mustParseSquare(t, thirdPawn), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "g8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Queen).
			Place(mustParseSquare(t, "b8"), chess.Black, chess.Rook).
			Place(mustParseSquare(t, "f7"), chess.Black, chess.Pawn).
			Place(mustParseSquare(t, "g7"), chess.Black, chess.Pawn).
			Place(mustParseSquare(t, "h7"), chess.Black, chess.Pawn),
	)
}

func endgameKingPosition(t *testing.T, whiteKing string) chess.Position {
	t.Helper()

	return mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, whiteKing), chess.White, chess.King).
			Place(mustParseSquare(t, "a2"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "h8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "h7"), chess.Black, chess.Pawn),
	)
}
