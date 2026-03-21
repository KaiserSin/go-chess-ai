package position

import chessmodel "github.com/KaiserSin/go-chess-ai/internal/domain/chess/model"

func IsFiftyMoveDraw(position Position) bool {
	return position.halfmoveClock >= 100
}

func (p Position) IsInCheck(side chessmodel.Side) bool {
	return p.isInCheck(side)
}

func (p Position) Status() chessmodel.Status {
	return p.status()
}

func (p Position) isInCheck(side chessmodel.Side) bool {
	kingSquare, ok := p.board.kingSquare(side)
	if !ok {
		return false
	}

	return p.isSquareAttacked(kingSquare, side.Opponent())
}

func (p Position) status() chessmodel.Status {
	inCheck := p.isInCheck(p.sideToMove)
	if len(p.generateLegalMoves()) == 0 {
		if inCheck {
			return chessmodel.Checkmate
		}
		return chessmodel.Stalemate
	}

	if inCheck {
		return chessmodel.Check
	}

	return chessmodel.Ongoing
}

func (p Position) isFinishedPosition() bool {
	status := p.status()
	return status == chessmodel.Checkmate || status == chessmodel.Stalemate
}

func (p Position) effectiveEnPassantSquare() optionalSquare {
	if !p.enPassant.ok {
		return noSquare()
	}

	for _, move := range p.generateLegalMoves() {
		if move.To != p.enPassant.value || move.From.File() == move.To.File() {
			continue
		}

		piece, ok := p.board.pieceAt(move.From)
		if ok && piece.Type() == chessmodel.Pawn {
			return p.enPassant
		}
	}

	return noSquare()
}
