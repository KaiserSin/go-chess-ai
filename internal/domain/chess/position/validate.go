package position

import chessmodel "github.com/KaiserSin/go-chess-ai/internal/domain/chess/model"

func (p Position) IsLegalMove(move chessmodel.Move) bool {
	if err := validateMoveSquares(move); err != nil {
		return false
	}

	return p.hasLegalMove(move)
}

func (p Position) ApplyMove(move chessmodel.Move) (Position, error) {
	if _, err := p.validateMoveRequest(move); err != nil {
		return Position{}, err
	}

	if !p.hasLegalMove(move) {
		return Position{}, chessmodel.ErrInvalidMove
	}

	return p.applyMoveUnchecked(move), nil
}

func (p Position) validateMoveRequest(move chessmodel.Move) (chessmodel.Piece, error) {
	if err := validateMoveSquares(move); err != nil {
		return chessmodel.Piece{}, err
	}

	if p.isFinishedPosition() {
		return chessmodel.Piece{}, chessmodel.ErrGameFinished
	}

	piece, ok := p.board.pieceAt(move.From)
	if !ok {
		return chessmodel.Piece{}, chessmodel.ErrNoPiece
	}

	if piece.Side() != p.sideToMove {
		return chessmodel.Piece{}, chessmodel.ErrWrongSide
	}

	if err := validatePromotion(piece, move); err != nil {
		return chessmodel.Piece{}, err
	}

	return piece, nil
}

func (p Position) hasLegalMove(move chessmodel.Move) bool {
	for _, candidate := range p.generateLegalMoves() {
		if candidate == move {
			return true
		}
	}

	return false
}

func validatePromotion(piece chessmodel.Piece, move chessmodel.Move) error {
	if move.Promotion != chessmodel.NoPieceType && !isPromotionChoice(move.Promotion) {
		return chessmodel.ErrInvalidPromotion
	}

	if piece.Type() != chessmodel.Pawn {
		if move.Promotion != chessmodel.NoPieceType {
			return chessmodel.ErrInvalidPromotion
		}
		return nil
	}

	targetRank := move.To.Rank()
	if (piece.Side() == chessmodel.White && targetRank == 7) ||
		(piece.Side() == chessmodel.Black && targetRank == 0) {
		if move.Promotion == chessmodel.NoPieceType {
			return chessmodel.ErrPromotionRequired
		}
		return nil
	}

	if move.Promotion != chessmodel.NoPieceType {
		return chessmodel.ErrInvalidPromotion
	}

	return nil
}
