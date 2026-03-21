package position

import chessmodel "github.com/KaiserSin/go-chess-ai/internal/domain/chess/model"

type RepetitionKey struct {
	pieces         [2][6]uint64
	sideToMove     chessmodel.Side
	castlingRights chessmodel.CastlingRights
	enPassant      optionalSquare
}

func NewRepetitionKey(position Position) RepetitionKey {
	return RepetitionKey{
		pieces:         position.board.pieces,
		sideToMove:     position.sideToMove,
		castlingRights: position.castlingRights,
		enPassant:      position.effectiveEnPassantSquare(),
	}
}
