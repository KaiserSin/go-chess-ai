package viewmodel

import (
	"fmt"

	"github.com/KaiserSin/go-chess-ai/internal/application/dto"
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
}

type SquareViewModel struct {
	Algebraic string
	X         int
	Y         int
	Size      int
	IsLight   bool
	Piece     PieceViewModel
}

type PieceViewModel struct {
	Visible bool
	CenterX int
	CenterY int
	Visual  theme.PieceVisual
}

type AxisLabelViewModel struct {
	Text    string
	CenterX int
	CenterY int
}

func NewMapper(theme theme.Theme) *Mapper {
	return &Mapper{theme: theme}
}

func (m *Mapper) Map(snapshot dto.GameSnapshot) BoardViewModel {
	board := BoardViewModel{
		Title:      "Go Chess AI",
		Status:     statusLine(snapshot),
		BoardX:     m.theme.BoardX,
		BoardY:     m.theme.BoardY,
		BoardSize:  m.theme.BoardSize,
		Squares:    make([]SquareViewModel, 0, len(snapshot.Squares)),
		FileLabels: make([]AxisLabelViewModel, 0, 8),
		RankLabels: make([]AxisLabelViewModel, 0, 8),
	}

	for _, square := range snapshot.Squares {
		x := square.File * m.theme.SquareSize
		y := (7 - square.Rank) * m.theme.SquareSize

		squareView := SquareViewModel{
			Algebraic: square.Algebraic,
			X:         x,
			Y:         y,
			Size:      m.theme.SquareSize,
			IsLight:   (square.File+square.Rank)%2 != 0,
		}

		if square.Occupied {
			squareView.Piece = PieceViewModel{
				Visible: true,
				CenterX: x + m.theme.SquareSize/2,
				CenterY: y + m.theme.SquareSize/2,
				Visual:  m.theme.PieceCatalog.Lookup(square.PieceKey),
			}
		}

		board.Squares = append(board.Squares, squareView)
	}

	for file := 0; file < 8; file++ {
		board.FileLabels = append(board.FileLabels, AxisLabelViewModel{
			Text:    string(rune('a' + file)),
			CenterX: file*m.theme.SquareSize + m.theme.SquareSize/2,
			CenterY: m.theme.BoardSize + 22,
		})
	}

	for rank := 0; rank < 8; rank++ {
		board.RankLabels = append(board.RankLabels, AxisLabelViewModel{
			Text:    fmt.Sprintf("%d", rank+1),
			CenterX: -18,
			CenterY: (7-rank)*m.theme.SquareSize + m.theme.SquareSize/2,
		})
	}

	return board
}

func statusLine(snapshot dto.GameSnapshot) string {
	if snapshot.OutcomeReason != "" && snapshot.OutcomeReason != "none" {
		if snapshot.HasWinner {
			return fmt.Sprintf("%s won by %s", snapshot.Winner, snapshot.OutcomeReason)
		}

		return fmt.Sprintf("draw by %s", snapshot.OutcomeReason)
	}

	return fmt.Sprintf("%s to move", snapshot.SideToMove)
}
