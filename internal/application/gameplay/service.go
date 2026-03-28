package gameplay

import (
	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"

	"github.com/KaiserSin/go-chess-ai/internal/application/dto"
)

type Service struct {
	game *chess.Game
}

func NewService() *Service {
	return &Service{
		game: chess.NewGame(),
	}
}

func (s *Service) Snapshot() dto.GameSnapshot {
	position := s.game.Position()
	outcome := s.game.Outcome()

	snapshot := dto.GameSnapshot{
		Squares:        make([]dto.SquareSnapshot, 0, 64),
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

	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			square := mustSquare(file, rank)
			squareSnapshot := dto.SquareSnapshot{
				File:      file,
				Rank:      rank,
				Algebraic: square.String(),
			}

			if piece, ok := position.PieceAt(square); ok {
				squareSnapshot.Occupied = true
				squareSnapshot.PieceKey = pieceKey(piece)
			}

			snapshot.Squares = append(snapshot.Squares, squareSnapshot)
		}
	}

	return snapshot
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
