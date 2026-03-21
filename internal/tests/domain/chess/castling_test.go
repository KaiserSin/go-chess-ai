package chess_test

import "testing"

import chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"

func TestCastlingWorks(t *testing.T) {
	testCases := []struct {
		name     string
		builder  *chess.PositionBuilder
		move     chess.Move
		kingTo   string
		rookTo   string
		rookFrom string
		side     chess.Side
	}{
		{
			name: "white king side",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				WithCastlingRights(chess.WhiteKingSide).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "h1"), chess.White, chess.Rook).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
			move:     mustMove(t, "e1", "g1"),
			kingTo:   "g1",
			rookTo:   "f1",
			rookFrom: "h1",
			side:     chess.White,
		},
		{
			name: "white queen side",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				WithCastlingRights(chess.WhiteQueenSide).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
			move:     mustMove(t, "e1", "c1"),
			kingTo:   "c1",
			rookTo:   "d1",
			rookFrom: "a1",
			side:     chess.White,
		},
		{
			name: "black king side",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.Black).
				WithCastlingRights(chess.BlackKingSide).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
				Place(mustParseSquare(t, "h8"), chess.Black, chess.Rook),
			move:     mustMove(t, "e8", "g8"),
			kingTo:   "g8",
			rookTo:   "f8",
			rookFrom: "h8",
			side:     chess.Black,
		},
		{
			name: "black queen side",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.Black).
				WithCastlingRights(chess.BlackQueenSide).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
				Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook),
			move:     mustMove(t, "e8", "c8"),
			kingTo:   "c8",
			rookTo:   "d8",
			rookFrom: "a8",
			side:     chess.Black,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			position := mustBuildPosition(t, tc.builder)
			if !position.IsLegalMove(tc.move) {
				t.Fatal("want good castle")
			}

			next, err := position.ApplyMove(tc.move)
			if err != nil {
				t.Fatalf("want good castle, got %v", err)
			}

			assertPieceAt(t, next, tc.kingTo, tc.side, chess.King)
			assertPieceAt(t, next, tc.rookTo, tc.side, chess.Rook)

			if _, ok := next.PieceAt(mustParseSquare(t, tc.rookFrom)); ok {
				t.Fatalf("want empty %s", tc.rookFrom)
			}
		})
	}
}

func TestCastlingBlocked(t *testing.T) {
	testCases := []struct {
		name    string
		builder *chess.PositionBuilder
		move    chess.Move
	}{
		{
			name: "king side",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				WithCastlingRights(chess.WhiteKingSide).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "f1"), chess.White, chess.Bishop).
				Place(mustParseSquare(t, "h1"), chess.White, chess.Rook).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
			move: mustMove(t, "e1", "g1"),
		},
		{
			name: "queen side",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				WithCastlingRights(chess.WhiteQueenSide).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
				Place(mustParseSquare(t, "h8"), chess.Black, chess.King).
				Place(mustParseSquare(t, "d8"), chess.Black, chess.Rook),
			move: mustMove(t, "e1", "c1"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			position := mustBuildPosition(t, tc.builder)
			if position.IsLegalMove(tc.move) {
				t.Fatal("want bad castle")
			}
		})
	}
}

func TestCastlingNeedsRook(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			WithCastlingRights(chess.WhiteKingSide|chess.WhiteQueenSide).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	if position.IsLegalMove(mustMove(t, "e1", "g1")) {
		t.Fatal("want bad castle")
	}

	if position.IsLegalMove(mustMove(t, "e1", "c1")) {
		t.Fatal("want bad castle")
	}
}

