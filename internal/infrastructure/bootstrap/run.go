package bootstrap

import (
	"github.com/KaiserSin/go-chess-ai/internal/application/gameplay"
	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/assets"
	boardinput "github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/input"
	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/scene"
	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/theme"
	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/viewmodel"
	"github.com/hajimehoshi/ebiten/v2"
)

func Run() error {
	uiTheme := theme.NewTheme()
	service := gameplay.NewService()
	mapper := viewmodel.NewMapper(uiTheme)
	input := boardinput.NewTranslator(uiTheme)
	sprites := assets.LoadPieceSprites()
	game, err := scene.NewGame(service, mapper, input, uiTheme, sprites)
	if err != nil {
		return err
	}

	ebiten.SetWindowSize(uiTheme.WindowWidth, uiTheme.WindowHeight)
	ebiten.SetWindowTitle("Go Chess AI")

	return ebiten.RunGame(game)
}
