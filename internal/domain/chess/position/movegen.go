package position

import (
	"github.com/KaiserSin/go-chess-ai/internal/domain/chess/internal/geom"
	chessmodel "github.com/KaiserSin/go-chess-ai/internal/domain/chess/model"
)

func (p Position) generateLegalMoves() []chessmodel.Move {
	pseudoMoves := p.generatePseudoLegalMoves()
	legalMoves := make([]chessmodel.Move, 0, len(pseudoMoves))

	for _, move := range pseudoMoves {
		next := p.applyMoveUnchecked(move)
		if next.isInCheck(p.sideToMove) {
			continue
		}

		legalMoves = append(legalMoves, move)
	}

	return legalMoves
}

func (p Position) generatePseudoLegalMoves() []chessmodel.Move {
	moves := make([]chessmodel.Move, 0, 64)
	side := p.sideToMove

	for pieceType := chessmodel.Pawn; pieceType <= chessmodel.King; pieceType++ {
		forEachSquare(p.board.bitboard(side, pieceType), func(from chessmodel.Square) {
			switch pieceType {
			case chessmodel.Pawn:
				p.appendPawnMoves(&moves, from, side)
			case chessmodel.Knight:
				p.appendKnightMoves(&moves, from, side)
			case chessmodel.Bishop:
				p.appendSlidingMoves(&moves, from, side, geom.BishopDirections)
			case chessmodel.Rook:
				p.appendSlidingMoves(&moves, from, side, geom.RookDirections)
			case chessmodel.Queen:
				p.appendSlidingMoves(&moves, from, side, geom.QueenDirections)
			case chessmodel.King:
				p.appendKingMoves(&moves, from, side)
			}
		})
	}

	return moves
}

func (p Position) appendPawnMoves(moves *[]chessmodel.Move, from chessmodel.Square, side chessmodel.Side) {
	direction := pawnDirection(side)
	startRank := pawnStartRank(side)
	promotionRank := pawnPromotionRank(side)

	oneForward, ok := squareOffset(from, 0, direction)
	if ok && !p.board.occupied(oneForward) {
		p.appendPawnMove(moves, from, oneForward, promotionRank)

		twoForward, canDouble := squareOffset(from, 0, direction*2)
		if from.Rank() == startRank && canDouble && !p.board.occupied(twoForward) {
			*moves = append(*moves, chessmodel.Move{From: from, To: twoForward})
		}
	}

	for _, fileDelta := range geom.PawnCaptureFiles {
		target, canCapture := squareOffset(from, fileDelta, direction)
		if !canCapture {
			continue
		}

		if p.board.occupiedBy(side.Opponent(), target) {
			p.appendPawnMove(moves, from, target, promotionRank)
			continue
		}

		if p.enPassant.ok && target == p.enPassant.value {
			*moves = append(*moves, chessmodel.Move{From: from, To: target})
		}
	}
}

func (p Position) appendPawnMove(moves *[]chessmodel.Move, from, to chessmodel.Square, promotionRank int) {
	if to.Rank() != promotionRank {
		*moves = append(*moves, chessmodel.Move{From: from, To: to})
		return
	}

	for _, promotion := range promotionChoices {
		*moves = append(*moves, chessmodel.Move{
			From:      from,
			To:        to,
			Promotion: promotion,
		})
	}
}

func (p Position) appendKnightMoves(moves *[]chessmodel.Move, from chessmodel.Square, side chessmodel.Side) {
	for _, offset := range geom.KnightOffsets {
		target, ok := squareOffset(from, offset[0], offset[1])
		if !ok || p.board.occupiedBy(side, target) {
			continue
		}

		*moves = append(*moves, chessmodel.Move{From: from, To: target})
	}
}

func (p Position) appendSlidingMoves(moves *[]chessmodel.Move, from chessmodel.Square, side chessmodel.Side, directions [][2]int) {
	for _, direction := range directions {
		current := from
		for {
			next, ok := squareOffset(current, direction[0], direction[1])
			if !ok {
				break
			}

			if p.board.occupiedBy(side, next) {
				break
			}

			*moves = append(*moves, chessmodel.Move{From: from, To: next})
			if p.board.occupiedBy(side.Opponent(), next) {
				break
			}

			current = next
		}
	}
}

func (p Position) appendKingMoves(moves *[]chessmodel.Move, from chessmodel.Square, side chessmodel.Side) {
	for _, offset := range geom.KingOffsets {
		target, ok := squareOffset(from, offset[0], offset[1])
		if !ok || p.board.occupiedBy(side, target) {
			continue
		}

		*moves = append(*moves, chessmodel.Move{From: from, To: target})
	}

	if p.canCastleKingside(side) {
		*moves = append(*moves, chessmodel.Move{From: from, To: castleKingSquare(side, kingSide)})
	}

	if p.canCastleQueenside(side) {
		*moves = append(*moves, chessmodel.Move{From: from, To: castleKingSquare(side, queenSide)})
	}
}
