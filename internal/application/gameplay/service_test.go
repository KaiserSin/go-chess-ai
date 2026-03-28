package gameplay

import (
	"testing"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"

	"github.com/KaiserSin/go-chess-ai/internal/application/dto"
)

func TestSnapshotInitialPosition(t *testing.T) {
	service := NewService()

	snapshot := service.Snapshot()

	if got := len(snapshot.Squares); got != 64 {
		t.Fatalf("want 64 squares, got %d", got)
	}

	if snapshot.SideToMove != "white" {
		t.Fatalf("want white to move, got %q", snapshot.SideToMove)
	}

	if square := squareByAlgebraic(t, snapshot, "e1"); square.PieceKey != "white-king" {
		t.Fatalf("want white-king on e1, got %q", square.PieceKey)
	}

	if square := squareByAlgebraic(t, snapshot, "d1"); square.PieceKey != "white-queen" {
		t.Fatalf("want white-queen on d1, got %q", square.PieceKey)
	}

	if square := squareByAlgebraic(t, snapshot, "e8"); square.PieceKey != "black-king" {
		t.Fatalf("want black-king on e8, got %q", square.PieceKey)
	}

	if square := squareByAlgebraic(t, snapshot, "a2"); square.PieceKey != "white-pawn" {
		t.Fatalf("want white-pawn on a2, got %q", square.PieceKey)
	}

	if square := squareByAlgebraic(t, snapshot, "h7"); square.PieceKey != "black-pawn" {
		t.Fatalf("want black-pawn on h7, got %q", square.PieceKey)
	}
}

func TestSelectSquareMarksSelectedPiece(t *testing.T) {
	service := NewService()

	service.SelectSquareAt(4, 1)

	snapshot := service.Snapshot()
	if square := squareByAlgebraic(t, snapshot, "e2"); !square.Selected {
		t.Fatal("want e2 selected")
	}
}

func TestSelectSquareMarksLegalTargets(t *testing.T) {
	service := NewService()

	service.SelectSquareAt(4, 1)

	snapshot := service.Snapshot()
	if square := squareByAlgebraic(t, snapshot, "e3"); !square.LegalTarget {
		t.Fatal("want e3 as legal target")
	}

	if square := squareByAlgebraic(t, snapshot, "e4"); !square.LegalTarget {
		t.Fatal("want e4 as legal target")
	}
}

func TestSelectingSameSquareClearsSelection(t *testing.T) {
	service := NewService()

	service.SelectSquareAt(4, 1)
	service.SelectSquareAt(4, 1)

	snapshot := service.Snapshot()
	if square := squareByAlgebraic(t, snapshot, "e2"); square.Selected {
		t.Fatal("did not expect e2 selected")
	}
}

func TestTryMoveUpdatesBoardAndSide(t *testing.T) {
	service := NewService()

	service.SelectSquareAt(4, 1)
	if err := service.TryMoveAt(4, 3); err != nil {
		t.Fatalf("want legal move, got %v", err)
	}

	snapshot := service.Snapshot()
	if snapshot.SideToMove != "black" {
		t.Fatalf("want black to move, got %q", snapshot.SideToMove)
	}

	if square := squareByAlgebraic(t, snapshot, "e4"); square.PieceKey != "white-pawn" {
		t.Fatalf("want white-pawn on e4, got %q", square.PieceKey)
	}

	if square := squareByAlgebraic(t, snapshot, "e2"); square.Occupied {
		t.Fatal("did not expect piece on e2")
	}
}

func TestMoveClearsSelection(t *testing.T) {
	service := NewService()

	service.SelectSquareAt(4, 1)
	if err := service.TryMoveAt(4, 3); err != nil {
		t.Fatalf("want legal move, got %v", err)
	}

	snapshot := service.Snapshot()
	if square := squareByAlgebraic(t, snapshot, "e2"); square.Selected {
		t.Fatal("did not expect e2 selected after move")
	}

	if square := squareByAlgebraic(t, snapshot, "e4"); square.Selected {
		t.Fatal("did not expect e4 selected after move")
	}
}

func TestPromotionShowsChoicesAndChoosePromotionWorks(t *testing.T) {
	service := newServiceWithGame(mustGameFromPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustSquareAt(t, 4, 0), chess.White, chess.King).
			Place(mustSquareAt(t, 4, 7), chess.Black, chess.King).
			Place(mustSquareAt(t, 6, 6), chess.White, chess.Pawn),
	))

	service.SelectSquareAt(6, 6)
	if err := service.TryMoveAt(6, 7); err != nil {
		t.Fatalf("want promotion selection, got %v", err)
	}

	snapshot := service.Snapshot()
	if snapshot.Promotion == nil || !snapshot.Promotion.Visible {
		t.Fatal("want visible promotion snapshot")
	}

	if got := len(snapshot.Promotion.Options); got != 4 {
		t.Fatalf("want 4 promotion options, got %d", got)
	}

	wantOrder := []string{"queen", "rook", "bishop", "knight"}
	for index, want := range wantOrder {
		if got := snapshot.Promotion.Options[index].PieceType; got != want {
			t.Fatalf("want promotion option %q at %d, got %q", want, index, got)
		}
	}

	if err := service.ChoosePromotion(chess.Queen); err != nil {
		t.Fatalf("want promotion success, got %v", err)
	}

	final := service.Snapshot()
	if final.Promotion != nil {
		t.Fatal("did not expect pending promotion after choosing piece")
	}

	if final.SideToMove != "black" {
		t.Fatalf("want black to move after promotion, got %q", final.SideToMove)
	}

	if square := squareByAlgebraic(t, final, "g8"); square.PieceKey != "white-queen" {
		t.Fatalf("want white-queen on g8, got %q", square.PieceKey)
	}
}

func squareByAlgebraic(t *testing.T, snapshot dto.GameSnapshot, algebraic string) dto.SquareSnapshot {
	t.Helper()

	for _, square := range snapshot.Squares {
		if square.Algebraic == algebraic {
			return square
		}
	}

	t.Fatalf("square %q not found", algebraic)
	return dto.SquareSnapshot{}
}

func mustGameFromPosition(t *testing.T, builder *chess.PositionBuilder) *chess.Game {
	t.Helper()

	position, err := builder.Build()
	if err != nil {
		t.Fatalf("build position: %v", err)
	}

	game, err := chess.NewGameFromPosition(position)
	if err != nil {
		t.Fatalf("new game from position: %v", err)
	}

	return game
}

func mustSquareAt(t *testing.T, file, rank int) chess.Square {
	t.Helper()

	square, err := chess.NewSquare(file, rank)
	if err != nil {
		t.Fatalf("new square: %v", err)
	}

	return square
}
