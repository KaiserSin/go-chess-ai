package viewmodel_test

import (
	"testing"

	"github.com/KaiserSin/go-chess-ai/internal/application/dto"
	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/theme"
	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/viewmodel"
)

func TestMapperPromotionOverlayOrderAndPosition(t *testing.T) {
	mapper := viewmodel.NewMapper(theme.NewTheme())

	board := mapper.Map(dto.GameSnapshot{
		SideToMove: "white",
		Promotion: &dto.PromotionSnapshot{
			Visible: true,
			Options: []dto.PromotionOptionSnapshot{
				{PieceType: "queen", PieceKey: "white-queen"},
				{PieceType: "rook", PieceKey: "white-rook"},
				{PieceType: "bishop", PieceKey: "white-bishop"},
				{PieceType: "knight", PieceKey: "white-knight"},
			},
		},
	})

	if board.Promotion == nil {
		t.Fatal("want promotion overlay")
	}

	if got := board.Promotion.Options[0].PieceType; got != "queen" {
		t.Fatalf("want queen first, got %s", got)
	}

	if got := board.Promotion.Options[3].PieceType; got != "knight" {
		t.Fatalf("want knight last, got %s", got)
	}

	if option := board.Promotion.Options[0]; option.X != 200 || option.Y != 368 {
		t.Fatalf("want first promotion option at (200, 368), got (%d, %d)", option.X, option.Y)
	}
}
