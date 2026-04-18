package gameplay

import (
	"strings"

	"github.com/KaiserSin/go-chess-ai/internal/application/ai"
	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"

	"github.com/KaiserSin/go-chess-ai/internal/application/dto"
)

const FixedAISearchDepth = ai.FixedSearchDepth

type Service struct {
	game             *chess.Game
	selectedSquare   optionalSquare
	pendingPromotion *pendingPromotion
}

type optionalSquare struct {
	value chess.Square
	ok    bool
}

type pendingPromotion struct {
	from    chess.Square
	to      chess.Square
	options []chess.PieceType
}

func NewService() *Service {
	service := &Service{}
	service.NewGame()
	return service
}

func NewGame() *Service {
	return NewService()
}

func (s *Service) NewGame() {
	s.game = chess.NewGame()
	s.clearSelection()
	s.pendingPromotion = nil
}

func (s *Service) SelectSquare(square chess.Square) {
	if s.pendingPromotion != nil || s.game.IsFinished() {
		return
	}

	position := s.game.Position()

	if !s.selectedSquare.ok {
		s.selectOwnPiece(position, square)
		return
	}

	if s.selectedSquare.value == square {
		s.clearSelection()
		return
	}

	if s.selectOwnPiece(position, square) {
		return
	}

	if err := s.TryMove(square); err != nil {
		return
	}
}

func (s *Service) SelectSquareAt(file, rank int) {
	square, err := chess.NewSquare(file, rank)
	if err != nil {
		return
	}

	s.SelectSquare(square)
}

func (s *Service) TryMove(to chess.Square) error {
	if s.game.IsFinished() {
		return chess.ErrGameFinished
	}

	if s.pendingPromotion != nil {
		return chess.ErrInvalidMove
	}

	if !s.selectedSquare.ok {
		return chess.ErrInvalidMove
	}

	matches := s.movesToTarget(to)
	switch len(matches) {
	case 0:
		s.clearSelection()
		return chess.ErrInvalidMove
	case 1:
		return s.applyMove(matches[0])
	default:
		s.pendingPromotion = &pendingPromotion{
			from:    matches[0].From,
			to:      matches[0].To,
			options: orderedPromotionOptions(matches),
		}

		return nil
	}
}

func (s *Service) TryMoveAt(file, rank int) error {
	square, err := chess.NewSquare(file, rank)
	if err != nil {
		return chess.ErrInvalidSquare
	}

	return s.TryMove(square)
}

func (s *Service) ChoosePromotion(pieceType chess.PieceType) error {
	promotion := s.pendingPromotion
	if promotion == nil {
		return chess.ErrInvalidMove
	}

	if !containsPromotionChoice(promotion.options, pieceType) {
		return chess.ErrInvalidPromotion
	}

	move := chess.Move{
		From:      promotion.from,
		To:        promotion.to,
		Promotion: pieceType,
	}

	return s.applyMove(move)
}

func (s *Service) ChoosePromotionByName(pieceType string) error {
	promotion, ok := parsePromotionName(pieceType)
	if !ok {
		return chess.ErrInvalidPromotion
	}

	return s.ChoosePromotion(promotion)
}

func (s *Service) ApplyAIMove() error {
	if s.game.IsFinished() {
		return chess.ErrGameFinished
	}

	if s.pendingPromotion != nil {
		return chess.ErrInvalidMove
	}

	result := ai.BestMove(s.game.Position())
	if !result.HasMove {
		return chess.ErrGameFinished
	}

	return s.applyMove(result.Move)
}

func (s *Service) Snapshot() dto.GameSnapshot {
	position := s.game.Position()
	outcome := s.game.Outcome()
	legalTargets := s.legalTargets()

	snapshot := dto.GameSnapshot{
		Squares:        s.snapshotSquares(position, legalTargets),
		SideToMove:     position.SideToMove().String(),
		Status:         s.game.Status().String(),
		OutcomeReason:  outcome.Reason().String(),
		HalfmoveClock:  position.HalfmoveClock(),
		FullmoveNumber: position.FullmoveNumber(),
	}

	if winner, ok := outcome.Winner(); ok {
		snapshot.Winner = winner.String()
		snapshot.HasWinner = true
	}

	snapshot.Promotion = s.snapshotPromotion(position.SideToMove())

	return snapshot
}

func newServiceWithGame(game *chess.Game) *Service {
	return &Service{
		game: game,
	}
}

func mustSquare(file, rank int) chess.Square {
	square, err := chess.NewSquare(file, rank)
	if err != nil {
		panic(err)
	}

	return square
}