func TestRookMoveRemovesRight(t *testing.T) {
	testCases := []struct {
		name      string
		builder   *chess.PositionBuilder
		move      chess.Move
		side      chess.Side
		wantKing  bool
		wantQueen bool
	}{
		{
			name: "white h1",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				WithCastlingRights(chess.WhiteKingSide|chess.WhiteQueenSide).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "h1"), chess.White, chess.Rook).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
			move:      mustMove(t, "h1", "h3"),
			side:      chess.White,
			wantKing:  false,
			wantQueen: true,
		},
		{
			name: "white a1",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				WithCastlingRights(chess.WhiteKingSide|chess.WhiteQueenSide).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
			move:      mustMove(t, "a1", "a3"),
			side:      chess.White,
			wantKing:  true,
			wantQueen: false,
		},
		{
			name: "black h8",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.Black).
				WithCastlingRights(chess.BlackKingSide|chess.BlackQueenSide).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
				Place(mustParseSquare(t, "h8"), chess.Black, chess.Rook),
			move:      mustMove(t, "h8", "h6"),
			side:      chess.Black,
			wantKing:  false,
			wantQueen: true,
		},
		{
			name: "black a8",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.Black).
				WithCastlingRights(chess.BlackKingSide|chess.BlackQueenSide).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
				Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook),
			move:      mustMove(t, "a8", "a6"),
			side:      chess.Black,
			wantKing:  true,
			wantQueen: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			position := mustBuildPosition(t, tc.builder)
			next, err := position.ApplyMove(tc.move)
			if err != nil {
				t.Fatalf("want good rook move, got %v", err)
			}

			rights := next.CastlingRights()
			if rights.CanCastleKingside(tc.side) != tc.wantKing {
				t.Fatalf("bad king side right for %s", tc.side)
			}

			if rights.CanCastleQueenside(tc.side) != tc.wantQueen {
				t.Fatalf("bad queen side right for %s", tc.side)
			}
		})
	}
}

func TestRookCaptureRemovesRights(t *testing.T) {
	testCases := []struct {
		name      string
		builder   *chess.PositionBuilder
		move      chess.Move
		side      chess.Side
		wantKing  bool
		wantQueen bool
	}{
		{
			name: "a1",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.Black).
				WithCastlingRights(chess.WhiteKingSide|chess.WhiteQueenSide|chess.BlackKingSide|chess.BlackQueenSide).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
				Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook),
			move:      mustMove(t, "a8", "a1"),
			side:      chess.White,
			wantKing:  true,
			wantQueen: false,
		},
		{
			name: "h1",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.Black).
				WithCastlingRights(chess.WhiteKingSide|chess.WhiteQueenSide|chess.BlackKingSide|chess.BlackQueenSide).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "h1"), chess.White, chess.Rook).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
				Place(mustParseSquare(t, "h8"), chess.Black, chess.Rook),
			move:      mustMove(t, "h8", "h1"),
			side:      chess.White,
			wantKing:  false,
			wantQueen: true,
		},
		{
			name: "a8",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				WithCastlingRights(chess.WhiteKingSide|chess.WhiteQueenSide|chess.BlackKingSide|chess.BlackQueenSide).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
				Place(mustParseSquare(t, "a8"), chess.Black, chess.Rook),
			move:      mustMove(t, "a1", "a8"),
			side:      chess.Black,
			wantKing:  true,
			wantQueen: false,
		},
		{
			name: "h8",
			builder: chess.NewPositionBuilder().
				WithSideToMove(chess.White).
				WithCastlingRights(chess.WhiteKingSide|chess.WhiteQueenSide|chess.BlackKingSide|chess.BlackQueenSide).
				Place(mustParseSquare(t, "e1"), chess.White, chess.King).
				Place(mustParseSquare(t, "h1"), chess.White, chess.Rook).
				Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
				Place(mustParseSquare(t, "h8"), chess.Black, chess.Rook),
			move:      mustMove(t, "h1", "h8"),
			side:      chess.Black,
			wantKing:  false,
			wantQueen: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			position := mustBuildPosition(t, tc.builder)
			next, err := position.ApplyMove(tc.move)
			if err != nil {
				t.Fatalf("want good capture, got %v", err)
			}

			rights := next.CastlingRights()
			if rights.CanCastleKingside(tc.side) != tc.wantKing {
				t.Fatalf("bad king side right for %s", tc.side)
			}

			if rights.CanCastleQueenside(tc.side) != tc.wantQueen {
				t.Fatalf("bad queen side right for %s", tc.side)
			}
		})
	}
}
