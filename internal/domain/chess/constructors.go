package chess

import (
	chessgame "github.com/KaiserSin/go-chess-ai/internal/domain/chess/game"
	chessmodel "github.com/KaiserSin/go-chess-ai/internal/domain/chess/model"
	chessposition "github.com/KaiserSin/go-chess-ai/internal/domain/chess/position"
)

func NewGame() *Game {
	return chessgame.NewGame()
}

func NewGameFromPosition(position Position) (*Game, error) {
	return chessgame.NewGameFromPosition(position)
}

func NewInitialPosition() Position {
	return chessposition.NewInitialPosition()
}

func NewPositionBuilder() *PositionBuilder {
	return chessposition.NewPositionBuilder()
}

func NewSquare(file, rank int) (Square, error) {
	return chessmodel.NewSquare(file, rank)
}

func ParseSquare(raw string) (Square, error) {
	return chessmodel.ParseSquare(raw)
}
