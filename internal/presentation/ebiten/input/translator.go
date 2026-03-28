package input

import (
	"fmt"

	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/theme"
)

type Translator struct {
	theme theme.Theme
}

type SquareTarget struct {
	File int
	Rank int
}

type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

func NewTranslator(theme theme.Theme) *Translator {
	return &Translator{theme: theme}
}

func (t *Translator) SquareAt(screenX, screenY int) (SquareTarget, bool) {
	if screenX < t.theme.BoardX || screenX >= t.theme.BoardX+t.theme.BoardSize {
		return SquareTarget{}, false
	}

	if screenY < t.theme.BoardY || screenY >= t.theme.BoardY+t.theme.BoardSize {
		return SquareTarget{}, false
	}

	file := (screenX - t.theme.BoardX) / t.theme.SquareSize
	rank := 7 - (screenY-t.theme.BoardY)/t.theme.SquareSize

	return SquareTarget{
		File: file,
		Rank: rank,
	}, true
}

func (t *Translator) PromotionChoiceAt(screenX, screenY int, pieceTypes []string) (string, bool) {
	rects := PromotionOptionRects(t.theme, len(pieceTypes))
	for index, rect := range rects {
		if rect.Contains(screenX, screenY) {
			return pieceTypes[index], true
		}
	}

	return "", false
}

func PromotionOptionRects(theme theme.Theme, optionCount int) []Rect {
	if optionCount <= 0 {
		return nil
	}

	optionSize := theme.SquareSize
	totalWidth := optionCount * optionSize
	startX := theme.BoardX + (theme.BoardSize-totalWidth)/2
	startY := theme.BoardY + (theme.BoardSize-optionSize)/2

	rects := make([]Rect, 0, optionCount)
	for index := 0; index < optionCount; index++ {
		rects = append(rects, Rect{
			X:      startX + index*optionSize,
			Y:      startY,
			Width:  optionSize,
			Height: optionSize,
		})
	}

	return rects
}

func (s SquareTarget) Algebraic() string {
	return fmt.Sprintf("%c%d", 'a'+s.File, s.Rank+1)
}

func (r Rect) Contains(x, y int) bool {
	return x >= r.X && x < r.X+r.Width && y >= r.Y && y < r.Y+r.Height
}
