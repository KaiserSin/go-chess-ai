package model

import "errors"

var (
	ErrInvalidSquare     = errors.New("bad square")
	ErrInvalidPosition   = errors.New("bad position")
	ErrNoPiece           = errors.New("no piece")
	ErrWrongSide         = errors.New("wrong side")
	ErrInvalidMove       = errors.New("bad move")
	ErrGameFinished      = errors.New("game is over")
	ErrPromotionRequired = errors.New("need promotion piece")
	ErrInvalidPromotion  = errors.New("bad promotion piece")
)
