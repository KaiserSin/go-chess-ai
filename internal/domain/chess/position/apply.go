package position

import chessmodel "github.com/KaiserSin/go-chess-ai/internal/domain/chess/model"

func (p Position) applyMoveUnchecked(move chessmodel.Move) Position {
	next := p
	next.enPassant = noSquare()

	piece, _ := next.board.pieceAt(move.From)
	capturedPiece, captured := next.board.pieceAt(move.To)
	enPassantCaptureSquare, enPassantCapture := p.enPassantCaptureSquare(piece, move, captured)

	if enPassantCapture {
		capturedPiece, captured = next.board.removePiece(enPassantCaptureSquare)
	}

	if piece.Type() == chessmodel.Pawn || captured {
		next.halfmoveClock = 0
	} else {
		next.halfmoveClock++
	}

	next.updateCastlingRightsForMove(piece, move, capturedPiece, captured)

	next.board.removePiece(move.From)
	if !enPassantCapture {
		next.board.removePiece(move.To)
	}

	if isCastlingMove(piece, move) {
		next.board.placePiece(piece, move.To)
		next.moveCastlingRook(piece.Side(), move)
	} else {
		next.board.placePiece(promotedPiece(piece, move), move.To)
	}

	if piece.Type() == chessmodel.Pawn && absInt(move.To.Rank()-move.From.Rank()) == 2 {
		passedRank := (move.From.Rank() + move.To.Rank()) / 2
		next.enPassant = someSquare(squareMust(move.From.File(), passedRank))
	}

	next.sideToMove = p.sideToMove.Opponent()
	if piece.Side() == chessmodel.Black {
		next.fullmoveNumber++
	}

	return next
}

func (p *Position) updateCastlingRightsForMove(piece chessmodel.Piece, move chessmodel.Move, captured chessmodel.Piece, capturedPresent bool) {
	switch piece.Type() {
	case chessmodel.King:
		p.castlingRights = p.castlingRights.WithoutSide(piece.Side())
	case chessmodel.Rook:
		if side, castle, ok := rookCastlingRef(move.From); ok {
			if castle == kingSide {
				p.castlingRights = p.castlingRights.WithoutKingside(side)
			} else {
				p.castlingRights = p.castlingRights.WithoutQueenside(side)
			}
		}
	}

	if capturedPresent && captured.Type() == chessmodel.Rook {
		if side, castle, ok := rookCastlingRef(move.To); ok {
			if castle == kingSide {
				p.castlingRights = p.castlingRights.WithoutKingside(side)
			} else {
				p.castlingRights = p.castlingRights.WithoutQueenside(side)
			}
		}
	}
}

func (p *Position) moveCastlingRook(side chessmodel.Side, move chessmodel.Move) {
	castle, ok := castleSideFromKingTarget(side, move.To)
	if !ok {
		return
	}

	rook, ok := p.board.removePiece(rookStartSquare(side, castle))
	if !ok {
		return
	}

	p.board.placePiece(rook, castleRookSquare(side, castle))
}

func (p Position) enPassantCaptureSquare(piece chessmodel.Piece, move chessmodel.Move, captured bool) (chessmodel.Square, bool) {
	if piece.Type() != chessmodel.Pawn || !p.enPassant.ok || captured {
		return 0, false
	}

	if move.To != p.enPassant.value || move.From.File() == move.To.File() {
		return 0, false
	}

	return squareMust(move.To.File(), move.To.Rank()-pawnDirection(piece.Side())), true
}

func isCastlingMove(piece chessmodel.Piece, move chessmodel.Move) bool {
	return piece.Type() == chessmodel.King && absInt(move.To.File()-move.From.File()) == 2
}

func promotedPiece(piece chessmodel.Piece, move chessmodel.Move) chessmodel.Piece {
	if piece.Type() == chessmodel.Pawn && move.Promotion != chessmodel.NoPieceType {
		return chessmodel.NewPiece(piece.Side(), move.Promotion)
	}

	return piece
}
