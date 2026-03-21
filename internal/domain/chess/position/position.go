package position

import (
	"github.com/KaiserSin/go-chess-ai/internal/domain/chess/internal/bitboard"
	chessmodel "github.com/KaiserSin/go-chess-ai/internal/domain/chess/model"
)

type Position struct {
	board          Board
	sideToMove     chessmodel.Side
	castlingRights chessmodel.CastlingRights
	enPassant      optionalSquare
	halfmoveClock  int
	fullmoveNumber int
}

func NewInitialPosition() Position {
	return Position{
		board:          newInitialBoard(),
		sideToMove:     chessmodel.White,
		castlingRights: initialCastlingRights(),
		enPassant:      noSquare(),
		halfmoveClock:  0,
		fullmoveNumber: 1,
	}
}

func Validate(position Position) error {
	return position.validate()
}

func (p Position) PieceAt(square chessmodel.Square) (chessmodel.Piece, bool) {
	return p.board.PieceAt(square)
}

func (p Position) SideToMove() chessmodel.Side {
	return p.sideToMove
}

func (p Position) CastlingRights() chessmodel.CastlingRights {
	return p.castlingRights
}

func (p Position) EnPassantSquare() (chessmodel.Square, bool) {
	return p.enPassant.value, p.enPassant.ok
}

func (p Position) HalfmoveClock() int {
	return p.halfmoveClock
}

func (p Position) FullmoveNumber() int {
	return p.fullmoveNumber
}

func (p Position) LegalMoves() []chessmodel.Move {
	return p.generateLegalMoves()
}

func (p Position) validate() error {
	if !isValidSide(p.sideToMove) {
		return chessmodel.ErrInvalidPosition
	}

	if bitboard.Count(p.board.bitboard(chessmodel.White, chessmodel.King)) != 1 {
		return chessmodel.ErrInvalidPosition
	}

	if bitboard.Count(p.board.bitboard(chessmodel.Black, chessmodel.King)) != 1 {
		return chessmodel.ErrInvalidPosition
	}

	return nil
}
