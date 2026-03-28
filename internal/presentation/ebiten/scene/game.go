package scene

import (
	"bytes"
	"image/color"

	"github.com/KaiserSin/go-chess-ai/internal/application/gameplay"
	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/assets"
	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/theme"
	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/viewmodel"
	"github.com/hajimehoshi/ebiten/v2"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"
)

type Game struct {
	service    *gameplay.Service
	mapper     *viewmodel.Mapper
	sprites    assets.PieceSprites
	theme      theme.Theme
	titleFace  text.Face
	statusFace text.Face
	labelFace  text.Face
	pieceFace  text.Face
	board      viewmodel.BoardViewModel
}

func NewGame(service *gameplay.Service, mapper *viewmodel.Mapper, uiTheme theme.Theme, sprites assets.PieceSprites) (*Game, error) {
	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		return nil, err
	}

	return &Game{
		service:    service,
		mapper:     mapper,
		sprites:    sprites,
		theme:      uiTheme,
		titleFace:  &text.GoTextFace{Source: fontSource, Size: 24},
		statusFace: &text.GoTextFace{Source: fontSource, Size: 18},
		labelFace:  &text.GoTextFace{Source: fontSource, Size: 18},
		pieceFace:  &text.GoTextFace{Source: fontSource, Size: 42},
	}, nil
}

func (g *Game) Update() error {
	g.board = g.mapper.Map(g.service.Snapshot())
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.theme.BackgroundColor)

	g.drawBoard(screen)
	g.drawAxisLabels(screen)
	g.drawHeader(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.theme.WindowWidth, g.theme.WindowHeight
}

func (g *Game) drawHeader(screen *ebiten.Image) {
	drawTopLeftText(screen, g.board.Title, g.titleFace, g.theme.BoardX, 18, g.theme.TitleColor)
	drawTopLeftText(screen, g.board.Status, g.statusFace, g.theme.BoardX, 50, g.theme.StatusColor)
}

func (g *Game) drawBoard(screen *ebiten.Image) {
	borderInset := float32(4)
	vector.FillRect(
		screen,
		float32(g.board.BoardX)-borderInset,
		float32(g.board.BoardY)-borderInset,
		float32(g.board.BoardSize)+borderInset*2,
		float32(g.board.BoardSize)+borderInset*2,
		g.theme.BorderColor,
		false,
	)

	for _, square := range g.board.Squares {
		squareColor := g.theme.DarkSquareColor
		if square.IsLight {
			squareColor = g.theme.LightSquareColor
		}

		vector.FillRect(
			screen,
			float32(g.board.BoardX+square.X),
			float32(g.board.BoardY+square.Y),
			float32(square.Size),
			float32(square.Size),
			squareColor,
			false,
		)

		if square.Piece.Visible {
			g.drawPiece(screen, square)
		}
	}
}

func (g *Game) drawAxisLabels(screen *ebiten.Image) {
	for _, label := range g.board.FileLabels {
		drawCenteredText(screen, label.Text, g.labelFace, g.board.BoardX+label.CenterX, g.board.BoardY+label.CenterY, g.theme.LabelColor)
	}

	for _, label := range g.board.RankLabels {
		drawCenteredText(screen, label.Text, g.labelFace, g.board.BoardX+label.CenterX, g.board.BoardY+label.CenterY, g.theme.LabelColor)
	}
}

func drawTopLeftText(screen *ebiten.Image, value string, face text.Face, x, y int, clr color.Color) {
	if value == "" {
		return
	}

	var options text.DrawOptions
	options.GeoM.Translate(float64(x), float64(y))
	options.ColorScale.ScaleWithColor(clr)
	text.Draw(screen, value, face, &options)
}

func drawCenteredText(screen *ebiten.Image, value string, face text.Face, centerX, centerY int, clr color.Color) {
	if value == "" {
		return
	}

	width, height := text.Measure(value, face, 0)

	var options text.DrawOptions
	options.GeoM.Translate(float64(centerX)-width/2, float64(centerY)-height/2)
	options.ColorScale.ScaleWithColor(clr)
	text.Draw(screen, value, face, &options)
}

func (g *Game) drawPiece(screen *ebiten.Image, square viewmodel.SquareViewModel) {
	if sprite, ok := g.sprites.Lookup(square.Piece.Visual.AssetKey); ok {
		placement := spritePlacementForSquare(
			g.board.BoardX+square.X,
			g.board.BoardY+square.Y,
			square.Size,
			sprite.Bounds().Dx(),
			sprite.Bounds().Dy(),
		)

		var options ebiten.DrawImageOptions
		options.GeoM.Scale(placement.ScaleX, placement.ScaleY)
		options.GeoM.Translate(placement.X, placement.Y)
		options.Filter = ebiten.FilterLinear
		screen.DrawImage(sprite, &options)
		return
	}

	centerX := g.board.BoardX + square.Piece.CenterX
	centerY := g.board.BoardY + square.Piece.CenterY
	drawCenteredText(screen, square.Piece.Visual.Label, g.pieceFace, centerX+1, centerY+1, color.RGBA{A: 64})
	drawCenteredText(screen, square.Piece.Visual.Label, g.pieceFace, centerX, centerY, square.Piece.Visual.Color)
}

type spritePlacement struct {
	X      float64
	Y      float64
	ScaleX float64
	ScaleY float64
}

func spritePlacementForSquare(squareX, squareY, squareSize, spriteWidth, spriteHeight int) spritePlacement {
	return spritePlacement{
		X:      float64(squareX),
		Y:      float64(squareY),
		ScaleX: float64(squareSize) / float64(spriteWidth),
		ScaleY: float64(squareSize) / float64(spriteHeight),
	}
}
