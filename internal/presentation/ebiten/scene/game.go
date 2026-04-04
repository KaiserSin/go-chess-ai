package scene

import (
	"bytes"
	"image/color"
	"strconv"

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

type screenState uint8

const (
	menuScreen screenState = iota
	playingScreen
)

const (
	defaultAISearchDepth = 2
	menuTitleY           = 150
	menuDepthLabelY      = 220
	menuDepthHintY       = 328
	menuSideLabelY       = 356
	menuBorderWidth      = 2
	maxAISearchDepth     = 20
	menuInputPaddingX    = 14
	menuInputPaddingY    = 13
)

type Game struct {
	service       *gameplay.Service
	mapper        *viewmodel.Mapper
	input         *boardinput.Translator
	sprites       assets.PieceSprites
	theme         theme.Theme
	titleFace     text.Face
	statusFace    text.Face
	labelFace     text.Face
	pieceFace     text.Face
	screen        screenState
	playerSide    string
	aiSearchDepth int
	depthInput    string
	depthFocused  bool
	board         viewmodel.BoardViewModel
}

func NewGame(service *gameplay.Service, mapper *viewmodel.Mapper, input *boardinput.Translator, uiTheme theme.Theme, sprites assets.PieceSprites) (*Game, error) {
	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		return nil, err
	}

	return &Game{
		service:       service,
		mapper:        mapper,
		input:         input,
		sprites:       sprites,
		theme:         uiTheme,
		screen:        menuScreen,
		aiSearchDepth: defaultAISearchDepth,
		depthInput:    "",
		depthFocused:  true,
		titleFace:     &text.GoTextFace{Source: fontSource, Size: 24},
		statusFace:    &text.GoTextFace{Source: fontSource, Size: 18},
		labelFace:     &text.GoTextFace{Source: fontSource, Size: 18},
		pieceFace:     &text.GoTextFace{Source: fontSource, Size: 42},
	}, nil
}

func (g *Game) Update() error {
	if g.screen == menuScreen {
		return g.updateMenu()
	}

	return g.updateGame()
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(g.theme.BackgroundColor)

	if g.screen == menuScreen {
		g.drawMenu(screen)
		return
	}

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

func (g *Game) drawMenu(screen *ebiten.Image) {
	drawCenteredText(screen, "Go Chess AI", g.titleFace, g.theme.WindowWidth/2, menuTitleY, g.theme.TitleColor)
	drawCenteredText(screen, "Choose depth", g.statusFace, g.theme.WindowWidth/2, menuDepthLabelY, g.theme.StatusColor)

	depthRect := boardinput.DepthInputRect(g.theme.WindowWidth)
	vector.FillRect(
		screen,
		float32(depthRect.X),
		float32(depthRect.Y),
		float32(depthRect.Width),
		float32(depthRect.Height),
		color.RGBA{R: 28, G: 28, B: 28, A: 255},
		false,
	)

	depthBorderColor := g.theme.BorderColor
	if g.depthFocused {
		depthBorderColor = g.theme.SelectedSquareColor
	}
	drawRectBorder(
		screen,
		depthRect.X,
		depthRect.Y,
		depthRect.Width,
		depthRect.Height,
		menuBorderWidth,
		depthBorderColor,
	)

	drawTopLeftText(
		screen,
		g.depthFieldDisplay(),
		g.labelFace,
		depthRect.X+menuInputPaddingX,
		depthRect.Y+menuInputPaddingY,
		g.theme.LabelColor,
	)
	drawCenteredText(screen, "Type a number from 1 to 20", g.labelFace, g.theme.WindowWidth/2, menuDepthHintY, g.theme.LabelColor)

	drawCenteredText(screen, "Choose your side", g.statusFace, g.theme.WindowWidth/2, menuSideLabelY, g.theme.StatusColor)

	for _, choice := range boardinput.SideChoiceRects(g.theme.WindowWidth) {
		drawRectBorder(
			screen,
			choice.Rect.X,
			choice.Rect.Y,
			choice.Rect.Width,
			choice.Rect.Height,
			menuBorderWidth,
			g.theme.BorderColor,
		)
		drawCenteredText(
			screen,
			choice.Label,
			g.labelFace,
			choice.Rect.X+choice.Rect.Width/2,
			choice.Rect.Y+choice.Rect.Height/2,
			g.theme.LabelColor,
		)
	}
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
	g.drawAxisLabelSet(screen, g.board.FileLabels)
	g.drawAxisLabelSet(screen, g.board.RankLabels)
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
		g.drawPromotionOption(screen, option)
	}
}

