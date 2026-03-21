package chess

type PositionBuilder struct {
	board          Board
	sideToMove     Side
	castlingRights CastlingRights
	enPassant      optionalSquare
	halfmoveClock  int
	fullmoveNumber int
	err            error
}

func NewPositionBuilder() *PositionBuilder {
	return &PositionBuilder{
		board:          newEmptyBoard(),
		sideToMove:     White,
		castlingRights: 0,
		enPassant:      noSquare(),
		halfmoveClock:  0,
		fullmoveNumber: 1,
	}
}

func (b *PositionBuilder) WithSideToMove(side Side) *PositionBuilder {
	if b.err != nil {
		return b
	}

	if !side.isValid() {
		b.err = ErrInvalidPosition
		return b
	}

	b.sideToMove = side
	return b
}

func (b *PositionBuilder) WithCastlingRights(rights CastlingRights) *PositionBuilder {
	if b.err != nil {
		return b
	}

	b.castlingRights = rights
	return b
}

func (b *PositionBuilder) WithEnPassantSquare(square Square) *PositionBuilder {
	if b.err != nil {
		return b
	}

	if !square.isValid() {
		b.err = ErrInvalidSquare
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
		b.err = ErrInvalidPosition
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
		b.err = ErrInvalidPosition
		return b
	}

	b.fullmoveNumber = number
	return b
}

func (b *PositionBuilder) Place(square Square, side Side, pieceType PieceType) *PositionBuilder {
	if b.err != nil {
		return b
	}

	if !square.isValid() {
		b.err = ErrInvalidSquare
		return b
	}

	if !side.isValid() || pieceType < Pawn || pieceType > King {
		b.err = ErrInvalidPosition
		return b
	}

	b.board.placePiece(newPiece(side, pieceType), square)
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
