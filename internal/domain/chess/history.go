package chess

type positionKey struct {
	pieces         [2][6]uint64
	sideToMove     Side
	castlingRights CastlingRights
	enPassant      optionalSquare
}

func (p Position) repetitionKey() positionKey {
	return positionKey{
		pieces:         p.board.pieces,
		sideToMove:     p.sideToMove,
		castlingRights: p.castlingRights,
		enPassant:      p.effectiveEnPassantSquare(),
	}
}

type repetitionHistory struct {
	counts map[positionKey]int
}

func newRepetitionHistory() repetitionHistory {
	return repetitionHistory{
		counts: make(map[positionKey]int),
	}
}

func (h *repetitionHistory) record(position Position) int {
	key := position.repetitionKey()
	h.counts[key]++
	return h.counts[key]
}

func (h repetitionHistory) count(position Position) int {
	return h.counts[position.repetitionKey()]
}
