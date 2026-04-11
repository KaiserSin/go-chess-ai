package ai_test

import (
	"testing"

	"github.com/KaiserSin/go-chess-ai/internal/application/ai"
	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestEvaluateInitialPosition(t *testing.T) {
	position := chess.NewInitialPosition()

	if got := ai.Evaluate(position, chess.White); got != 0 {
		t.Fatalf("want 0 for white, got %d", got)
	}

	if got := ai.Evaluate(position, chess.Black); got != 0 {
		t.Fatalf("want 0 for black, got %d", got)
	}
}

func TestEvaluateMaterialAdvantageForWhite(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "d1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	if got := ai.Evaluate(position, chess.White); got != 900 {
		t.Fatalf("want 900 for white, got %d", got)
	}

	if got := ai.Evaluate(position, chess.Black); got != -900 {
		t.Fatalf("want -900 for black, got %d", got)
	}
}

func TestEvaluateMaterialAdvantageForBlack(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook),
	)

	if got := ai.Evaluate(position, chess.White); got != -500 {
		t.Fatalf("want -500 for white, got %d", got)
	}

	if got := ai.Evaluate(position, chess.Black); got != 500 {
		t.Fatalf("want 500 for black, got %d", got)
	}
}

func TestEvaluateCheckmate(t *testing.T) {
	game := chess.NewGame()
	applyMoves(t, game,
		mustMove(t, "f2", "f3"),
		mustMove(t, "e7", "e5"),
		mustMove(t, "g2", "g4"),
		mustMove(t, "d8", "h4"),
	)

	position := game.Position()
	if got := ai.Evaluate(position, chess.Black); got != 100000 {
		t.Fatalf("want 100000 for black, got %d", got)
	}

	if got := ai.Evaluate(position, chess.White); got != -100000 {
		t.Fatalf("want -100000 for white, got %d", got)
	}
}

func TestEvaluateStalemateReturnsZero(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.Black).
			Place(mustParseSquare(t, "c6"), chess.White, chess.King).
			Place(mustParseSquare(t, "b6"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.King),
	)

	if got := ai.Evaluate(position, chess.White); got != 0 {
		t.Fatalf("want 0 for white, got %d", got)
	}

	if got := ai.Evaluate(position, chess.Black); got != 0 {
		t.Fatalf("want 0 for black, got %d", got)
	}
}

func TestEvaluateInsufficientMaterialReturnsZero(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	if got := ai.Evaluate(position, chess.White); got != 0 {
		t.Fatalf("want 0 for white, got %d", got)
	}

	if got := ai.Evaluate(position, chess.Black); got != 0 {
		t.Fatalf("want 0 for black, got %d", got)
	}
}

func TestEvaluateFiftyMoveDrawReturnsZero(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			WithHalfmoveClock(100).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	if got := ai.Evaluate(position, chess.White); got != 0 {
		t.Fatalf("want 0 for white, got %d", got)
	}

	if got := ai.Evaluate(position, chess.Black); got != 0 {
		t.Fatalf("want 0 for black, got %d", got)
	}
}

func TestEvaluateInvalidPerspectiveReturnsZero(t *testing.T) {
	position := chess.NewInitialPosition()

	if got := ai.Evaluate(position, chess.Side(9)); got != 0 {
		t.Fatalf("want 0 for invalid side, got %d", got)
	}
}

func TestEvaluateKnightInCenterBetterThanKnightOnRim(t *testing.T) {
	centerKnight := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "d4"), chess.White, chess.Knight).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)
	rimKnight := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Knight).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	if ai.Evaluate(centerKnight, chess.White) <= ai.Evaluate(rimKnight, chess.White) {
		t.Fatal("want center knight to evaluate better than rim knight")
	}
}

func TestEvaluateCentralPawnBetterThanEdgePawn(t *testing.T) {
	centerPawn := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "d1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "d4"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "d8"), chess.Black, chess.Queen).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook),
	)
	edgePawn := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "d1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "a4"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "d8"), chess.Black, chess.Queen).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook),
	)

	if ai.Evaluate(centerPawn, chess.White) <= ai.Evaluate(edgePawn, chess.White) {
		t.Fatal("want central pawn to evaluate better than edge pawn")
	}
}

func TestEvaluateActiveBishopBetterThanPassiveBishop(t *testing.T) {
	activeBishop := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "d4"), chess.White, chess.Bishop).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)
	passiveBishop := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Bishop).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	if ai.Evaluate(activeBishop, chess.White) <= ai.Evaluate(passiveBishop, chess.White) {
		t.Fatal("want active bishop to evaluate better than passive bishop")
	}
}

