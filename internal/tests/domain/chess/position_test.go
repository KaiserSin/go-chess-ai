package chess_test

import (
	"errors"
	"testing"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestSquare(t *testing.T) {
	square, err := chess.NewSquare(4, 3)
	if err != nil {
		t.Fatalf("NewSquare error: %v", err)
	}

	parsed, err := chess.ParseSquare("e4")
	if err != nil {
		t.Fatalf("ParseSquare error: %v", err)
	}

	if parsed != square {
		t.Fatalf("want %v, got %v", square, parsed)
	}

	if _, err := chess.ParseSquare("i9"); !errors.Is(err, chess.ErrInvalidSquare) {
		t.Fatalf("want bad square error, got %v", err)
	}

	if _, err := chess.ParseSquare("a"); !errors.Is(err, chess.ErrInvalidSquare) {
		t.Fatalf("want bad square error, got %v", err)
	}

	if _, err := chess.NewSquare(8, 0); !errors.Is(err, chess.ErrInvalidSquare) {
		t.Fatalf("want bad square error, got %v", err)
	}
}

func TestStartPosition(t *testing.T) {
	position := chess.NewInitialPosition()

	if position.SideToMove() != chess.White {
		t.Fatalf("want white, got %s", position.SideToMove())
	}

	if position.HalfmoveClock() != 0 {
		t.Fatalf("want 0, got %d", position.HalfmoveClock())
	}

	if position.FullmoveNumber() != 1 {
		t.Fatalf("want 1, got %d", position.FullmoveNumber())
	}

	if got := len(position.LegalMoves()); got != 20 {
		t.Fatalf("want 20 moves, got %d", got)
	}

	assertPieceAt(t, position, "e2", chess.White, chess.Pawn)

	game := chess.NewGame()
	if got := len(game.LegalMoves()); got != 20 {
		t.Fatalf("want 20 moves, got %d", got)
	}
}

func TestBadMoveInput(t *testing.T) {
	position := chess.NewInitialPosition()

	_, err := position.ApplyMove(chess.Move{From: mustParseSquare(t, "e3"), To: mustParseSquare(t, "e4")})
	if !errors.Is(err, chess.ErrNoPiece) {
		t.Fatalf("want no piece error, got %v", err)
	}

	_, err = position.ApplyMove(chess.Move{From: mustParseSquare(t, "e7"), To: mustParseSquare(t, "e6")})
	if !errors.Is(err, chess.ErrWrongSide) {
		t.Fatalf("want wrong side error, got %v", err)
	}

	game := chess.NewGame()
	err = game.ApplyMove(mustMove(t, "e3", "e4"))
	if !errors.Is(err, chess.ErrNoPiece) {
		t.Fatalf("want no piece error, got %v", err)
	}
}

func TestBadMoveSquares(t *testing.T) {
	position := chess.NewInitialPosition()
	move := chess.Move{
		From: chess.Square(99),
		To:   mustParseSquare(t, "e4"),
	}

	if position.IsLegalMove(move) {
		t.Fatal("want bad move")
	}

	_, err := position.ApplyMove(move)
	if !errors.Is(err, chess.ErrInvalidSquare) {
		t.Fatalf("want bad square error, got %v", err)
	}
}

func TestPinnedMove(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "e2"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.Rook).
			Place(mustParseSquare(t, "h8"), chess.Black, chess.King),
	)

	move := chess.Move{From: mustParseSquare(t, "e2"), To: mustParseSquare(t, "d2")}
	if position.IsLegalMove(move) {
		t.Fatal("want bad move")
	}

	_, err := position.ApplyMove(move)
	if !errors.Is(err, chess.ErrInvalidMove) {
		t.Fatalf("want bad move error, got %v", err)
	}
}
