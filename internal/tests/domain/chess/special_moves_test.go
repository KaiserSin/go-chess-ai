package chess_test

import (
	"errors"
	"testing"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestKingsideCastlingMovesKingAndRook(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			WithCastlingRights(chess.WhiteKingSide|chess.WhiteQueenSide).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "h1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	move := mustMove(t, "e1", "g1")
	if !position.IsLegalMove(move) {
		t.Fatal("want legal kingside castling")
	}

	next, err := position.ApplyMove(move)
	if err != nil {
		t.Fatalf("ApplyMove error: %v", err)
	}

	assertPieceAt(t, next, "g1", chess.White, chess.King)
	assertPieceAt(t, next, "f1", chess.White, chess.Rook)

	if _, ok := next.PieceAt(mustParseSquare(t, "h1")); ok {
		t.Fatal("want empty h1 after castling")
	}

	if next.CastlingRights().CanCastleKingside(chess.White) || next.CastlingRights().CanCastleQueenside(chess.White) {
		t.Fatal("want no white castling rights")
	}
}

func TestQueensideCastlingBlockedByAttack(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			WithCastlingRights(chess.WhiteKingSide|chess.WhiteQueenSide).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "h8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "d8"), chess.Black, chess.Rook),
	)

	move := mustMove(t, "e1", "c1")
	if position.IsLegalMove(move) {
		t.Fatal("want illegal queenside castling")
	}
}

func TestEnPassantCaptureAndExpiration(t *testing.T) {
	game := chess.NewGame()
	applyMoves(t, game,
		mustMove(t, "e2", "e4"),
		mustMove(t, "a7", "a6"),
		mustMove(t, "e4", "e5"),
		mustMove(t, "d7", "d5"),
	)

	enPassant := mustMove(t, "e5", "d6")
	if !game.Position().IsLegalMove(enPassant) {
		t.Fatal("want legal en passant")
	}

	target, ok := game.Position().EnPassantSquare()
	if !ok || target != mustParseSquare(t, "d6") {
		t.Fatalf("want en passant d6, got %v, %t", target, ok)
	}

	if err := game.ApplyMove(enPassant); err != nil {
		t.Fatalf("ApplyMove error: %v", err)
	}

	position := game.Position()
	assertPieceAt(t, position, "d6", chess.White, chess.Pawn)
	if _, ok := position.PieceAt(mustParseSquare(t, "d5")); ok {
		t.Fatal("want no pawn on d5")
	}

	game = chess.NewGame()
	applyMoves(t, game,
		mustMove(t, "e2", "e4"),
		mustMove(t, "a7", "a6"),
		mustMove(t, "e4", "e5"),
		mustMove(t, "d7", "d5"),
		mustMove(t, "h2", "h3"),
		mustMove(t, "a6", "a5"),
	)

	if game.Position().IsLegalMove(enPassant) {
		t.Fatal("want illegal en passant")
	}
}

func TestPromotionRequiresExplicitChoice(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "g7"), chess.White, chess.Pawn),
	)

	_, err := position.ApplyMove(mustMove(t, "g7", "g8"))
	if !errors.Is(err, chess.ErrPromotionRequired) {
		t.Fatalf("expected ErrPromotionRequired, got %v", err)
	}

	_, err = position.ApplyMove(mustMove(t, "g7", "g8", chess.King))
	if !errors.Is(err, chess.ErrInvalidPromotion) {
		t.Fatalf("expected ErrInvalidPromotion, got %v", err)
	}

	next, err := position.ApplyMove(mustMove(t, "g7", "g8", chess.Queen))
	if err != nil {
		t.Fatalf("want good promotion, got %v", err)
	}

	assertPieceAt(t, next, "g8", chess.White, chess.Queen)
}

func TestCastlingRightsAndMoveCountersUpdate(t *testing.T) {
	game := chess.NewGame()
	applyMoves(t, game,
		mustMove(t, "h2", "h4"),
		mustMove(t, "a7", "a6"),
		mustMove(t, "h1", "h3"),
	)

	position := game.Position()
	if position.CastlingRights().CanCastleKingside(chess.White) {
		t.Fatal("want no white kingside castling")
	}

	if position.FullmoveNumber() != 2 {
		t.Fatalf("expected fullmove number 2, got %d", position.FullmoveNumber())
	}

	if position.HalfmoveClock() != 1 {
		t.Fatalf("expected halfmove clock 1, got %d", position.HalfmoveClock())
	}

	if position.SideToMove() != chess.Black {
		t.Fatalf("expected black to move, got %s", position.SideToMove())
	}
}

func assertPieceAt(t *testing.T, position chess.Position, rawSquare string, side chess.Side, pieceType chess.PieceType) {
	t.Helper()

	piece, ok := position.PieceAt(mustParseSquare(t, rawSquare))
	if !ok {
		t.Fatalf("want piece on %s", rawSquare)
	}

	if piece.Side() != side || piece.Type() != pieceType {
		t.Fatalf("expected %s %s on %s, got %s", side, pieceType, rawSquare, piece)
	}
}