func TestEvaluateRookAndQueenPieceSquareTablesAreModerate(t *testing.T) {
	centerQueen := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "d4"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)
	cornerQueen := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)
	centerRook := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "d4"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)
	cornerRook := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	queenDiff := ai.Evaluate(centerQueen, chess.White) - ai.Evaluate(cornerQueen, chess.White)
	if queenDiff <= 0 || queenDiff >= 50 {
		t.Fatalf("want moderate positive queen PST diff, got %d", queenDiff)
	}

	rookDiff := ai.Evaluate(centerRook, chess.White) - ai.Evaluate(cornerRook, chess.White)
	if rookDiff <= 0 || rookDiff >= 50 {
		t.Fatalf("want moderate positive rook PST diff, got %d", rookDiff)
	}
}

func TestEvaluateKingUsesMiddlegameTableWhenHeavyPiecesRemain(t *testing.T) {
	safeKing := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "g1"), chess.White, chess.King).
			Place(mustParseSquare(t, "d1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "d8"), chess.Black, chess.Queen).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook).
			Place(mustParseSquare(t, "g8"), chess.Black, chess.King),
	)
	exposedKing := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e4"), chess.White, chess.King).
			Place(mustParseSquare(t, "d1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "d8"), chess.Black, chess.Queen).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook).
			Place(mustParseSquare(t, "g8"), chess.Black, chess.King),
	)

	if ai.Evaluate(safeKing, chess.White) <= ai.Evaluate(exposedKing, chess.White) {
		t.Fatal("want castled middlegame king to evaluate better than exposed king")
	}
}

func TestEvaluateKingUsesEndgameTableWhenMaterialIsLow(t *testing.T) {
	activeKing := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "d4"), chess.White, chess.King).
			Place(mustParseSquare(t, "a2"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)
	passiveKing := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "g1"), chess.White, chess.King).
			Place(mustParseSquare(t, "a2"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	if ai.Evaluate(activeKing, chess.White) <= ai.Evaluate(passiveKing, chess.White) {
		t.Fatal("want active endgame king to evaluate better than passive king")
	}
}

func TestEvaluateKingPhaseDependsOnMaterialNotSideToMove(t *testing.T) {
	whiteToMove := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "d4"), chess.White, chess.King).
			Place(mustParseSquare(t, "a2"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)
	blackToMove := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.Black).
			Place(mustParseSquare(t, "d4"), chess.White, chess.King).
			Place(mustParseSquare(t, "a2"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	if ai.Evaluate(whiteToMove, chess.White) != ai.Evaluate(blackToMove, chess.White) {
		t.Fatal("want same evaluation when only side to move changes in equal endgame material")
	}
}

func TestEvaluateDoubledPawnsAreWorseThanHealthyPawns(t *testing.T) {
	healthy := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "g1"), chess.White, chess.King).
			Place(mustParseSquare(t, "c2"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "d2"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "g8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "c7"), chess.Black, chess.Pawn).
			Place(mustParseSquare(t, "d7"), chess.Black, chess.Pawn),
	)
	doubled := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "g1"), chess.White, chess.King).
			Place(mustParseSquare(t, "c2"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "c3"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "g8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "c7"), chess.Black, chess.Pawn).
			Place(mustParseSquare(t, "d7"), chess.Black, chess.Pawn),
	)

	if ai.Evaluate(doubled, chess.White) >= ai.Evaluate(healthy, chess.White) {
		t.Fatal("want doubled pawns to evaluate worse than healthy pawns")
	}
}

func TestEvaluateIsolatedPawnIsWorseThanSupportedPawn(t *testing.T) {
	supported := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "d1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "c4"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "d4"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "d8"), chess.Black, chess.Queen).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook).
			Place(mustParseSquare(t, "c6"), chess.Black, chess.Pawn).
			Place(mustParseSquare(t, "e6"), chess.Black, chess.Pawn),
	)
	isolated := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "d1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "a4"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "d4"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "d8"), chess.Black, chess.Queen).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook).
			Place(mustParseSquare(t, "a6"), chess.Black, chess.Pawn).
			Place(mustParseSquare(t, "e6"), chess.Black, chess.Pawn),
	)

	if ai.Evaluate(isolated, chess.White) >= ai.Evaluate(supported, chess.White) {
		t.Fatal("want isolated pawn structure to evaluate worse than supported pawns")
	}
}