func (g *Game) drawAxisLabelSet(screen *ebiten.Image, labels []viewmodel.AxisLabelViewModel) {
	for _, label := range labels {
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

func (g *Game) drawPiece(screen *ebiten.Image, visual theme.PieceVisual, x, y, size int) {
	if sprite, ok := g.sprites.Lookup(visual.AssetKey); ok {
		g.drawSprite(screen, sprite, x, y, size)
		return
	}

	centerX := x + size/2
	centerY := y + size/2
	drawCenteredText(screen, visual.Label, g.pieceFace, centerX+1, centerY+1, color.RGBA{A: 64})
	drawCenteredText(screen, visual.Label, g.pieceFace, centerX, centerY, visual.Color)
}

func (g *Game) drawSprite(screen, sprite *ebiten.Image, x, y, size int) {
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
}

type spritePlacement struct {
	X      float64
	Y      float64
	ScaleX float64
	ScaleY float64
}

type targetDotPlacement struct {
	CenterX float32
	CenterY float32
	Radius  float32
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
		dot := targetDotForSquare(x, y, square.Size)
		vector.FillCircle(screen, dot.CenterX, dot.CenterY, dot.Radius, g.theme.SelectedSquareColor, false)
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

func targetDotForSquare(squareX, squareY, squareSize int) targetDotPlacement {
	return targetDotPlacement{
		CenterX: float32(squareX + squareSize/2),
		CenterY: float32(squareY + squareSize/2),
		Radius:  float32(squareSize) / 8,
	}
}

func (g *Game) refreshBoard() {
	if g.screen != playingScreen {
		return
	}

	snapshot := g.service.Snapshot()
	g.board = g.mapper.Map(snapshot, g.blackPerspective())
}

func (g *Game) drawPromotionOption(screen *ebiten.Image, option viewmodel.PromotionOptionViewModel) {
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

func (g *Game) handleClick(screenX, screenY int) {
	snapshot := g.service.Snapshot()
	if g.isAITurn(snapshot.SideToMove, snapshot.OutcomeReason) {
		return
	}

	if g.board.Promotion != nil {
		g.handlePromotionClick(screenX, screenY)
		return
	}

	if square, ok := g.input.SquareAt(screenX, screenY, g.blackPerspective()); ok {
		g.service.SelectSquareAt(square.File, square.Rank)
	}
}

func (g *Game) handleMenuClick(screenX, screenY int) error {
	if g.input.DepthInputAt(screenX, screenY) {
		g.focusDepthInput()
		return nil
	}

	g.blurDepthInput()

	side, ok := g.input.SideChoiceAt(screenX, screenY)
	if !ok {
		return nil
	}

	return g.startGameAs(side)
}

func (g *Game) handlePromotionClick(screenX, screenY int) {
	choice, ok := g.input.PromotionChoiceAt(screenX, screenY, g.promotionChoices())
	if !ok {
		return
	}

	if err := g.service.ChoosePromotionByName(choice); err != nil {
		return
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

func (g *Game) startGameAs(side string) error {
	if side != "black" {
		side = "white"
	}

	g.aiSearchDepth = g.menuSearchDepth()
	g.playerSide = side
	g.service.SetAISearchDepth(g.aiSearchDepth)
	g.service.NewGame()
	g.screen = playingScreen
	g.refreshBoard()

	return g.advanceAITurnIfNeeded()
}

func (g *Game) blackPerspective() bool {
	return g.playerSide == "black"
}

func (g *Game) isAITurn(sideToMove, outcomeReason string) bool {
	if g.screen != playingScreen || g.playerSide == "" || gameFinished(outcomeReason) {
		return false
	}

	return sideToMove == opponentSide(g.playerSide)
}

func (g *Game) updateMenu() error {
	g.handleMenuKeyboardInput()

	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return nil
	}

	mouseX, mouseY := ebiten.CursorPosition()
	return g.handleMenuClick(mouseX, mouseY)
}

func (g *Game) advanceAITurnIfNeeded() error {
	snapshot := g.service.Snapshot()
	if !g.isAITurn(snapshot.SideToMove, snapshot.OutcomeReason) {
		return nil
	}

	if err := g.service.ApplyAIMove(); err != nil {
		return err
	}

	g.refreshBoard()
	return nil
}

func (g *Game) handleMenuKeyboardInput() {
	if !g.depthFocused {
		return
	}

	if digit, ok := menuDigitKeyJustPressed(); ok {
		g.appendDepthDigit(digit)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		g.deleteDepthDigit()
	}
}

func (g *Game) appendDepthDigit(inputRune rune) {
	if inputRune < '0' || inputRune > '9' {
		return
	}

	next := g.depthInput + string(inputRune)
	depth, ok := parseDepthInput(next)
	if !ok {
		return
	}

	g.depthInput = next
	g.aiSearchDepth = depth
}

func (g *Game) deleteDepthDigit() {
	if g.depthInput == "" {
		return
	}

	g.depthInput = g.depthInput[:len(g.depthInput)-1]

	if depth, ok := parseDepthInput(g.depthInput); ok {
		g.aiSearchDepth = depth
	}
}

func (g *Game) menuSearchDepth() int {
	depth, ok := parseDepthInput(g.depthInput)
	if !ok {
		return defaultAISearchDepth
	}

	return depth
}

func (g *Game) updateGame() error {
	g.refreshBoard()
	if err := g.advanceAITurnIfNeeded(); err != nil {
		return err
	}

	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return nil
	}

	mouseX, mouseY := ebiten.CursorPosition()
	g.handleClick(mouseX, mouseY)
	g.refreshBoard()

	return g.advanceAITurnIfNeeded()
}

func opponentSide(side string) string {
	if side == "black" {
		return "white"
	}

	return "black"
}

func gameFinished(outcomeReason string) bool {
	return outcomeReason != "" && outcomeReason != "none"
}

func parseDepthInput(value string) (int, bool) {
	if value == "" {
		return 0, false
	}

	depth, err := strconv.Atoi(value)
	if err != nil || depth < 1 || depth > maxAISearchDepth {
		return 0, false
	}

	return depth, true
}

func (g *Game) focusDepthInput() {
	g.depthFocused = true
}

func (g *Game) blurDepthInput() {
	g.depthFocused = false
}

func (g *Game) depthFieldDisplay() string {
	value := g.depthInput
	if !g.depthFocused || !caretVisible() {
		return value
	}

	return value + "|"
}

func caretVisible() bool {
	return ebiten.Tick()/30%2 == 0
}

func menuDigitKeyJustPressed() (rune, bool) {
	for digit := 0; digit <= 9; digit++ {
		if inpututil.IsKeyJustPressed(ebiten.KeyDigit0 + ebiten.Key(digit)) {
			return rune('0' + digit), true
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyNumpad0 + ebiten.Key(digit)) {
			return rune('0' + digit), true
		}
	}

	return 0, false
}
