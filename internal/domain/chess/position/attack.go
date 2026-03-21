package position

import (
	"github.com/KaiserSin/go-chess-ai/internal/domain/chess/internal/geom"
	chessmodel "github.com/KaiserSin/go-chess-ai/internal/domain/chess/model"
)

func (p Position) isSquareAttacked(square chessmodel.Square, by chessmodel.Side) bool {
	for _, pawnFileDelta := range geom.PawnCaptureFiles {
		source, ok := squareOffset(square, pawnFileDelta, -pawnDirection(by))
		if !ok {
			continue
		}

		if p.hasPiece(source, by, chessmodel.Pawn) {
			return true
		}
	}

	for _, offset := range geom.KnightOffsets {
		source, ok := squareOffset(square, offset[0], offset[1])
		if !ok {
			continue
		}

		if p.hasPiece(source, by, chessmodel.Knight) {
			return true
		}
	}

	if p.attackedBySliding(square, by, geom.RookDirections, chessmodel.Rook, chessmodel.Queen) {
		return true
	}

	if p.attackedBySliding(square, by, geom.BishopDirections, chessmodel.Bishop, chessmodel.Queen) {
		return true
	}

	for _, offset := range geom.KingOffsets {
		source, ok := squareOffset(square, offset[0], offset[1])
		if !ok {
			continue
		}

		if p.hasPiece(source, by, chessmodel.King) {
			return true
		}
	}

	return false
}

func (p Position) attackedBySliding(square chessmodel.Square, by chessmodel.Side, directions [][2]int, matchA, matchB chessmodel.PieceType) bool {
	for _, direction := range directions {
		current := square
		for {
			next, ok := squareOffset(current, direction[0], direction[1])
			if !ok {
				break
			}

			piece, found := p.board.pieceAt(next)
			if !found {
				current = next
				continue
			}

			if piece.Side() == by && (piece.Type() == matchA || piece.Type() == matchB) {
				return true
			}

			break
		}
	}

	return false
}

func (p Position) hasPiece(square chessmodel.Square, side chessmodel.Side, kind chessmodel.PieceType) bool {
	piece, found := p.board.pieceAt(square)
	return found && piece.Side() == side && piece.Type() == kind
}