func TestEvaluatePassedPawnIsBetterThanContestedPawn(t *testing.T) {
	passed := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "d1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "d6"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "d8"), chess.Black, chess.Queen).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook).
			Place(mustParseSquare(t, "a6"), chess.Black, chess.Pawn),
	)
	contested := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "d1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "d6"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "d8"), chess.Black, chess.Queen).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook).
			Place(mustParseSquare(t, "e7"), chess.Black, chess.Pawn),
	)

	if ai.Evaluate(passed, chess.White) <= ai.Evaluate(contested, chess.White) {
		t.Fatal("want passed pawn to evaluate better than contested pawn")
	}
}

func TestEvaluateKingWithPawnShieldIsSaferInMiddlegame(t *testing.T) {
	shielded := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "d1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "d2"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "e2"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "f2"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "d8"), chess.Black, chess.Queen).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook).
			Place(mustParseSquare(t, "d7"), chess.Black, chess.Pawn).
			Place(mustParseSquare(t, "e7"), chess.Black, chess.Pawn).
			Place(mustParseSquare(t, "f7"), chess.Black, chess.Pawn),
	)
	exposed := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "d1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "a2"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "b2"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "c2"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "d8"), chess.Black, chess.Queen).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook).
			Place(mustParseSquare(t, "d7"), chess.Black, chess.Pawn).
			Place(mustParseSquare(t, "e7"), chess.Black, chess.Pawn).
			Place(mustParseSquare(t, "f7"), chess.Black, chess.Pawn),
	)

	if ai.Evaluate(shielded, chess.White) <= ai.Evaluate(exposed, chess.White) {
		t.Fatal("want king with pawn shield to evaluate better than exposed king")
	}
}

func TestEvaluateKingInCenterIsWorseInMiddlegame(t *testing.T) {
	safe := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "g1"), chess.White, chess.King).
			Place(mustParseSquare(t, "d1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "f2"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "g2"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "h2"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "g8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "d8"), chess.Black, chess.Queen).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook).
			Place(mustParseSquare(t, "f7"), chess.Black, chess.Pawn).
			Place(mustParseSquare(t, "g7"), chess.Black, chess.Pawn).
			Place(mustParseSquare(t, "h7"), chess.Black, chess.Pawn),
	)
	center := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e4"), chess.White, chess.King).
			Place(mustParseSquare(t, "d1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "f2"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "g2"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "h2"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "g8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "d8"), chess.Black, chess.Queen).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook).
			Place(mustParseSquare(t, "f7"), chess.Black, chess.Pawn).
			Place(mustParseSquare(t, "g7"), chess.Black, chess.Pawn).
			Place(mustParseSquare(t, "h7"), chess.Black, chess.Pawn),
	)

	if ai.Evaluate(center, chess.White) >= ai.Evaluate(safe, chess.White) {
		t.Fatal("want king in center to evaluate worse than safer king in middlegame")
	}
}

func TestEvaluateAdvancedPassedPawnIsBetterInEndgame(t *testing.T) {
	advanced := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e4"), chess.White, chess.King).
			Place(mustParseSquare(t, "d6"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "a7"), chess.Black, chess.Pawn),
	)
	lessAdvanced := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e4"), chess.White, chess.King).
			Place(mustParseSquare(t, "d4"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "a7"), chess.Black, chess.Pawn),
	)

	if ai.Evaluate(advanced, chess.White) <= ai.Evaluate(lessAdvanced, chess.White) {
		t.Fatal("want advanced passed pawn to evaluate better in endgame")
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

func mustParseSquare(t *testing.T, raw string) chess.Square {
	t.Helper()

	square, err := chess.ParseSquare(raw)
	if err != nil {
		t.Fatalf("ParseSquare(%q) error: %v", raw, err)
	}

	return square
}

func mustMove(t *testing.T, from, to string, promotion ...chess.PieceType) chess.Move {
	t.Helper()

	move := chess.Move{
		From: mustParseSquare(t, from),
		To:   mustParseSquare(t, to),
	}

	if len(promotion) > 0 {
		move.Promotion = promotion[0]
	}

	return move
}

func applyMoves(t *testing.T, game *chess.Game, moves ...chess.Move) {
	t.Helper()

	for _, move := range moves {
		if err := game.ApplyMove(move); err != nil {
			t.Fatalf("ApplyMove(%s) error: %v", move, err)
		}
	}
}
