package viewmodel

import (
	"strconv"

	"github.com/KaiserSin/go-chess-ai/internal/application/dto"
	boardinput "github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/input"
	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/theme"
)

type Mapper struct {
	theme theme.Theme
}

type BoardViewModel struct {
	Title      string
	Status     string
	BoardX     int
	BoardY     int
	BoardSize  int
	Squares    []SquareViewModel
	FileLabels []AxisLabelViewModel
	RankLabels []AxisLabelViewModel
	Promotion  *PromotionViewModel
}

type SquareViewModel struct {
	Algebraic   string
	X           int
	Y           int
	Size        int
	IsLight     bool
	Selected    bool
	LegalTarget bool
	Piece       PieceViewModel
}

type PieceViewModel struct {
	Visible bool
	Visual  theme.PieceVisual
}

type AxisLabelViewModel struct {
	Text    string
	CenterX int
	CenterY int
}

type PromotionViewModel struct {
	Title   string
	Options []PromotionOptionViewModel
}

type PromotionOptionViewModel struct {
	PieceType string
	X         int
	Y         int
	Size      int
	Visual    theme.PieceVisual
}

func NewMapper(theme theme.Theme) *Mapper {
	return &Mapper{theme: theme}
}

func (m *Mapper) Map(snapshot dto.GameSnapshot) BoardViewModel {
	return BoardViewModel{
		Title:      "Go Chess AI",
		Status:     statusLine(snapshot),
		BoardX:     m.theme.BoardX,
		BoardY:     m.theme.BoardY,
		BoardSize:  m.theme.BoardSize,
		Squares:    m.mapSquares(snapshot.Squares),
		FileLabels: m.fileLabels(),
		RankLabels: m.rankLabels(),
		Promotion:  m.mapPromotion(snapshot.Promotion),
	}
}

func (m *Mapper) mapSquares(squares []dto.SquareSnapshot) []SquareViewModel {
	mapped := make([]SquareViewModel, 0, len(squares))
	for _, square := range squares {
		mapped = append(mapped, m.mapSquare(square))
	}

	return mapped
}

func (m *Mapper) mapSquare(square dto.SquareSnapshot) SquareViewModel {
	x := square.File * m.theme.SquareSize
	y := (7 - square.Rank) * m.theme.SquareSize

	mapped := SquareViewModel{
		Algebraic:   square.Algebraic,
		X:           x,
		Y:           y,
		Size:        m.theme.SquareSize,
		IsLight:     (square.File+square.Rank)%2 != 0,
		Selected:    square.Selected,
		LegalTarget: square.LegalTarget,
	}

	if square.Occupied {
		mapped.Piece = PieceViewModel{
			Visible: true,
			Visual:  m.theme.PieceCatalog.Lookup(square.PieceKey),
		}
	}

	return mapped
}

func (m *Mapper) fileLabels() []AxisLabelViewModel {
	labels := make([]AxisLabelViewModel, 0, 8)
	for file := 0; file < 8; file++ {
		labels = append(labels, AxisLabelViewModel{
			Text:    string(rune('a' + file)),
			CenterX: file*m.theme.SquareSize + m.theme.SquareSize/2,
			CenterY: m.theme.BoardSize + 22,
		})
	}

	return labels
}

func (m *Mapper) rankLabels() []AxisLabelViewModel {
	labels := make([]AxisLabelViewModel, 0, 8)
	for rank := 0; rank < 8; rank++ {
		labels = append(labels, AxisLabelViewModel{
			Text:    strconv.Itoa(rank + 1),
			CenterX: -18,
			CenterY: (7-rank)*m.theme.SquareSize + m.theme.SquareSize/2,
		})
	}

	return labels
}

func (m *Mapper) mapPromotion(promotion *dto.PromotionSnapshot) *PromotionViewModel {
	if promotion == nil || !promotion.Visible {
		return nil
	}

	rects := boardinput.PromotionOptionRects(m.theme, len(promotion.Options))
	mapped := &PromotionViewModel{
		Title:   "Choose promotion",
		Options: make([]PromotionOptionViewModel, 0, len(promotion.Options)),
	}

	for index, option := range promotion.Options {
		rect := rects[index]
		mapped.Options = append(mapped.Options, PromotionOptionViewModel{
			PieceType: option.PieceType,
			X:         rect.X,
			Y:         rect.Y,
			Size:      rect.Width,
			Visual:    m.theme.PieceCatalog.Lookup(option.PieceKey),
		})
	}

	return mapped
}

func statusLine(snapshot dto.GameSnapshot) string {
	if snapshot.OutcomeReason != "" && snapshot.OutcomeReason != "none" {
		if snapshot.HasWinner {
			return snapshot.Winner + " won by " + snapshot.OutcomeReason
		}

		return "draw by " + snapshot.OutcomeReason
	}

	if snapshot.Status == "check" {
		return snapshot.SideToMove + " to move · check"
	}

	return snapshot.SideToMove + " to move"
}
