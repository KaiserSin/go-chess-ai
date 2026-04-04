package scene

import (
	"testing"

	"github.com/KaiserSin/go-chess-ai/internal/application/dto"
	"github.com/KaiserSin/go-chess-ai/internal/application/gameplay"
	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/assets"
	boardinput "github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/input"
	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/theme"
	"github.com/KaiserSin/go-chess-ai/internal/presentation/ebiten/viewmodel"
)

func TestNewGameStartsOnMenuScreen(t *testing.T) {
	game := newTestGame(t)

	if game.screen != menuScreen {
		t.Fatalf("want menu screen, got %d", game.screen)
	}
}

func TestWhiteSelectionStartsPlayableGame(t *testing.T) {
	game := newTestGame(t)
	clickSideChoice(t, game, "white")

	if game.screen != playingScreen {
		t.Fatalf("want playing screen, got %d", game.screen)
	}

	if game.playerSide != "white" {
		t.Fatalf("want player side white, got %s", game.playerSide)
	}

	if got := game.board.Status; got != "white to move" {
		t.Fatalf("want white turn status, got %q", got)
	}

	x, y := squareCenter(game.theme, false, 4, 1)
	game.handleClick(x, y)

	if !squareSelected(game.service.Snapshot(), "e2") {
		t.Fatal("want e2 selected after white move input")
	}
}

func TestBlackSelectionBlocksBoardInputUntilAIExists(t *testing.T) {
	game := newTestGame(t)
	clickSideChoice(t, game, "black")

	if game.playerSide != "black" {
		t.Fatalf("want player side black, got %s", game.playerSide)
	}

	if got := game.board.Status; got != "white to move · AI move not implemented yet" {
		t.Fatalf("want ai placeholder status, got %q", got)
	}

	x, y := squareCenter(game.theme, true, 4, 6)
	game.handleClick(x, y)

	if hasSelectedSquare(game.service.Snapshot()) {
		t.Fatal("did not expect selection while ai turn is blocked")
	}
}

func newTestGame(t *testing.T) *Game {
	t.Helper()

	uiTheme := theme.NewTheme()
	game, err := NewGame(
		gameplay.NewService(),
		viewmodel.NewMapper(uiTheme),
		boardinput.NewTranslator(uiTheme),
		uiTheme,
		assets.PieceSprites{},
	)
	if err != nil {
		t.Fatalf("new game: %v", err)
	}

	return game
}

func clickSideChoice(t *testing.T, game *Game, side string) {
	t.Helper()

	for _, choice := range boardinput.SideChoiceRects(game.theme.WindowWidth) {
		if choice.Side != side {
			continue
		}

		game.handleMenuClick(choice.Rect.X+choice.Rect.Width/2, choice.Rect.Y+choice.Rect.Height/2)
		return
	}

	t.Fatalf("side choice %q not found", side)
}

func squareCenter(uiTheme theme.Theme, blackPerspective bool, file, rank int) (int, int) {
	column := file
	row := 7 - rank
	if blackPerspective {
		column = 7 - file
		row = rank
	}

	return uiTheme.BoardX + column*uiTheme.SquareSize + uiTheme.SquareSize/2,
		uiTheme.BoardY + row*uiTheme.SquareSize + uiTheme.SquareSize/2
}

func squareSelected(snapshot dto.GameSnapshot, algebraic string) bool {
	for _, square := range snapshot.Squares {
		if square.Algebraic == algebraic {
			return square.Selected
		}
	}

	return false
}

func hasSelectedSquare(snapshot dto.GameSnapshot) bool {
	for _, square := range snapshot.Squares {
		if square.Selected {
			return true
		}
	}

	return false
}
