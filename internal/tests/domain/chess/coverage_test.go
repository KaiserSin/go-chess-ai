package chess_test

import (
	"errors"
	"testing"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestText(t *testing.T) {
	position := chess.NewInitialPosition()
	piece, ok := position.PieceAt(mustParseSquare(t, "e1"))
	if !ok {
		t.Fatal("want piece on e1")
	}

	var empty chess.Piece

	testCases := []struct {
		name string
		got  string
		want string
	}{
		{name: "side", got: chess.White.String(), want: "white"},
		{name: "black side", got: chess.Black.String(), want: "black"},
		{name: "bad side", got: chess.Side(9).String(), want: "side(9)"},
		{name: "none type", got: chess.NoPieceType.String(), want: "none"},
		{name: "pawn type", got: chess.Pawn.String(), want: "pawn"},
		{name: "piece type", got: chess.Knight.String(), want: "knight"},
		{name: "bishop type", got: chess.Bishop.String(), want: "bishop"},
		{name: "rook type", got: chess.Rook.String(), want: "rook"},
		{name: "queen type", got: chess.Queen.String(), want: "queen"},
		{name: "king type", got: chess.King.String(), want: "king"},
		{name: "bad piece type", got: chess.PieceType(9).String(), want: "piece_type(9)"},
		{name: "ongoing status", got: chess.Ongoing.String(), want: "ongoing"},
		{name: "check status", got: chess.Check.String(), want: "check"},
		{name: "status", got: chess.Checkmate.String(), want: "checkmate"},
		{name: "stalemate status", got: chess.Stalemate.String(), want: "stalemate"},
		{name: "bad status", got: chess.Status(9).String(), want: "status(9)"},
		{name: "no reason", got: chess.NoOutcomeReason.String(), want: "none"},
		{name: "mate reason", got: chess.OutcomeByCheckmate.String(), want: "checkmate"},
		{name: "stalemate reason", got: chess.OutcomeByStalemate.String(), want: "stalemate"},
		{name: "three reason", got: chess.OutcomeByThreefoldRepetition.String(), want: "same position 3 times"},
		{name: "fifty reason", got: chess.OutcomeByFiftyMoveRule.String(), want: "50-move rule"},
		{name: "reason", got: chess.OutcomeByInsufficientMaterial.String(), want: "not enough material"},
		{name: "bad reason", got: chess.OutcomeReason(9).String(), want: "outcome_reason(9)"},
		{name: "empty piece", got: empty.String(), want: "empty"},
		{name: "piece", got: piece.String(), want: "white king"},
		{name: "square", got: mustParseSquare(t, "a1").String(), want: "a1"},
		{name: "bad square", got: chess.Square(99).String(), want: "<invalid>"},
		{name: "move", got: mustMove(t, "e2", "e4").String(), want: "e2e4"},
		{name: "promotion move", got: mustMove(t, "a7", "a8", chess.Queen).String(), want: "a7a8=q"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got != tc.want {
				t.Fatalf("want %q, got %q", tc.want, tc.got)
			}
		})
	}

	if _, err := chess.ParseSquare("a"); !errors.Is(err, chess.ErrInvalidSquare) {
		t.Fatalf("want bad square error, got %v", err)
	}
}

func TestOutcomeFlags(t *testing.T) {
	winGame := chess.NewGame()
	applyMoves(t, winGame,
		mustMove(t, "f2", "f3"),
		mustMove(t, "e7", "e5"),
		mustMove(t, "g2", "g4"),
		mustMove(t, "d8", "h4"),
	)

	if !winGame.Outcome().IsDecisive() {
		t.Fatal("want win")
	}

	if got := len(chess.NewGame().LegalMoves()); got != 20 {
		t.Fatalf("want 20 moves, got %d", got)
	}

	if moves := winGame.LegalMoves(); moves != nil {
		t.Fatal("want no moves")
	}

	drawGame := mustBuildGame(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.Black).
			Place(mustParseSquare(t, "c6"), chess.White, chess.King).
			Place(mustParseSquare(t, "b6"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.King),
	)

	if drawGame.Outcome().IsDecisive() {
		t.Fatal("want no win")
	}

	if moves := drawGame.LegalMoves(); moves != nil {
		t.Fatal("want no moves")
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

func TestBadPromotionUse(t *testing.T) {
	rookPosition := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook),
	)

	_, err := rookPosition.ApplyMove(mustMove(t, "a1", "a2", chess.Queen))
	if !errors.Is(err, chess.ErrInvalidPromotion) {
		t.Fatalf("want bad promotion error, got %v", err)
	}

	pawnPosition := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "g6"), chess.White, chess.Pawn),
	)

	_, err = pawnPosition.ApplyMove(mustMove(t, "g6", "g7", chess.Queen))
	if !errors.Is(err, chess.ErrInvalidPromotion) {
		t.Fatalf("want bad promotion error, got %v", err)
	}

	_, err = pawnPosition.ApplyMove(mustMove(t, "g6", "g7", chess.PieceType(9)))
	if !errors.Is(err, chess.ErrInvalidPromotion) {
		t.Fatalf("want bad promotion error, got %v", err)
	}
}

