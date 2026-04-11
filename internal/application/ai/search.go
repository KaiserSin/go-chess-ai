package ai

import (
	"sync"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

const searchInfinity = mateScore + 1

const (
	movePriorityPromotion = iota
	movePriorityCapture
	movePriorityQuiet
)

type Node struct {
	Position chess.Position
	Move     chess.Move
	Score    int
	Children []Node
}

type SearchResult struct {
	Move    chess.Move
	Score   int
	HasMove bool
}

type rootMoveResult struct {
	move          chess.Move
	score         int
	originalIndex int
}

type searchHooks struct {
	leafEvaluations int
	cutoffs         int
	ttHits          int
}

func BuildTree(position chess.Position, depth int) Node {
	root := Node{
		Position: position,
	}

	if depth <= 0 || isTerminalPosition(position) {
		return root
	}

	moves := position.LegalMoves()
	if len(moves) == 0 {
		return root
	}

	root.Children = make([]Node, 0, len(moves))
	for _, move := range moves {
		next, err := position.ApplyMove(move)
		if err != nil {
			panic(err)
		}

		child := BuildTree(next, depth-1)
		child.Move = move
		root.Children = append(root.Children, child)
	}

	return root
}

func BestMove(position chess.Position, depth int) SearchResult {
	rootPerspective := position.SideToMove()
	if depth <= 0 || isTerminalPosition(position) {
		return noMoveSearchResult(position, rootPerspective)
	}

	table := newTranspositionTable()
	bestResult := noMoveSearchResult(position, rootPerspective)
	for currentDepth := 1; currentDepth <= depth; currentDepth++ {
		result := searchAtDepth(position, currentDepth, table, nil)
		if result.HasMove {
			bestResult = result
		}
	}

	return bestResult
}

func bestMoveAtDepth(position chess.Position, depth int, table *transpositionTable) SearchResult {
	return searchAtDepth(position, depth, table, nil)
}

func bestMoveAtDepthWithHooks(position chess.Position, depth int, table *transpositionTable, hooks *searchHooks) SearchResult {
	return searchAtDepth(position, depth, table, hooks)
}

func searchAtDepth(position chess.Position, depth int, table *transpositionTable, hooks *searchHooks) SearchResult {
	rootPerspective := position.SideToMove()
	if depth <= 0 || isTerminalPosition(position) {
		return noMoveSearchResult(position, rootPerspective)
	}

	originalMoves := position.LegalMoves()
	if len(originalMoves) == 0 {
		return noMoveSearchResult(position, rootPerspective)
	}

	rootEntry, hasRootEntry := table.probe(position)
	orderedMoves := orderMoves(position, originalMoves, rootEntry.bestMove, hasRootEntry)
	results := collectRootResults(position, originalMoves, orderedMoves, depth-1, rootPerspective, table, hooks)
	bestResult := pickBestRootResult(results)
	table.store(position, ttEntry{
		depth:    depth,
		score:    bestResult.Score,
		bound:    ttExact,
		bestMove: bestResult.Move,
	})

	return bestResult
}

func noMoveSearchResult(position chess.Position, rootPerspective chess.Side) SearchResult {
	return SearchResult{
		Score:   Evaluate(position, rootPerspective),
		HasMove: false,
	}
}

func searchRootMove(position chess.Position, move chess.Move, depth int, rootPerspective chess.Side, alpha int, table *transpositionTable, hooks *searchHooks) int {
	next, err := position.ApplyMove(move)
	if err != nil {
		panic(err)
	}

	return alphaBeta(next, depth, alpha, searchInfinity, rootPerspective, table, hooks)
}

func collectRootResults(position chess.Position, originalMoves []chess.Move, orderedMoves []chess.Move, depth int, rootPerspective chess.Side, table *transpositionTable, hooks *searchHooks) []rootMoveResult {
	firstMove := orderedMoves[0]
	firstResult := rootMoveResult{
		move:          firstMove,
		score:         searchRootMove(position, firstMove, depth, rootPerspective, -searchInfinity, table, hooks),
		originalIndex: findMoveIndex(originalMoves, firstMove),
	}

	results := []rootMoveResult{firstResult}
	if len(orderedMoves) == 1 {
		return results
	}

	alpha := firstResult.score
	remainingMoves := orderedMoves[1:]
	if hooks != nil {
		for _, move := range remainingMoves {
			results = append(results, rootMoveResult{
				move:          move,
				score:         searchRootMove(position, move, depth, rootPerspective, alpha, table, hooks),
				originalIndex: findMoveIndex(originalMoves, move),
			})
		}

		return results
	}

	remainingResults := make([]rootMoveResult, len(remainingMoves))
	var wait sync.WaitGroup

	for index, move := range remainingMoves {
		wait.Add(1)
		go func(index int, move chess.Move) {
			defer wait.Done()
			remainingResults[index] = rootMoveResult{
				move:          move,
				score:         searchRootMove(position, move, depth, rootPerspective, alpha, table, nil),
				originalIndex: findMoveIndex(originalMoves, move),
			}
		}(index, move)
	}

	wait.Wait()
	return append(results, remainingResults...)
}

func findMoveIndex(moves []chess.Move, target chess.Move) int {
	for index, move := range moves {
		if move == target {
			return index
		}
	}

	panic("move not found in original move list")
}

func pickBestRootResult(results []rootMoveResult) SearchResult {
	best := results[0]

	for _, result := range results[1:] {
		if result.score > best.score || (result.score == best.score && result.originalIndex < best.originalIndex) {
			best = result
		}
	}

	return SearchResult{
		Move:    best.move,
		Score:   best.score,
		HasMove: true,
	}
}

func alphaBeta(position chess.Position, depth int, alpha, beta int, rootPerspective chess.Side, table *transpositionTable, hooks *searchHooks) int {
	alphaOrig := alpha
	betaOrig := beta

	entry, hasEntry := table.probe(position)
	if hasEntry {
		if hooks != nil {
			hooks.ttHits++
		}

		if entry.depth >= depth {
			switch entry.bound {
			case ttExact:
				return entry.score
			case ttLower:
				if entry.score > alpha {
					alpha = entry.score
				}
			case ttUpper:
				if entry.score < beta {
					beta = entry.score
				}
			}

			if alpha >= beta {
				return entry.score
			}
		}
	}

	if depth <= 0 || isTerminalPosition(position) {
		return evaluateAndStore(position, depth, rootPerspective, table, hooks)
	}

	moves := orderMoves(position, position.LegalMoves(), entry.bestMove, hasEntry)
	if len(moves) == 0 {
		return evaluateAndStore(position, depth, rootPerspective, table, hooks)
	}

	maximizing := position.SideToMove() == rootPerspective
	bestScore := 0
	bestMove := chess.Move{}

	for index, move := range moves {
		next, err := position.ApplyMove(move)
		if err != nil {
			panic(err)
		}

		score := alphaBeta(next, depth-1, alpha, beta, rootPerspective, table, hooks)
		if index == 0 || betterScore(score, bestScore, maximizing) {
			bestScore = score
			bestMove = move
		}

		if maximizing {
			if score > alpha {
				alpha = score
			}
		} else if score < beta {
			beta = score
		}

		if alpha >= beta {
			if hooks != nil {
				hooks.cutoffs++
			}

			break
		}
	}

	entryBound := ttExact
	if bestScore <= alphaOrig {
		entryBound = ttUpper
	} else if bestScore >= betaOrig {
		entryBound = ttLower
	}

	table.store(position, ttEntry{
		depth:    depth,
		score:    bestScore,
		bound:    entryBound,
		bestMove: bestMove,
	})

	return bestScore
}

func betterScore(candidate, current int, maximizing bool) bool {
	if maximizing {
		return candidate > current
	}

	return candidate < current
}

func isTerminalPosition(position chess.Position) bool {
	status := position.Status()
	if status == chess.Checkmate || status == chess.Stalemate {
		return true
	}

	return chess.HasInsufficientMaterial(position) || chess.IsFiftyMoveDraw(position)
}

func evaluateAndStore(position chess.Position, depth int, rootPerspective chess.Side, table *transpositionTable, hooks *searchHooks) int {
	if hooks != nil {
		hooks.leafEvaluations++
	}

	score := Evaluate(position, rootPerspective)
	table.store(position, ttEntry{
		depth: depth,
		score: score,
		bound: ttExact,
	})

	return score
}

func orderMoves(position chess.Position, moves []chess.Move, hashMove chess.Move, hasHashMove bool) []chess.Move {
	hashes := make([]chess.Move, 0, 1)
	promotions := make([]chess.Move, 0, len(moves))
	captures := make([]chess.Move, 0, len(moves))
	quiets := make([]chess.Move, 0, len(moves))

	for _, move := range moves {
		if hasHashMove && move == hashMove {
			hashes = append(hashes, move)
			continue
		}

		switch movePriority(position, move) {
		case movePriorityPromotion:
			promotions = append(promotions, move)
		case movePriorityCapture:
			captures = append(captures, move)
		default:
			quiets = append(quiets, move)
		}
	}

	ordered := append(hashes, promotions...)
	ordered = append(ordered, captures...)
	ordered = append(ordered, quiets...)
	return ordered
}

func movePriority(position chess.Position, move chess.Move) int {
	if move.Promotion != chess.NoPieceType {
		return movePriorityPromotion
	}

	if isCaptureMove(position, move) {
		return movePriorityCapture
	}

	return movePriorityQuiet
}

func isCaptureMove(position chess.Position, move chess.Move) bool {
	if _, ok := position.PieceAt(move.To); ok {
		return true
	}

	piece, ok := position.PieceAt(move.From)
	if !ok || piece.Type() != chess.Pawn {
		return false
	}

	return move.From.File() != move.To.File()
}
