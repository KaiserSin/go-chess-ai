package ai

import (
	"sync"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

type ttBound uint8

const (
	ttExact ttBound = iota
	ttLower
	ttUpper
)

type ttEntry struct {
	depth    int
	score    int
	bound    ttBound
	bestMove chess.Move
}

type transpositionTable struct {
	mu      sync.RWMutex
	entries map[chess.RepetitionKey]ttEntry
}

func newTranspositionTable() *transpositionTable {
	return &transpositionTable{
		entries: make(map[chess.RepetitionKey]ttEntry),
	}
}

func (tt *transpositionTable) probe(position chess.Position) (ttEntry, bool) {
	if tt == nil {
		return ttEntry{}, false
	}

	key := chess.NewRepetitionKey(position)

	tt.mu.RLock()
	entry, ok := tt.entries[key]
	tt.mu.RUnlock()

	return entry, ok
}

func (tt *transpositionTable) store(position chess.Position, entry ttEntry) {
	if tt == nil {
		return
	}

	key := chess.NewRepetitionKey(position)

	tt.mu.Lock()
	existing, ok := tt.entries[key]
	if !ok || entry.depth >= existing.depth {
		tt.entries[key] = entry
	}
	tt.mu.Unlock()
}
