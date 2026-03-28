package gameplay

import (
	"testing"

	"github.com/KaiserSin/go-chess-ai/internal/application/dto"
	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

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
