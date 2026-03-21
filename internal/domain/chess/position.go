package chess

import "math/bits"

type Position struct {
	board          Board
	sideToMove     Side
	castlingRights CastlingRights
	enPassant      optionalSquare
	halfmoveClock  int
	fullmoveNumber int
}

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

func (p Position) PieceAt(square Square) (Piece, bool) {
	return p.board.PieceAt(square)
}

func (p Position) SideToMove() Side {
	return p.sideToMove
}

func (p Position) CastlingRights() CastlingRights {
	return p.castlingRights
}

func (p Position) EnPassantSquare() (Square, bool) {
	return p.enPassant.value, p.enPassant.ok
}

func (p Position) HalfmoveClock() int {
	return p.halfmoveClock
}

func (p Position) FullmoveNumber() int {
	return p.fullmoveNumber
}

func (p Position) LegalMoves() []Move {
	return p.generateLegalMoves()
}

func (p Position) IsLegalMove(move Move) bool {
	if err := move.validateSquares(); err != nil {
		return false
	}

	return p.hasLegalMove(move)
}

func (p Position) ApplyMove(move Move) (Position, error) {
	if _, err := p.validateMoveRequest(move); err != nil {
		return Position{}, err
	}

	if !p.hasLegalMove(move) {
		return Position{}, ErrInvalidMove
	}

	return p.applyMoveUnchecked(move), nil
}

func (p Position) IsInCheck(side Side) bool {
	return p.isInCheck(side)
}

func (p Position) Status() Status {
	return p.status()
}

func (p Position) validate() error {
	if !p.sideToMove.isValid() {
		return ErrInvalidPosition
	}

	if bits.OnesCount64(p.board.bitboard(White, King)) != 1 {
		return ErrInvalidPosition
	}

	if bits.OnesCount64(p.board.bitboard(Black, King)) != 1 {
		return ErrInvalidPosition
	}

	return nil
}

func (p Position) isInCheck(side Side) bool {
	kingSquare, ok := p.board.kingSquare(side)
	if !ok {
		return false
	}

	return p.isSquareAttacked(kingSquare, side.Opponent())
}

func (p Position) status() Status {
	inCheck := p.isInCheck(p.sideToMove)
	if len(p.generateLegalMoves()) == 0 {
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

func (p Position) validateMoveRequest(move Move) (Piece, error) {
	if err := move.validateSquares(); err != nil {
		return Piece{}, err
	}

	if p.isFinishedPosition() {
		return Piece{}, ErrGameFinished
	}

	piece, ok := p.board.pieceAt(move.From)
	if !ok {
		return Piece{}, ErrNoPiece
	}

	if piece.side != p.sideToMove {
		return Piece{}, ErrWrongSide
	}

	if err := validatePromotion(piece, move); err != nil {
		return Piece{}, err
	}

	return piece, nil
}

func (p Position) isFinishedPosition() bool {
	status := p.status()
	return status == Checkmate || status == Stalemate
}

func (p Position) hasLegalMove(move Move) bool {
	for _, candidate := range p.generateLegalMoves() {
		if candidate == move {
			return true
		}
	}

	return false
}

func (p Position) isFiftyMoveDraw() bool {
	return p.halfmoveClock >= 100
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
		if ok && piece.kind == Pawn {
			return p.enPassant
		}
	}

	return noSquare()
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
	enPassantCaptureSquare, enPassantCapture := p.enPassantCaptureSquare(piece, move, captured)

	if enPassantCapture {
		capturedPiece, captured = next.board.removePiece(enPassantCaptureSquare)
	}

	if piece.kind == Pawn || captured {
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
		next.moveCastlingRook(piece.side, move)
	} else {
		next.board.placePiece(promotedPiece(piece, move), move.To)
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
		if side, castle, ok := rookCastlingRef(move.From); ok {
			p.castlingRights.remove(side, castle)
		}
	}

	if capturedPresent && captured.kind == Rook {
		if side, castle, ok := rookCastlingRef(move.To); ok {
			p.castlingRights.remove(side, castle)
		}
	}
}

func (p *Position) moveCastlingRook(side Side, move Move) {
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

func (p Position) enPassantCaptureSquare(piece Piece, move Move, captured bool) (Square, bool) {
	if piece.kind != Pawn || !p.enPassant.ok || captured {
		return 0, false
	}

	if move.To != p.enPassant.value || move.From.File() == move.To.File() {
		return 0, false
	}

	return mustSquare(move.To.File(), move.To.Rank()-pawnDirection(piece.side)), true
}

func isCastlingMove(piece Piece, move Move) bool {
	return piece.kind == King && absInt(move.To.File()-move.From.File()) == 2
}

func promotedPiece(piece Piece, move Move) Piece {
	if piece.kind == Pawn && move.Promotion != NoPieceType {
		return newPiece(piece.side, move.Promotion)
	}

	return piece
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}

	return value
}
