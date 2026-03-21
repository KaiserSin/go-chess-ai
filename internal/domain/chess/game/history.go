package game

import "github.com/KaiserSin/go-chess-ai/internal/domain/chess/position"

type repetitionHistory struct {
	counts map[position.RepetitionKey]int
}

func newRepetitionHistory() repetitionHistory {
	return repetitionHistory{
		counts: make(map[position.RepetitionKey]int),
	}
}

func (h *repetitionHistory) record(pos position.Position) int {
	key := position.NewRepetitionKey(pos)
	h.counts[key]++
	return h.counts[key]
}

func (h repetitionHistory) count(pos position.Position) int {
	return h.counts[position.NewRepetitionKey(pos)]
}
