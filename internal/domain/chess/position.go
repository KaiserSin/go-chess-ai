package chess

// Position represents a full chess position.
type Position struct {
	board          Board
	sideToMove     Side
	castlingRights CastlingRights
	enPassant      optionalSquare
	halfmoveClock  int
	fullmoveNumber int
}

// NewInitialPosition creates the standard chess starting position.
func NewInitialPosition() Position {
	return Position{
		board:          newInitialBoard(),
		sideToMove:     White,
		castlingRights: initialCastlingRights(),
		enPassant:      noSquare(),
		halfmoveClock:  0,
		fullmoveNumber: 1,
	}
}

// PieceAt returns the piece on the target square.
func (p Position) PieceAt(square Square) (Piece, bool) {
	return p.board.PieceAt(square)
}

// SideToMove returns the side to move.
func (p Position) SideToMove() Side {
	return p.sideToMove
}

// CastlingRights returns the remaining castling rights for the position.
func (p Position) CastlingRights() CastlingRights {
	return p.castlingRights
}

// EnPassantSquare returns the current en passant target square, if any.
func (p Position) EnPassantSquare() (Square, bool) {
	return p.enPassant.value, p.enPassant.ok
}

// HalfmoveClock returns the number of halfmoves since the last pawn move or capture.
func (p Position) HalfmoveClock() int {
	return p.halfmoveClock
}

// FullmoveNumber returns the current fullmove number.
func (p Position) FullmoveNumber() int {
	return p.fullmoveNumber
}

// LegalMoves returns all legal moves for the current side.
func (p Position) LegalMoves() []Move {
	return p.generateLegalMoves()
}

// IsLegalMove reports whether a move is legal in the current position.
func (p Position) IsLegalMove(move Move) bool {
	if err := move.validateSquares(); err != nil {
		return false
	}

	for _, candidate := range p.generateLegalMoves() {
		if candidate == move {
			return true
		}
	}

	return false
}

// ApplyMove applies a legal move and returns the next position.
func (p Position) ApplyMove(move Move) (Position, error) {
	if err := move.validateSquares(); err != nil {
		return Position{}, err
	}

	status := p.status()
	if status == Checkmate || status == Stalemate {
		return Position{}, ErrGameFinished
	}

	piece, ok := p.board.pieceAt(move.From)
	if !ok {
		return Position{}, ErrNoPiece
	}

	if piece.side != p.sideToMove {
		return Position{}, ErrWrongSide
	}

	if err := validatePromotion(piece, move); err != nil {
		return Position{}, err
	}

	for _, candidate := range p.generateLegalMoves() {
		if candidate != move {
			continue
		}

		return p.applyMoveUnchecked(move), nil
	}

	return Position{}, ErrInvalidMove
}

// IsInCheck reports whether the given side is in check.
func (p Position) IsInCheck(side Side) bool {
	return p.isInCheck(side)
}

// Status returns the current play status.
func (p Position) Status() Status {
	return p.status()
}

func (p Position) isInCheck(side Side) bool {
	kingSquare, ok := p.board.kingSquare(side)
	if !ok {
		return false
	}

	return p.isSquareAttacked(kingSquare, side.Opponent())
}

func (p Position) status() Status {
	moves := p.generateLegalMoves()
	inCheck := p.isInCheck(p.sideToMove)
	if len(moves) == 0 {
		if inCheck {
			return Checkmate
		}
		return Stalemate
	}

	if inCheck {
		return Check
	}

	return Ongoing
}

func validatePromotion(piece Piece, move Move) error {
	if move.Promotion != NoPieceType && !move.Promotion.isPromotionChoice() {
		return ErrInvalidPromotion
	}

	if piece.kind != Pawn {
		if move.Promotion != NoPieceType {
			return ErrInvalidPromotion
		}
		return nil
	}

	targetRank := move.To.Rank()
	if (piece.side == White && targetRank == 7) || (piece.side == Black && targetRank == 0) {
		if move.Promotion == NoPieceType {
			return ErrPromotionRequired
		}
		if !move.Promotion.isPromotionChoice() {
			return ErrInvalidPromotion
		}
		return nil
	}

	if move.Promotion != NoPieceType {
		return ErrInvalidPromotion
	}

	return nil
}

