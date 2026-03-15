package chess

// PositionBuilder creates ad-hoc valid positions without exposing internal board details.
type PositionBuilder struct {
	board          Board
	sideToMove     Side
	castlingRights CastlingRights
	enPassant      optionalSquare
	halfmoveClock  int
	fullmoveNumber int
	err            error
}

// NewPositionBuilder creates a builder for custom positions.
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

// WithSideToMove sets the side to move.
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

// WithCastlingRights sets the castling rights.
func (b *PositionBuilder) WithCastlingRights(rights CastlingRights) *PositionBuilder {
	if b.err != nil {
		return b
	}

	b.castlingRights = rights
	return b
}

// WithEnPassantSquare sets the en passant target square.
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

// WithHalfmoveClock sets the halfmove clock.
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

// WithFullmoveNumber sets the fullmove number.
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

// Place adds a piece to the target square.
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

// Build validates and returns the assembled position.
func (b *PositionBuilder) Build() (Position, error) {
	if b.err != nil {
		return Position{}, b.err
	}

	if !b.sideToMove.isValid() {
		return Position{}, ErrInvalidPosition
	}

	if b.board.bitboard(White, King) == 0 || b.board.bitboard(Black, King) == 0 {
		return Position{}, ErrInvalidPosition
	}

	if b.board.bitboard(White, King)&(b.board.bitboard(White, King)-1) != 0 {
		return Position{}, ErrInvalidPosition
	}

	if b.board.bitboard(Black, King)&(b.board.bitboard(Black, King)-1) != 0 {
		return Position{}, ErrInvalidPosition
	}

	return Position{
		board:          b.board,
		sideToMove:     b.sideToMove,
		castlingRights: b.castlingRights,
		enPassant:      b.enPassant,
		halfmoveClock:  b.halfmoveClock,
		fullmoveNumber: b.fullmoveNumber,
	}, nil
}
