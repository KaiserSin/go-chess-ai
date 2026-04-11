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
		result := bestMoveAtDepth(position, currentDepth, table)
		if result.HasMove {
			bestResult = result
		}
	}

	return bestResult
}

func bestMoveAtDepth(position chess.Position, depth int, table *transpositionTable) SearchResult {
	return bestMoveAtDepthWithHooks(position, depth, table, nil)
}

func bestMoveAtDepthWithHooks(position chess.Position, depth int, table *transpositionTable, hooks *searchHooks) SearchResult {
	rootPerspective := position.SideToMove()
	if depth <= 0 || isTerminalPosition(position) {
		return noMoveSearchResult(position, rootPerspective)
	}

	moves := position.LegalMoves()
	if len(moves) == 0 {
		return noMoveSearchResult(position, rootPerspective)
	}

	rootEntry, hasRootEntry := table.probe(position)
	ordered := reorderMoves(position, moves, rootEntry.bestMove, hasRootEntry)
	moveIndex := indexMoves(moves)
	results := make([]SearchResult, len(moves))

	firstMove := ordered[0]
	firstScore := searchRootMove(position, firstMove, depth-1, rootPerspective, -searchInfinity, table, hooks)
	results[moveIndex[firstMove]] = SearchResult{
		Move:    firstMove,
		Score:   firstScore,
		HasMove: true,
	}

	if len(ordered) > 1 {
		if hooks != nil {
			searchRemainingRootMovesSequential(position, ordered[1:], moveIndex, depth-1, rootPerspective, firstScore, table, hooks, results)
		} else {
			searchRemainingRootMovesParallel(position, ordered[1:], moveIndex, depth-1, rootPerspective, firstScore, table, results)
		}
	}

	bestResult := pickBestResult(results)
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

func indexMoves(moves []chess.Move) map[chess.Move]int {
	indexes := make(map[chess.Move]int, len(moves))
	for index, move := range moves {
		indexes[move] = index
	}

	return indexes
}

func searchRootMove(position chess.Position, move chess.Move, depth int, rootPerspective chess.Side, alpha int, table *transpositionTable, hooks *searchHooks) int {
	next, err := position.ApplyMove(move)
	if err != nil {
		panic(err)
	}

	return alphaBeta(next, depth, alpha, searchInfinity, rootPerspective, table, hooks)
}

func searchRemainingRootMovesSequential(position chess.Position, moves []chess.Move, moveIndex map[chess.Move]int, depth int, rootPerspective chess.Side, alpha int, table *transpositionTable, hooks *searchHooks, results []SearchResult) {
	for _, move := range moves {
		score := searchRootMove(position, move, depth, rootPerspective, alpha, table, hooks)
		results[moveIndex[move]] = SearchResult{
			Move:    move,
			Score:   score,
			HasMove: true,
		}
	}
}

func searchRemainingRootMovesParallel(position chess.Position, moves []chess.Move, moveIndex map[chess.Move]int, depth int, rootPerspective chess.Side, alpha int, table *transpositionTable, results []SearchResult) {
	var wait sync.WaitGroup
	var mu sync.Mutex

	for _, move := range moves {
		index := moveIndex[move]
		wait.Add(1)
		go func(index int, move chess.Move) {
			defer wait.Done()

			score := searchRootMove(position, move, depth, rootPerspective, alpha, table, nil)
			mu.Lock()
			results[index] = SearchResult{
				Move:    move,
				Score:   score,
				HasMove: true,
			}
			mu.Unlock()
		}(index, move)
	}

	wait.Wait()
}

func pickBestResult(results []SearchResult) SearchResult {
	best := results[0]
	for index := 1; index < len(results); index++ {
		if betterScore(results[index].Score, best.Score, true) {
			best = results[index]
		}
	}

	return best
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

	moves := orderedMoves(position, entry.bestMove, hasEntry)
	if len(moves) == 0 {
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

func orderedMoves(position chess.Position, hashMove chess.Move, hasHashMove bool) []chess.Move {
	return reorderMoves(position, position.LegalMoves(), hashMove, hasHashMove)
}

func reorderMoves(position chess.Position, moves []chess.Move, hashMove chess.Move, hasHashMove bool) []chess.Move {
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