func (p Position) applyMoveUnchecked(move Move) Position {
	next := p
	next.enPassant = noSquare()

	piece, _ := next.board.pieceAt(move.From)
	capturedPiece, captured := next.board.pieceAt(move.To)

	isEnPassantCapture := piece.kind == Pawn &&
		p.enPassant.ok &&
		move.To == p.enPassant.value &&
		move.From.File() != move.To.File() &&
		!captured

	if isEnPassantCapture {
		captureRank := move.To.Rank()
		if piece.side == White {
			captureRank--
		} else {
			captureRank++
		}
		captureSquare := mustSquare(move.To.File(), captureRank)
		capturedPiece, captured = next.board.removePiece(captureSquare)
	}

	if piece.kind == Pawn || captured {
		next.halfmoveClock = 0
	} else {
		next.halfmoveClock++
	}

	next.updateCastlingRightsForMove(piece, move, capturedPiece, captured)

	next.board.removePiece(move.From)
	if !isEnPassantCapture {
		next.board.removePiece(move.To)
	}

	if piece.kind == King && absInt(move.To.File()-move.From.File()) == 2 {
		next.board.placePiece(piece, move.To)
		next.moveCastlingRook(piece.side, move)
	} else {
		finalPiece := piece
		if piece.kind == Pawn && move.Promotion != NoPieceType {
			finalPiece = newPiece(piece.side, move.Promotion)
		}
		next.board.placePiece(finalPiece, move.To)
	}

	if piece.kind == Pawn && absInt(move.To.Rank()-move.From.Rank()) == 2 {
		passedRank := (move.From.Rank() + move.To.Rank()) / 2
		next.enPassant = someSquare(mustSquare(move.From.File(), passedRank))
	}

	next.sideToMove = p.sideToMove.Opponent()
	if piece.side == Black {
		next.fullmoveNumber++
	}

	return next
}

func (p *Position) updateCastlingRightsForMove(piece Piece, move Move, captured Piece, capturedPresent bool) {
	switch piece.kind {
	case King:
		p.castlingRights.removeSide(piece.side)
	case Rook:
		switch move.From {
		case mustSquare(0, 0):
			p.castlingRights.removeQueenside(White)
		case mustSquare(7, 0):
			p.castlingRights.removeKingside(White)
		case mustSquare(0, 7):
			p.castlingRights.removeQueenside(Black)
		case mustSquare(7, 7):
			p.castlingRights.removeKingside(Black)
		}
	}

	if capturedPresent && captured.kind == Rook {
		switch move.To {
		case mustSquare(0, 0):
			p.castlingRights.removeQueenside(White)
		case mustSquare(7, 0):
			p.castlingRights.removeKingside(White)
		case mustSquare(0, 7):
			p.castlingRights.removeQueenside(Black)
		case mustSquare(7, 7):
			p.castlingRights.removeKingside(Black)
		}
	}
}

func (p *Position) moveCastlingRook(side Side, move Move) {
	var rookFrom Square
	var rookTo Square

	switch {
	case side == White && move.To == mustSquare(6, 0):
		rookFrom = mustSquare(7, 0)
		rookTo = mustSquare(5, 0)
	case side == White && move.To == mustSquare(2, 0):
		rookFrom = mustSquare(0, 0)
		rookTo = mustSquare(3, 0)
	case side == Black && move.To == mustSquare(6, 7):
		rookFrom = mustSquare(7, 7)
		rookTo = mustSquare(5, 7)
	case side == Black && move.To == mustSquare(2, 7):
		rookFrom = mustSquare(0, 7)
		rookTo = mustSquare(3, 7)
	default:
		return
	}

	rook, ok := p.board.removePiece(rookFrom)
	if !ok {
		return
	}

	p.board.placePiece(rook, rookTo)
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}

	return value
}
