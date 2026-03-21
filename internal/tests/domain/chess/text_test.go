package chess_test

import (
	"testing"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestSideText(t *testing.T) {
	testCases := []struct {
		name string
		side chess.Side
		want string
	}{
		{name: "white", side: chess.White, want: "white"},
		{name: "black", side: chess.Black, want: "black"},
		{name: "bad", side: chess.Side(9), want: "side(9)"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.side.String(); got != tc.want {
				t.Fatalf("want %q, got %q", tc.want, got)
			}
		})
	}
}

func TestPieceTypeText(t *testing.T) {
	testCases := []struct {
		name string
		kind chess.PieceType
		want string
	}{
		{name: "none", kind: chess.NoPieceType, want: "none"},
		{name: "pawn", kind: chess.Pawn, want: "pawn"},
		{name: "knight", kind: chess.Knight, want: "knight"},
		{name: "bishop", kind: chess.Bishop, want: "bishop"},
		{name: "rook", kind: chess.Rook, want: "rook"},
		{name: "queen", kind: chess.Queen, want: "queen"},
		{name: "king", kind: chess.King, want: "king"},
		{name: "bad", kind: chess.PieceType(9), want: "piece_type(9)"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.kind.String(); got != tc.want {
				t.Fatalf("want %q, got %q", tc.want, got)
			}
		})
	}
}

func TestStatusText(t *testing.T) {
	testCases := []struct {
		name   string
		status chess.Status
		want   string
	}{
		{name: "ongoing", status: chess.Ongoing, want: "ongoing"},
		{name: "check", status: chess.Check, want: "check"},
		{name: "checkmate", status: chess.Checkmate, want: "checkmate"},
		{name: "stalemate", status: chess.Stalemate, want: "stalemate"},
		{name: "bad", status: chess.Status(9), want: "status(9)"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.status.String(); got != tc.want {
				t.Fatalf("want %q, got %q", tc.want, got)
			}
		})
	}
}

func TestOutcomeReasonText(t *testing.T) {
	testCases := []struct {
		name   string
		reason chess.OutcomeReason
		want   string
	}{
		{name: "none", reason: chess.NoOutcomeReason, want: "none"},
		{name: "checkmate", reason: chess.OutcomeByCheckmate, want: "checkmate"},
		{name: "stalemate", reason: chess.OutcomeByStalemate, want: "stalemate"},
		{name: "threefold", reason: chess.OutcomeByThreefoldRepetition, want: "same position 3 times"},
		{name: "fifty", reason: chess.OutcomeByFiftyMoveRule, want: "50-move rule"},
		{name: "material", reason: chess.OutcomeByInsufficientMaterial, want: "not enough material"},
		{name: "bad", reason: chess.OutcomeReason(9), want: "outcome_reason(9)"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.reason.String(); got != tc.want {
				t.Fatalf("want %q, got %q", tc.want, got)
			}
		})
	}
}

func TestPieceText(t *testing.T) {
	position := chess.NewInitialPosition()
	piece, ok := position.PieceAt(mustParseSquare(t, "e1"))
	if !ok {
		t.Fatal("want piece on e1")
	}

	var empty chess.Piece

	if got := empty.String(); got != "empty" {
		t.Fatalf("want %q, got %q", "empty", got)
	}

	if got := piece.String(); got != "white king" {
		t.Fatalf("want %q, got %q", "white king", got)
	}
}

func TestSquareText(t *testing.T) {
	if got := mustParseSquare(t, "a1").String(); got != "a1" {
		t.Fatalf("want %q, got %q", "a1", got)
	}

	if got := chess.Square(99).String(); got != "<invalid>" {
		t.Fatalf("want %q, got %q", "<invalid>", got)
	}
}

func TestMoveText(t *testing.T) {
	if got := mustMove(t, "e2", "e4").String(); got != "e2e4" {
		t.Fatalf("want %q, got %q", "e2e4", got)
	}

	if got := mustMove(t, "a7", "a8", chess.Queen).String(); got != "a7a8=q" {
		t.Fatalf("want %q, got %q", "a7a8=q", got)
	}
}
