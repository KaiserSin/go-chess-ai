package position

import chessmodel "github.com/KaiserSin/go-chess-ai/internal/domain/chess/model"

func (p Position) canCastleKingside(side chessmodel.Side) bool {
	return p.canCastle(side, kingSide)
}

func (p Position) canCastleQueenside(side chessmodel.Side) bool {
	return p.canCastle(side, queenSide)
}

func (p Position) canCastle(side chessmodel.Side, castle castleSide) bool {
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

func (p Position) hasCastlingPieces(side chessmodel.Side, castle castleSide) bool {
	king, kingPresent := p.board.pieceAt(kingStartSquare(side))
	rook, rookPresent := p.board.pieceAt(rookStartSquare(side, castle))
	if !kingPresent || !rookPresent {
		return false
	}

	return king.Type() == chessmodel.King &&
		king.Side() == side &&
		rook.Type() == chessmodel.Rook &&
		rook.Side() == side
}
