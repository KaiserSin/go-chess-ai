package position

import chessmodel "github.com/KaiserSin/go-chess-ai/internal/domain/chess/model"

type PositionBuilder struct {
	board          Board
	sideToMove     chessmodel.Side
	castlingRights chessmodel.CastlingRights
	enPassant      optionalSquare
	halfmoveClock  int
	fullmoveNumber int
	err            error
}

func NewPositionBuilder() *PositionBuilder {
	return &PositionBuilder{
		board:          newEmptyBoard(),
		sideToMove:     chessmodel.White,
		castlingRights: 0,
		enPassant:      noSquare(),
		halfmoveClock:  0,
		fullmoveNumber: 1,
	}
}

func (b *PositionBuilder) WithSideToMove(side chessmodel.Side) *PositionBuilder {
	if b.err != nil {
		return b
	}

	if !isValidSide(side) {
		b.err = chessmodel.ErrInvalidPosition
		return b
	}

	b.sideToMove = side
	return b
}

func (b *PositionBuilder) WithCastlingRights(rights chessmodel.CastlingRights) *PositionBuilder {
	if b.err != nil {
		return b
	}

	b.castlingRights = rights
	return b
}

func (b *PositionBuilder) WithEnPassantSquare(square chessmodel.Square) *PositionBuilder {
	if b.err != nil {
		return b
	}

	if !squareIsValid(square) {
		b.err = chessmodel.ErrInvalidSquare
		return b
	}

	b.enPassant = someSquare(square)
	return b
}

func (b *PositionBuilder) WithHalfmoveClock(clock int) *PositionBuilder {
	if b.err != nil {
		return b
	}

	if clock < 0 {
		b.err = chessmodel.ErrInvalidPosition
		return b
	}

	b.halfmoveClock = clock
	return b
}

func (b *PositionBuilder) WithFullmoveNumber(number int) *PositionBuilder {
	if b.err != nil {
		return b
	}

	if number <= 0 {
		b.err = chessmodel.ErrInvalidPosition
		return b
	}

	b.fullmoveNumber = number
	return b
}

func (b *PositionBuilder) Place(square chessmodel.Square, side chessmodel.Side, pieceType chessmodel.PieceType) *PositionBuilder {
	if b.err != nil {
		return b
	}

	if !squareIsValid(square) {
		b.err = chessmodel.ErrInvalidSquare
		return b
	}

	if !isValidSide(side) || pieceType < chessmodel.Pawn || pieceType > chessmodel.King {
		b.err = chessmodel.ErrInvalidPosition
		return b
	}

	b.board.placePiece(chessmodel.NewPiece(side, pieceType), square)
	return b
}

func (b *PositionBuilder) Build() (Position, error) {
	if b.err != nil {
		return Position{}, b.err
	}

	position := Position{
		board:          b.board,
		sideToMove:     b.sideToMove,
		castlingRights: b.castlingRights,
		enPassant:      b.enPassant,
		halfmoveClock:  b.halfmoveClock,
		fullmoveNumber: b.fullmoveNumber,
	}

	if err := position.validate(); err != nil {
		return Position{}, err
	}

	return position, nil
}
