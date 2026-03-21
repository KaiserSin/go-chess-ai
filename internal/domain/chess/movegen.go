package chess

func (p Position) generateLegalMoves() []Move {
	pseudoMoves := p.generatePseudoLegalMoves()
	legalMoves := make([]Move, 0, len(pseudoMoves))

	for _, move := range pseudoMoves {
		next := p.applyMoveUnchecked(move)
		if next.isInCheck(p.sideToMove) {
			continue
		}

		legalMoves = append(legalMoves, move)
	}

	return legalMoves
}

func (p Position) generatePseudoLegalMoves() []Move {
	moves := make([]Move, 0, 64)
	side := p.sideToMove

	for pieceType := Pawn; pieceType <= King; pieceType++ {
		forEachSquare(p.board.bitboard(side, pieceType), func(from Square) {
			switch pieceType {
			case Pawn:
				p.appendPawnMoves(&moves, from, side)
			case Knight:
				p.appendKnightMoves(&moves, from, side)
			case Bishop:
				p.appendSlidingMoves(&moves, from, side, bishopDirections)
			case Rook:
				p.appendSlidingMoves(&moves, from, side, rookDirections)
			case Queen:
				p.appendSlidingMoves(&moves, from, side, queenDirections)
			case King:
				p.appendKingMoves(&moves, from, side)
			}
		})
	}

	return moves
}

func (p Position) appendPawnMoves(moves *[]Move, from Square, side Side) {
	direction := pawnDirection(side)
	startRank := pawnStartRank(side)
	promotionRank := pawnPromotionRank(side)

	oneForward, ok := from.offset(0, direction)
	if ok && !p.board.occupied(oneForward) {
		p.appendPawnMove(moves, from, oneForward, promotionRank)

		twoForward, canDouble := from.offset(0, direction*2)
		if from.Rank() == startRank && canDouble && !p.board.occupied(twoForward) {
			*moves = append(*moves, Move{From: from, To: twoForward})
		}
	}

	for _, fileDelta := range pawnCaptureFiles {
		target, canCapture := from.offset(fileDelta, direction)
		if !canCapture {
			continue
		}

		if p.board.occupiedBy(side.Opponent(), target) {
			p.appendPawnMove(moves, from, target, promotionRank)
			continue
		}

		if p.enPassant.ok && target == p.enPassant.value {
			*moves = append(*moves, Move{From: from, To: target})
		}
	}
}

func (p Position) appendPawnMove(moves *[]Move, from, to Square, promotionRank int) {
	if to.Rank() != promotionRank {
		*moves = append(*moves, Move{From: from, To: to})
		return
	}

	for _, promotion := range promotionChoices {
		*moves = append(*moves, Move{
			From:      from,
			To:        to,
			Promotion: promotion,
		})
	}
}

func (p Position) appendKnightMoves(moves *[]Move, from Square, side Side) {
	for _, offset := range knightOffsets {
		target, ok := from.offset(offset[0], offset[1])
		if !ok || p.board.occupiedBy(side, target) {
			continue
		}

		*moves = append(*moves, Move{From: from, To: target})
	}
}

func (p Position) appendSlidingMoves(moves *[]Move, from Square, side Side, directions [][2]int) {
	for _, direction := range directions {
		current := from
		for {
			next, ok := current.offset(direction[0], direction[1])
			if !ok {
				break
			}

			if p.board.occupiedBy(side, next) {
				break
			}

			*moves = append(*moves, Move{From: from, To: next})
			if p.board.occupiedBy(side.Opponent(), next) {
				break
			}

			current = next
		}
	}
}

func (p Position) appendKingMoves(moves *[]Move, from Square, side Side) {
	for _, offset := range kingOffsets {
		target, ok := from.offset(offset[0], offset[1])
		if !ok || p.board.occupiedBy(side, target) {
			continue
		}

		*moves = append(*moves, Move{From: from, To: target})
	}

	if p.canCastleKingside(side) {
		*moves = append(*moves, Move{From: from, To: castleKingSquare(side, kingSide)})
	}

	if p.canCastleQueenside(side) {
		*moves = append(*moves, Move{From: from, To: castleKingSquare(side, queenSide)})
	}
}

func (p Position) canCastleKingside(side Side) bool {
	return p.canCastle(side, kingSide)
}

func (p Position) canCastleQueenside(side Side) bool {
	return p.canCastle(side, queenSide)
}

func (p Position) canCastle(side Side, castle castleSide) bool {
	if castle == kingSide && !p.castlingRights.CanCastleKingside(side) {
		return false
	}

	if castle == queenSide && !p.castlingRights.CanCastleQueenside(side) {
		return false
	}

	if !p.hasCastlingPieces(side, castle) {
		return false
	}

	for _, square := range castleTravelSquares(side, castle) {
		if p.board.occupied(square) {
			return false
		}
	}

	opponent := side.Opponent()
	for _, square := range castleSafeSquares(side, castle) {
		if p.isSquareAttacked(square, opponent) {
			return false
		}
	}

	return true
}

func (p Position) hasCastlingPieces(side Side, castle castleSide) bool {
	king, kingPresent := p.board.pieceAt(kingStartSquare(side))
	rook, rookPresent := p.board.pieceAt(rookStartSquare(side, castle))
	if !kingPresent || !rookPresent {
		return false
	}

	return king.kind == King && king.side == side && rook.kind == Rook && rook.side == side
}

func (p Position) isSquareAttacked(square Square, by Side) bool {
	for _, pawnFileDelta := range pawnCaptureFiles {
		source, ok := square.offset(pawnFileDelta, -pawnDirection(by))
		if !ok {
			continue
		}

		if p.hasPiece(source, by, Pawn) {
			return true
		}
	}

	for _, offset := range knightOffsets {
		source, ok := square.offset(offset[0], offset[1])
		if !ok {
			continue
		}

		if p.hasPiece(source, by, Knight) {
			return true
		}
	}

	if p.attackedBySliding(square, by, rookDirections, Rook, Queen) {
		return true
	}

	if p.attackedBySliding(square, by, bishopDirections, Bishop, Queen) {
		return true
	}

	for _, offset := range kingOffsets {
		source, ok := square.offset(offset[0], offset[1])
		if !ok {
			continue
		}

		if p.hasPiece(source, by, King) {
			return true
		}
	}

	return false
}

func (p Position) attackedBySliding(square Square, by Side, directions [][2]int, matchA, matchB PieceType) bool {
	for _, direction := range directions {
		current := square
		for {
			next, ok := current.offset(direction[0], direction[1])
			if !ok {
				break
			}

			piece, found := p.board.pieceAt(next)
			if !found {
				current = next
				continue
			}

			if piece.side == by && (piece.kind == matchA || piece.kind == matchB) {
				return true
			}

			break
		}
	}

	return false
}

func (p Position) hasPiece(square Square, side Side, kind PieceType) bool {
	piece, found := p.board.pieceAt(square)
	return found && piece.side == side && piece.kind == kind
}
