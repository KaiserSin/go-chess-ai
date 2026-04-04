package ai

import chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"

const mateScore = 100000

func Evaluate(position chess.Position, perspective chess.Side) int {
	if perspective != chess.White && perspective != chess.Black {
		return 0
	}

	switch position.Status() {
	case chess.Checkmate:
		if perspective == position.SideToMove() {
			return -mateScore
		}

		return mateScore
	case chess.Stalemate:
		return 0
	}

	if chess.HasInsufficientMaterial(position) || chess.IsFiftyMoveDraw(position) {
		return 0
	}

	score := materialScore(position)
	if perspective == chess.Black {
		return -score
	}

	return score
}

func materialScore(position chess.Position) int {
	score := 0

	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			square := squareMust(file, rank)
			piece, ok := position.PieceAt(square)
			if !ok {
				continue
			}

			value := pieceValue(piece.Type())
			if piece.Side() == chess.White {
				score += value
			} else {
				score -= value
			}
		}
	}

	return score
}

func pieceValue(pieceType chess.PieceType) int {
	switch pieceType {
	case chess.Pawn:
		return 100
	case chess.Knight:
		return 320
	case chess.Bishop:
		return 330
	case chess.Rook:
		return 500
	case chess.Queen:
		return 900
	default:
		return 0
	}
}

func squareMust(file, rank int) chess.Square {
	square, err := chess.NewSquare(file, rank)
	if err != nil {
		panic(err)
	}

	return square
}
