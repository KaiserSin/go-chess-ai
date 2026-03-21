package chess_test

import "testing"

import chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"

func TestWhiteEnPassant(t *testing.T) {
	game := chess.NewGame()
	applyMoves(t, game,
		mustMove(t, "e2", "e4"),
		mustMove(t, "a7", "a6"),
		mustMove(t, "e4", "e5"),
		mustMove(t, "d7", "d5"),
	)

	move := mustMove(t, "e5", "d6")
	if !game.Position().IsLegalMove(move) {
		t.Fatal("want good en passant")
	}

	target, ok := game.Position().EnPassantSquare()
	if !ok || target != mustParseSquare(t, "d6") {
		t.Fatalf("want d6, got %v, %t", target, ok)
	}

	if err := game.ApplyMove(move); err != nil {
		t.Fatalf("want good en passant, got %v", err)
	}

	position := game.Position()
	assertPieceAt(t, position, "d6", chess.White, chess.Pawn)
	if _, ok := position.PieceAt(mustParseSquare(t, "d5")); ok {
		t.Fatal("want no pawn on d5")
	}
}

func TestWhiteEnPassantExpires(t *testing.T) {
	game := chess.NewGame()
	applyMoves(t, game,
		mustMove(t, "e2", "e4"),
		mustMove(t, "a7", "a6"),
		mustMove(t, "e4", "e5"),
		mustMove(t, "d7", "d5"),
		mustMove(t, "h2", "h3"),
		mustMove(t, "a6", "a5"),
	)

	if game.Position().IsLegalMove(mustMove(t, "e5", "d6")) {
		t.Fatal("want bad en passant")
	}
}

func TestBlackEnPassant(t *testing.T) {
	game := chess.NewGame()
	applyMoves(t, game,
		mustMove(t, "a2", "a3"),
		mustMove(t, "h7", "h5"),
		mustMove(t, "a3", "a4"),
		mustMove(t, "h5", "h4"),
		mustMove(t, "g2", "g4"),
	)

	move := mustMove(t, "h4", "g3")
	if !game.Position().IsLegalMove(move) {
		t.Fatal("want good en passant")
	}

	if err := game.ApplyMove(move); err != nil {
		t.Fatalf("want good en passant, got %v", err)
	}

	position := game.Position()
	assertPieceAt(t, position, "g3", chess.Black, chess.Pawn)
	if _, ok := position.PieceAt(mustParseSquare(t, "g4")); ok {
		t.Fatal("want no pawn on g4")
	}
}