func pieceKey(piece chess.Piece) string {
	return piece.Side().String() + "-" + piece.Type().String()
}

func pieceKeyForSideAndType(side chess.Side, pieceType chess.PieceType) string {
	return side.String() + "-" + pieceType.String()
}

func someSquare(square chess.Square) optionalSquare {
	return optionalSquare{
		value: square,
		ok:    true,
	}
}

func (s *Service) selectOwnPiece(position chess.Position, square chess.Square) bool {
	piece, ok := position.PieceAt(square)
	if !ok || piece.Side() != position.SideToMove() {
		return false
	}

	s.selectedSquare = someSquare(square)
	return true
}

func (s *Service) clearSelection() {
	s.selectedSquare = optionalSquare{}
}

func (s *Service) clearTransientState() {
	s.clearSelection()
	s.pendingPromotion = nil
}

func (s *Service) applyMove(move chess.Move) error {
	if err := s.game.ApplyMove(move); err != nil {
		s.clearSelection()
		return err
	}

	s.clearTransientState()
	return nil
}

func (s *Service) snapshotSquares(position chess.Position, legalTargets map[chess.Square]bool) []dto.SquareSnapshot {
	squares := make([]dto.SquareSnapshot, 0, 64)
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			squares = append(squares, s.snapshotSquare(position, legalTargets, file, rank))
		}
	}

	return squares
}

func (s *Service) snapshotSquare(position chess.Position, legalTargets map[chess.Square]bool, file, rank int) dto.SquareSnapshot {
	square := mustSquare(file, rank)
	snapshot := dto.SquareSnapshot{
		File:        file,
		Rank:        rank,
		Algebraic:   square.String(),
		Selected:    s.selectedSquare.ok && s.selectedSquare.value == square,
		LegalTarget: legalTargets[square],
	}

	if piece, ok := position.PieceAt(square); ok {
		snapshot.Occupied = true
		snapshot.PieceKey = pieceKey(piece)
	}

	return snapshot
}

func (s *Service) snapshotPromotion(side chess.Side) *dto.PromotionSnapshot {
	if s.pendingPromotion == nil {
		return nil
	}

	snapshot := &dto.PromotionSnapshot{
		Visible:      true,
		TargetSquare: s.pendingPromotion.to.String(),
		Options:      make([]dto.PromotionOptionSnapshot, 0, len(s.pendingPromotion.options)),
	}

	for _, option := range s.pendingPromotion.options {
		snapshot.Options = append(snapshot.Options, dto.PromotionOptionSnapshot{
			PieceType: option.String(),
			PieceKey:  pieceKeyForSideAndType(side, option),
		})
	}

	return snapshot
}

func (s *Service) selectedMoves() []chess.Move {
	if !s.selectedSquare.ok {
		return nil
	}

	moves := make([]chess.Move, 0, 8)
	for _, move := range s.game.LegalMoves() {
		if move.From == s.selectedSquare.value {
			moves = append(moves, move)
		}
	}

	return moves
}

func (s *Service) movesToTarget(to chess.Square) []chess.Move {
	matches := make([]chess.Move, 0, 4)
	for _, move := range s.selectedMoves() {
		if move.To == to {
			matches = append(matches, move)
		}
	}

	return matches
}

func (s *Service) legalTargets() map[chess.Square]bool {
	targets := make(map[chess.Square]bool)
	for _, move := range s.selectedMoves() {
		targets[move.To] = true
	}

	return targets
}

func orderedPromotionOptions(moves []chess.Move) []chess.PieceType {
	available := make(map[chess.PieceType]bool, len(moves))
	for _, move := range moves {
		available[move.Promotion] = true
	}

	ordered := make([]chess.PieceType, 0, len(available))
	for _, pieceType := range []chess.PieceType{chess.Queen, chess.Rook, chess.Bishop, chess.Knight} {
		if available[pieceType] {
			ordered = append(ordered, pieceType)
		}
	}

	return ordered
}

func containsPromotionChoice(options []chess.PieceType, pieceType chess.PieceType) bool {
	for _, option := range options {
		if option == pieceType {
			return true
		}
	}

	return false
}

func parsePromotionName(name string) (chess.PieceType, bool) {
	switch strings.ToLower(name) {
	case "queen":
		return chess.Queen, true
	case "rook":
		return chess.Rook, true
	case "bishop":
		return chess.Bishop, true
	case "knight":
		return chess.Knight, true
	default:
		return chess.NoPieceType, false
	}
}
