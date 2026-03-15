package chess_test

import (
	"errors"
	"testing"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestSquareParsingAndFormatting(t *testing.T) {
	square, err := chess.NewSquare(4, 3)
	if err != nil {
		t.Fatalf("NewSquare returned error: %v", err)
	}

	if square.String() != "e4" {
		t.Fatalf("expected e4, got %s", square)
	}

	parsed, err := chess.ParseSquare("e4")
	if err != nil {
		t.Fatalf("ParseSquare returned error: %v", err)
	}

	if parsed != square {
		t.Fatalf("expected parsed square %v to equal %v", parsed, square)
	}

	if _, err := chess.ParseSquare("i9"); !errors.Is(err, chess.ErrInvalidSquare) {
		t.Fatalf("expected ErrInvalidSquare, got %v", err)
	}
}

func TestInitialPositionHasTwentyLegalMoves(t *testing.T) {
	position := chess.NewInitialPosition()

	if position.SideToMove() != chess.White {
		t.Fatalf("expected white to move, got %s", position.SideToMove())
	}

	if got := len(position.LegalMoves()); got != 20 {
		t.Fatalf("expected 20 legal moves, got %d", got)
	}

	piece, ok := position.PieceAt(mustParseSquare(t, "e2"))
	if !ok {
		t.Fatal("expected piece on e2")
	}

	if piece.Side() != chess.White || piece.Type() != chess.Pawn {
		t.Fatalf("expected white pawn on e2, got %s", piece)
	}
}

func TestApplyMoveRejectsEmptySquareAndWrongSide(t *testing.T) {
	position := chess.NewInitialPosition()

	_, err := position.ApplyMove(chess.Move{From: mustParseSquare(t, "e3"), To: mustParseSquare(t, "e4")})
	if !errors.Is(err, chess.ErrNoPiece) {
		t.Fatalf("expected ErrNoPiece, got %v", err)
	}

	_, err = position.ApplyMove(chess.Move{From: mustParseSquare(t, "e7"), To: mustParseSquare(t, "e6")})
	if !errors.Is(err, chess.ErrWrongSide) {
		t.Fatalf("expected ErrWrongSide, got %v", err)
	}
}

func TestPinnedPieceCannotExposeKing(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "e2"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.Rook).
			Place(mustParseSquare(t, "h8"), chess.Black, chess.King),
	)

	illegal := chess.Move{From: mustParseSquare(t, "e2"), To: mustParseSquare(t, "d2")}
	if position.IsLegalMove(illegal) {
		t.Fatal("expected pinned rook move to be illegal")
	}

	_, err := position.ApplyMove(illegal)
	if !errors.Is(err, chess.ErrInvalidMove) {
		t.Fatalf("expected ErrInvalidMove, got %v", err)
	}
}

func TestCheckStatusWithoutMate(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.Black).
			Place(mustParseSquare(t, "h1"), chess.White, chess.King).
			Place(mustParseSquare(t, "e1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	if position.Status() != chess.Check {
		t.Fatalf("expected check status, got %s", position.Status())
	}

	if !position.IsInCheck(chess.Black) {
		t.Fatal("expected black king to be in check")
	}

	if len(position.LegalMoves()) == 0 {
		t.Fatal("expected at least one legal move")
	}
}

func TestFoolsMateProducesCheckmate(t *testing.T) {
	game := chess.NewGame()

	applyMoves(t, game,
		mustMove(t, "f2", "f3"),
		mustMove(t, "e7", "e5"),
		mustMove(t, "g2", "g4"),
		mustMove(t, "d8", "h4"),
	)

	if game.Status() != chess.Checkmate {
		t.Fatalf("expected checkmate, got %s", game.Status())
	}

	if err := game.ApplyMove(mustMove(t, "g2", "g3")); !errors.Is(err, chess.ErrGameFinished) {
		t.Fatalf("expected ErrGameFinished after checkmate, got %v", err)
	}
}

func TestStalemateStatus(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.Black).
			Place(mustParseSquare(t, "c6"), chess.White, chess.King).
			Place(mustParseSquare(t, "b6"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.King),
	)

	if position.Status() != chess.Stalemate {
		t.Fatalf("expected stalemate, got %s", position.Status())
	}

	if len(position.LegalMoves()) != 0 {
		t.Fatalf("expected no legal moves, got %d", len(position.LegalMoves()))
	}
}

func mustBuildPosition(t *testing.T, builder *chess.PositionBuilder) chess.Position {
	t.Helper()

	position, err := builder.Build()
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}

	return position
}

func mustParseSquare(t *testing.T, raw string) chess.Square {
	t.Helper()

	square, err := chess.ParseSquare(raw)
	if err != nil {
		t.Fatalf("ParseSquare(%q) returned error: %v", raw, err)
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
			t.Fatalf("ApplyMove(%s) returned error: %v", move, err)
		}
	}
}
