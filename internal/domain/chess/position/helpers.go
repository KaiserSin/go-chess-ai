package position

import chessmodel "github.com/KaiserSin/go-chess-ai/internal/domain/chess/model"

var (
	allSides         = [...]chessmodel.Side{chessmodel.White, chessmodel.Black}
	promotionChoices = [...]chessmodel.PieceType{
		chessmodel.Queen,
		chessmodel.Rook,
		chessmodel.Bishop,
		chessmodel.Knight,
	}
)

type optionalSquare struct {
	value chessmodel.Square
	ok    bool
}

type castleSide uint8

const (
	kingSide castleSide = iota
	queenSide
)

func noSquare() optionalSquare {
	return optionalSquare{}
}

func someSquare(square chessmodel.Square) optionalSquare {
	return optionalSquare{
		value: square,
		ok:    true,
	}
}

func squareMust(file, rank int) chessmodel.Square {
	square, err := chessmodel.NewSquare(file, rank)
	if err != nil {
		panic(err)
	}

	return square
}

func squareIsValid(square chessmodel.Square) bool {
	return int(square) < 64
}

func squareBitboard(square chessmodel.Square) uint64 {
	return uint64(1) << square
}

func squareOffset(square chessmodel.Square, fileDelta, rankDelta int) (chessmodel.Square, bool) {
	if !squareIsValid(square) {
		return 0, false
	}

	file := square.File() + fileDelta
	rank := square.Rank() + rankDelta
	if file < 0 || file >= 8 || rank < 0 || rank >= 8 {
		return 0, false
	}

	return chessmodel.Square(rank*8 + file), true
}

func squareColor(square chessmodel.Square) int {
	return (square.File() + square.Rank()) % 2
}

func sideIndex(side chessmodel.Side) int {
	return int(side)
}

func isValidSide(side chessmodel.Side) bool {
	return side == chessmodel.White || side == chessmodel.Black
}

func bitboardIndex(pieceType chessmodel.PieceType) int {
	return int(pieceType - chessmodel.Pawn)
}

func isPromotionChoice(pieceType chessmodel.PieceType) bool {
	switch pieceType {
	case chessmodel.Queen, chessmodel.Rook, chessmodel.Bishop, chessmodel.Knight:
		return true
	default:
		return false
	}
}

func validateMoveSquares(move chessmodel.Move) error {
	if !squareIsValid(move.From) || !squareIsValid(move.To) {
		return chessmodel.ErrInvalidSquare
	}

	return nil
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}

	return value
}
