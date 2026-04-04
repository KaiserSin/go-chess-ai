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

	if game.aiSearchDepth != defaultAISearchDepth {
		t.Fatalf("want default ai depth %d, got %d", defaultAISearchDepth, game.aiSearchDepth)
	}

	if game.depthInput != "" {
		t.Fatalf("want empty depth input, got %q", game.depthInput)
	}

	if !game.depthFocused {
		t.Fatal("want depth input focused on menu start")
	}
}

func TestDepthInputUpdatesMenuChoice(t *testing.T) {
	game := newTestGame(t)
	focusDepthInput(t, game)
	game.appendDepthDigit('1')
	game.appendDepthDigit('2')

	if game.aiSearchDepth != 12 {
		t.Fatalf("want ai depth 12, got %d", game.aiSearchDepth)
	}

	if game.depthInput != "12" {
		t.Fatalf("want depth input 12, got %q", game.depthInput)
	}
}

func TestDepthInputStaysEmptyAfterBlur(t *testing.T) {
	game := newTestGame(t)
	focusDepthInput(t, game)
	game.deleteDepthDigit()
	game.blurDepthInput()

	if game.depthInput != "" {
		t.Fatalf("want empty depth input after blur, got %q", game.depthInput)
	}
}

func TestDepthInputRejectsValuesAboveTwenty(t *testing.T) {
	game := newTestGame(t)
	focusDepthInput(t, game)
	game.appendDepthDigit('2')
	game.appendDepthDigit('0')

	if game.aiSearchDepth != 20 {
		t.Fatalf("want ai depth 20, got %d", game.aiSearchDepth)
	}

	game.appendDepthDigit('1')

	if game.aiSearchDepth != 20 {
		t.Fatalf("want ai depth to stay 20, got %d", game.aiSearchDepth)
	}

	if game.depthInput != "20" {
		t.Fatalf("want depth input 20, got %q", game.depthInput)
	}
}

func TestEmptyDepthInputUsesDefaultDepth(t *testing.T) {
	game := newTestGame(t)
	clickSideChoice(t, game, "white")

	if game.aiSearchDepth != defaultAISearchDepth {
		t.Fatalf("want default ai depth %d, got %d", defaultAISearchDepth, game.aiSearchDepth)
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

func TestBlackSelectionTriggersAutomaticAIMove(t *testing.T) {
	game := newTestGame(t)
	clickSideChoice(t, game, "black")

	if game.playerSide != "black" {
		t.Fatalf("want player side black, got %s", game.playerSide)
	}

	if got := game.board.Status; got != "white to move" {
		t.Fatalf("want white turn status before ai move, got %q", got)
	}

	if err := game.updateGame(); err != nil {
		t.Fatalf("updateGame error: %v", err)
	}

	if got := game.board.Status; got != "black to move" {
		t.Fatalf("want black turn status after ai move, got %q", got)
	}

	snapshot := game.service.Snapshot()
	if square := squareByAlgebraic(t, snapshot, "a3"); square.PieceKey != "white-pawn" {
		t.Fatalf("want white-pawn on a3, got %q", square.PieceKey)
	}

	if square := squareByAlgebraic(t, snapshot, "a2"); square.Occupied {
		t.Fatal("did not expect piece on a2 after ai move")
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

func focusDepthInput(t *testing.T, game *Game) {
	t.Helper()

	rect := boardinput.DepthInputRect(game.theme.WindowWidth)
	game.handleMenuClick(rect.X+rect.Width/2, rect.Y+rect.Height/2)
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

func squareByAlgebraic(t *testing.T, snapshot dto.GameSnapshot, algebraic string) dto.SquareSnapshot {
	t.Helper()

	for _, square := range snapshot.Squares {
		if square.Algebraic == algebraic {
			return square
		}
	}

	t.Fatalf("square %q not found", algebraic)
	return dto.SquareSnapshot{}
}
