package scene

import (
	"bytes"
	"image/color"

	"github.com/KaiserSin/go-chess-ai/internal/application/gameplay"
	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/assets"
	boardinput "github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/input"
	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/theme"
	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/viewmodel"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	text "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"
)

type Game struct {
	service    *gameplay.Service
	mapper     *viewmodel.Mapper
	input      *boardinput.Translator
	sprites    assets.PieceSprites
	theme      theme.Theme
	titleFace  text.Face
	statusFace text.Face
	labelFace  text.Face
	pieceFace  text.Face
	board      viewmodel.BoardViewModel
}

func NewGame(service *gameplay.Service, mapper *viewmodel.Mapper, input *boardinput.Translator, uiTheme theme.Theme, sprites assets.PieceSprites) (*Game, error) {
	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		return nil, err
	}

	return &Game{
		service:    service,
		mapper:     mapper,
		input:      input,
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

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()
		g.handleClick(mouseX, mouseY)
		g.board = g.mapper.Map(g.service.Snapshot())
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.theme.BackgroundColor)

	g.drawBoard(screen)
	g.drawPromotionOverlay(screen)
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

		g.drawSquareState(screen, square)

		if square.Piece.Visible {
			g.drawPiece(screen, square.Piece.Visual, g.board.BoardX+square.X, g.board.BoardY+square.Y, square.Size)
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

func (g *Game) drawPromotionOverlay(screen *ebiten.Image) {
	if g.board.Promotion == nil {
		return
	}

	vector.FillRect(
		screen,
		float32(g.board.BoardX),
		float32(g.board.BoardY),
		float32(g.board.BoardSize),
		float32(g.board.BoardSize),
		g.theme.PromotionVeilColor,
		false,
	)

	titleY := g.board.Promotion.Options[0].Y - 28
	drawCenteredText(screen, g.board.Promotion.Title, g.statusFace, g.board.BoardX+g.board.BoardSize/2, titleY, g.theme.TitleColor)

	for _, option := range g.board.Promotion.Options {
		vector.FillRect(
			screen,
			float32(option.X),
			float32(option.Y),
			float32(option.Size),
			float32(option.Size),
			g.theme.PromotionButtonColor,
			false,
		)
		drawRectBorder(screen, option.X, option.Y, option.Size, option.Size, 4, g.theme.PromotionButtonBorderColor)
		g.drawPiece(screen, option.Visual, option.X, option.Y, option.Size)
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

func (g *Game) drawPiece(screen *ebiten.Image, visual theme.PieceVisual, x, y, size int) {
	if sprite, ok := g.sprites.Lookup(visual.AssetKey); ok {
		placement := spritePlacementForRect(
			x,
			y,
			size,
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

	centerX := x + size/2
	centerY := y + size/2
	drawCenteredText(screen, visual.Label, g.pieceFace, centerX+1, centerY+1, color.RGBA{A: 64})
	drawCenteredText(screen, visual.Label, g.pieceFace, centerX, centerY, visual.Color)
}

type spritePlacement struct {
	X      float64
	Y      float64
	ScaleX float64
	ScaleY float64
}

func spritePlacementForRect(squareX, squareY, squareSize, spriteWidth, spriteHeight int) spritePlacement {
	return spritePlacement{
		X:      float64(squareX),
		Y:      float64(squareY),
		ScaleX: float64(squareSize) / float64(spriteWidth),
		ScaleY: float64(squareSize) / float64(spriteHeight),
	}
}

func (g *Game) drawSquareState(screen *ebiten.Image, square viewmodel.SquareViewModel) {
	x := g.board.BoardX + square.X
	y := g.board.BoardY + square.Y

	if square.LegalTarget {
		vector.FillRect(
			screen,
			float32(x),
			float32(y),
			float32(square.Size),
			float32(square.Size),
			g.theme.LegalTargetColor,
			false,
		)
	}

	if square.Selected {
		drawRectBorder(screen, x, y, square.Size, square.Size, 4, g.theme.SelectedSquareColor)
	}
}

func drawRectBorder(screen *ebiten.Image, x, y, width, height, thickness int, clr color.Color) {
	vector.FillRect(screen, float32(x), float32(y), float32(width), float32(thickness), clr, false)
	vector.FillRect(screen, float32(x), float32(y+height-thickness), float32(width), float32(thickness), clr, false)
	vector.FillRect(screen, float32(x), float32(y), float32(thickness), float32(height), clr, false)
	vector.FillRect(screen, float32(x+width-thickness), float32(y), float32(thickness), float32(height), clr, false)
}

func (g *Game) handleClick(screenX, screenY int) {
	if g.board.Promotion != nil {
		if choice, ok := g.input.PromotionChoiceAt(screenX, screenY, g.promotionChoices()); ok {
			_ = g.service.ChoosePromotionByName(choice)
		}
		return
	}

	if square, ok := g.input.SquareAt(screenX, screenY); ok {
		g.service.SelectSquareAt(square.File, square.Rank)
	}
}

func (g *Game) promotionChoices() []string {
	if g.board.Promotion == nil {
		return nil
	}

	choices := make([]string, 0, len(g.board.Promotion.Options))
	for _, option := range g.board.Promotion.Options {
		choices = append(choices, option.PieceType)
	}

	return choices
}
