package chess

// generateLegalMoves collects all legal moves for the current side.
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
				p.appendSlidingMoves(&moves, from, side, [][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}})
			case Rook:
				p.appendSlidingMoves(&moves, from, side, [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}})
			case Queen:
				p.appendSlidingMoves(&moves, from, side, [][2]int{
					{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
					{1, 0}, {-1, 0}, {0, 1}, {0, -1},
				})
			case King:
				p.appendKingMoves(&moves, from, side)
			}
		})
	}

	return moves
}

func (p Position) appendPawnMoves(moves *[]Move, from Square, side Side) {
	direction := 1
	startRank := 1
	promotionRank := 7
	if side == Black {
		direction = -1
		startRank = 6
		promotionRank = 0
	}

	oneForward, ok := from.offset(0, direction)
	if ok && !p.board.occupied(oneForward) {
		p.appendPawnMove(moves, from, oneForward, promotionRank)

		twoForward, canDouble := from.offset(0, direction*2)
		if from.Rank() == startRank && canDouble && !p.board.occupied(twoForward) {
			*moves = append(*moves, Move{From: from, To: twoForward})
		}
	}

	for _, fileDelta := range []int{-1, 1} {
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
	offsets := [][2]int{
		{1, 2}, {2, 1}, {2, -1}, {1, -2},
		{-1, -2}, {-2, -1}, {-2, 1}, {-1, 2},
	}

	for _, offset := range offsets {
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
	for rankDelta := -1; rankDelta <= 1; rankDelta++ {
		for fileDelta := -1; fileDelta <= 1; fileDelta++ {
			if fileDelta == 0 && rankDelta == 0 {
				continue
			}

			target, ok := from.offset(fileDelta, rankDelta)
			if !ok || p.board.occupiedBy(side, target) {
				continue
			}

			*moves = append(*moves, Move{From: from, To: target})
		}
	}

	if p.canCastleKingside(side) {
		rank := 0
		if side == Black {
			rank = 7
		}
		*moves = append(*moves, Move{From: from, To: mustSquare(6, rank)})
	}

	if p.canCastleQueenside(side) {
		rank := 0
		if side == Black {
			rank = 7
		}
		*moves = append(*moves, Move{From: from, To: mustSquare(2, rank)})
	}
}

func (p Position) canCastleKingside(side Side) bool {
	if !p.castlingRights.CanCastleKingside(side) {
		return false
	}

	rank := 0
	if side == Black {
		rank = 7
	}

	kingFrom := mustSquare(4, rank)
	rookFrom := mustSquare(7, rank)
	fSquare := mustSquare(5, rank)
	gSquare := mustSquare(6, rank)

	king, kingPresent := p.board.pieceAt(kingFrom)
	rook, rookPresent := p.board.pieceAt(rookFrom)
	if !kingPresent || !rookPresent || king.kind != King || king.side != side || rook.kind != Rook || rook.side != side {
		return false
	}

	if p.board.occupied(fSquare) || p.board.occupied(gSquare) {
		return false
	}

	opponent := side.Opponent()
	if p.isSquareAttacked(kingFrom, opponent) || p.isSquareAttacked(fSquare, opponent) || p.isSquareAttacked(gSquare, opponent) {
		return false
	}

	return true
}

func (p Position) canCastleQueenside(side Side) bool {
	if !p.castlingRights.CanCastleQueenside(side) {
		return false
	}

	rank := 0
	if side == Black {
		rank = 7
	}

	kingFrom := mustSquare(4, rank)
	rookFrom := mustSquare(0, rank)
	bSquare := mustSquare(1, rank)
	cSquare := mustSquare(2, rank)
	dSquare := mustSquare(3, rank)

	king, kingPresent := p.board.pieceAt(kingFrom)
	rook, rookPresent := p.board.pieceAt(rookFrom)
	if !kingPresent || !rookPresent || king.kind != King || king.side != side || rook.kind != Rook || rook.side != side {
		return false
	}

	if p.board.occupied(bSquare) || p.board.occupied(cSquare) || p.board.occupied(dSquare) {
		return false
	}

	opponent := side.Opponent()
	if p.isSquareAttacked(kingFrom, opponent) || p.isSquareAttacked(dSquare, opponent) || p.isSquareAttacked(cSquare, opponent) {
		return false
	}

	return true
}

func (p Position) isSquareAttacked(square Square, by Side) bool {
	pawnRankDelta := -1
	if by == Black {
		pawnRankDelta = 1
	}

	for _, pawnFileDelta := range []int{-1, 1} {
		source, ok := square.offset(pawnFileDelta, pawnRankDelta)
		if !ok {
			continue
		}

		piece, found := p.board.pieceAt(source)
		if found && piece.side == by && piece.kind == Pawn {
			return true
		}
	}

	knightOffsets := [][2]int{
		{1, 2}, {2, 1}, {2, -1}, {1, -2},
		{-1, -2}, {-2, -1}, {-2, 1}, {-1, 2},
	}
	for _, offset := range knightOffsets {
		source, ok := square.offset(offset[0], offset[1])
		if !ok {
			continue
		}

		piece, found := p.board.pieceAt(source)
		if found && piece.side == by && piece.kind == Knight {
			return true
		}
	}

	if p.attackedBySliding(square, by, [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}, Rook, Queen) {
		return true
	}

	if p.attackedBySliding(square, by, [][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}, Bishop, Queen) {
		return true
	}

	for rankDelta := -1; rankDelta <= 1; rankDelta++ {
		for fileDelta := -1; fileDelta <= 1; fileDelta++ {
			if fileDelta == 0 && rankDelta == 0 {
				continue
			}

			source, ok := square.offset(fileDelta, rankDelta)
			if !ok {
				continue
			}

			piece, found := p.board.pieceAt(source)
			if found && piece.side == by && piece.kind == King {
				return true
			}
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
