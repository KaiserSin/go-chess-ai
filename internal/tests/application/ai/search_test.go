package ai_test

import (
	"testing"

	"github.com/KaiserSin/go-chess-ai/internal/application/ai"
	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestBuildTreeCreatesChildPerLegalMove(t *testing.T) {
	position := chess.NewInitialPosition()

	root := ai.BuildTree(position)

	if got, want := len(root.Children), len(position.LegalMoves()); got != want {
		t.Fatalf("want %d children, got %d", want, got)
	}
}

func TestBuildTreeUsesFixedDepthThree(t *testing.T) {
	position := chess.NewInitialPosition()

	root := ai.BuildTree(position)

	if len(root.Children) == 0 {
		t.Fatal("want root children")
	}

	firstChild := root.Children[0]
	if len(firstChild.Children) == 0 {
		t.Fatalf("want grandchildren for move %s", firstChild.Move)
	}

	if len(firstChild.Children[0].Children) == 0 {
		t.Fatalf("want great-grandchildren for move %s", firstChild.Move)
	}
}

func TestFixedSearchDepthIsThree(t *testing.T) {
	if ai.FixedSearchDepth != 3 {
		t.Fatalf("want fixed search depth 3, got %d", ai.FixedSearchDepth)
	}
}

func TestBestMoveChoosesWinningCapture(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Queen),
	)

	result := ai.BestMove(position)

	if !result.HasMove {
		t.Fatal("want best move")
	}

	if got := result.Move.String(); got != "a1a8" {
		t.Fatalf("want a1a8, got %s", got)
	}

	if result.Score <= 0 {
		t.Fatalf("want positive score, got %d", result.Score)
	}
}

func TestBestMoveTerminalPositionReturnsNoMove(t *testing.T) {
	game := chess.NewGame()
	applyMoves(t, game,
		mustMove(t, "f2", "f3"),
		mustMove(t, "e7", "e5"),
		mustMove(t, "g2", "g4"),
		mustMove(t, "d8", "h4"),
	)

	position := game.Position()
	result := ai.BestMove(position)

	if result.HasMove {
		t.Fatal("did not expect move")
	}

	if result.Score != -100000 {
		t.Fatalf("want -100000, got %d", result.Score)
	}
}

func TestBestMoveUsesFixedDepthChoiceFromInitialPosition(t *testing.T) {
	position := chess.NewInitialPosition()

	result := ai.BestMove(position)

	if !result.HasMove {
		t.Fatal("want best move")
	}

	if got := result.Move.String(); got != "b1c3" {
		t.Fatalf("want b1c3, got %s", got)
	}

	if result.Score != 20 {
		t.Fatalf("want 20, got %d", result.Score)
	}
}

func TestBestMoveAvoidsPoisonedCaptureAtFixedDepth(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "g1"), chess.White, chess.King).
			Place(mustParseSquare(t, "d1"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "g8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "d8"), chess.Black, chess.Queen).
			Place(mustParseSquare(t, "d7"), chess.Black, chess.Rook),
	)

	result := ai.BestMove(position)

	if !result.HasMove {
		t.Fatal("want best move")
	}

	if got := result.Move.String(); got == "d1d7" {
		t.Fatalf("did not expect poisoned capture d1d7, got %s", got)
	}
}