func TestBadBuilder(t *testing.T) {
	testCases := []struct {
		name    string
		builder *chess.PositionBuilder
		want    error
	}{
		{
			name:    "bad side",
			builder: chess.NewPositionBuilder().WithSideToMove(chess.Side(9)),
			want:    chess.ErrInvalidPosition,
		},
		{
			name:    "bad en passant",
			builder: chess.NewPositionBuilder().WithEnPassantSquare(chess.Square(99)),
			want:    chess.ErrInvalidSquare,
		},
		{
			name:    "bad halfmove",
			builder: chess.NewPositionBuilder().WithHalfmoveClock(-1),
			want:    chess.ErrInvalidPosition,
		},
		{
			name:    "bad fullmove",
			builder: chess.NewPositionBuilder().WithFullmoveNumber(0),
			want:    chess.ErrInvalidPosition,
		},
		{
			name: "bad piece",
			builder: chess.NewPositionBuilder().
				Place(mustParseSquare(t, "e4"), chess.White, chess.PieceType(9)),
			want: chess.ErrInvalidPosition,
		},
		{
			name:    "bad place square",
			builder: chess.NewPositionBuilder().Place(chess.Square(99), chess.White, chess.Pawn),
			want:    chess.ErrInvalidSquare,
		},
		{
			name: "missing king",
			builder: chess.NewPositionBuilder().
				Place(mustParseSquare(t, "e1"), chess.White, chess.King),
			want: chess.ErrInvalidPosition,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.builder.Build()
			if !errors.Is(err, tc.want) {
				t.Fatalf("want %v, got %v", tc.want, err)
			}
		})
	}
}

func TestBuilderStops(t *testing.T) {
	builder := chess.NewPositionBuilder().
		WithSideToMove(chess.Side(9)).
		WithSideToMove(chess.Black).
		WithCastlingRights(chess.WhiteKingSide).
		WithEnPassantSquare(mustParseSquare(t, "e3")).
		WithHalfmoveClock(4).
		WithFullmoveNumber(2).
		Place(mustParseSquare(t, "e1"), chess.White, chess.King)

	_, err := builder.Build()
	if !errors.Is(err, chess.ErrInvalidPosition) {
		t.Fatalf("want bad position error, got %v", err)
	}
}

func TestBadGame(t *testing.T) {
	if _, err := chess.NewGameFromPosition(chess.Position{}); !errors.Is(err, chess.ErrInvalidPosition) {
		t.Fatalf("want bad position error, got %v", err)
	}

	if err := chess.NewGame().ApplyMove(mustMove(t, "e3", "e4")); !errors.Is(err, chess.ErrNoPiece) {
		t.Fatalf("want no piece error, got %v", err)
	}

	if (chess.Position{}).IsInCheck(chess.White) {
		t.Fatal("want no check")
	}
}

func TestPositionOver(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.Black).
			Place(mustParseSquare(t, "c6"), chess.White, chess.King).
			Place(mustParseSquare(t, "b6"), chess.White, chess.Queen).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.King),
	)

	_, err := position.ApplyMove(mustMove(t, "a8", "a7"))
	if !errors.Is(err, chess.ErrGameFinished) {
		t.Fatalf("want game over error, got %v", err)
	}
}

func TestMoreCastles(t *testing.T) {
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
				t.Fatal("want legal castle")
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

func TestCastleNeedsRook(t *testing.T) {
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

func TestCastleKingSideBlocked(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			WithCastlingRights(chess.WhiteKingSide).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "f1"), chess.White, chess.Bishop).
			Place(mustParseSquare(t, "h1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	if position.IsLegalMove(mustMove(t, "e1", "g1")) {
		t.Fatal("want bad castle")
	}
}

func TestRookCaptureRights(t *testing.T) {
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

func TestLowMaterialWithPawn(t *testing.T) {
	game := mustBuildGame(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "a2"), chess.White, chess.Pawn).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King),
	)

	if game.IsFinished() {
		t.Fatal("want game on")
	}
}
