package theme

import "image/color"

type Theme struct {
	WindowWidth      int
	WindowHeight     int
	BoardX           int
	BoardY           int
	BoardSize        int
	SquareSize       int
	BackgroundColor  color.RGBA
	BorderColor      color.RGBA
	LightSquareColor color.RGBA
	DarkSquareColor  color.RGBA
	TitleColor       color.RGBA
	StatusColor      color.RGBA
	LabelColor       color.RGBA
	WhitePieceColor  color.RGBA
	BlackPieceColor  color.RGBA
	PieceCatalog     PieceCatalog
}

func NewTheme() Theme {
	whitePieceColor := color.RGBA{R: 248, G: 245, B: 240, A: 255}
	blackPieceColor := color.RGBA{R: 29, G: 29, B: 29, A: 255}

	return Theme{
		WindowWidth:      720,
		WindowHeight:     760,
		BoardX:           40,
		BoardY:           60,
		BoardSize:        640,
		SquareSize:       80,
		BackgroundColor:  color.RGBA{R: 22, G: 22, B: 22, A: 255},
		BorderColor:      color.RGBA{R: 58, G: 42, B: 27, A: 255},
		LightSquareColor: color.RGBA{R: 240, G: 217, B: 181, A: 255},
		DarkSquareColor:  color.RGBA{R: 181, G: 136, B: 99, A: 255},
		TitleColor:       color.RGBA{R: 248, G: 245, B: 240, A: 255},
		StatusColor:      color.RGBA{R: 226, G: 214, B: 197, A: 255},
		LabelColor:       color.RGBA{R: 218, G: 204, B: 185, A: 255},
		WhitePieceColor:  whitePieceColor,
		BlackPieceColor:  blackPieceColor,
		PieceCatalog:     NewPieceCatalog(whitePieceColor, blackPieceColor),
	}
}
