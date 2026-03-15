package chess

import "errors"

var (
	// ErrInvalidSquare reports that a square is outside the board.
	ErrInvalidSquare = errors.New("invalid square")
	// ErrInvalidPosition reports that a constructed position violates domain invariants.
	ErrInvalidPosition = errors.New("invalid position")
	// ErrNoPiece reports that the source square is empty.
	ErrNoPiece = errors.New("no piece on source square")
	// ErrWrongSide reports that the selected piece does not belong to the side to move.
	ErrWrongSide = errors.New("piece does not belong to side to move")
	// ErrInvalidMove reports that the move is not legal in the current position.
	ErrInvalidMove = errors.New("invalid move")
	// ErrGameFinished reports that no more moves can be played in the current position.
	ErrGameFinished = errors.New("game already finished")
	// ErrPromotionRequired reports that a pawn move to the back rank needs a promotion piece.
	ErrPromotionRequired = errors.New("promotion piece is required")
	// ErrInvalidPromotion reports that the requested promotion piece is not allowed.
	ErrInvalidPromotion = errors.New("invalid promotion piece")
)
