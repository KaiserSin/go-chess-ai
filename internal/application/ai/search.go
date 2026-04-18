package ai

import (
	"sync"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

const (
	FixedSearchDepth = 3
	searchInfinity   = mateScore + 1
)

const (
	movePriorityPromotion = iota
	movePriorityCapture
	movePriorityQuiet
)

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

func BestMove(position chess.Position) SearchResult {
	return bestMove(position, FixedSearchDepth)
}

func bestMove(position chess.Position, depth int) SearchResult {
	rootPerspective := position.SideToMove()
	if depth <= 0 || isTerminalPosition(position) {
		return noMoveSearchResult(position, rootPerspective)
	}

	bestResult := noMoveSearchResult(position, rootPerspective)
	for currentDepth := 1; currentDepth <= depth; currentDepth++ {
		result := searchAtDepth(position, currentDepth)
		if result.HasMove {
			bestResult = result
		}
	}

	return bestResult
}

func searchAtDepth(position chess.Position, depth int) SearchResult {
	rootPerspective := position.SideToMove()
	if depth <= 0 || isTerminalPosition(position) {
		return noMoveSearchResult(position, rootPerspective)
	}

	originalMoves := position.LegalMoves()
	if len(originalMoves) == 0 {
		return noMoveSearchResult(position, rootPerspective)
	}

	orderedMoves := orderMoves(position, originalMoves)
	return pickBestRootResult(collectRootResults(position, originalMoves, orderedMoves, depth-1, rootPerspective))
}

func noMoveSearchResult(position chess.Position, rootPerspective chess.Side) SearchResult {
	return SearchResult{
		Score:   Evaluate(position, rootPerspective),
		HasMove: false,
	}
}

func searchRootMove(position chess.Position, move chess.Move, depth int, rootPerspective chess.Side, alpha int) int {
	next, err := position.ApplyMove(move)
	if err != nil {
		panic(err)
	}

	return alphaBeta(next, depth, alpha, searchInfinity, rootPerspective)
}

func collectRootResults(position chess.Position, originalMoves []chess.Move, orderedMoves []chess.Move, depth int, rootPerspective chess.Side) []rootMoveResult {
	firstMove := orderedMoves[0]
	firstResult := rootMoveResult{
		move:          firstMove,
		score:         searchRootMove(position, firstMove, depth, rootPerspective, -searchInfinity),
		originalIndex: findMoveIndex(originalMoves, firstMove),
	}

	results := []rootMoveResult{firstResult}
	if len(orderedMoves) == 1 {
		return results
	}

	alpha := firstResult.score
	remainingMoves := orderedMoves[1:]
	remainingResults := make([]rootMoveResult, len(remainingMoves))
	var wait sync.WaitGroup

	for index, move := range remainingMoves {
		wait.Add(1)
		go func(index int, move chess.Move) {
			defer wait.Done()
			remainingResults[index] = rootMoveResult{
				move:          move,
				score:         searchRootMove(position, move, depth, rootPerspective, alpha),
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

func alphaBeta(position chess.Position, depth int, alpha, beta int, rootPerspective chess.Side) int {
	if depth <= 0 || isTerminalPosition(position) {
		if isTerminalPosition(position) {
			return evaluateStatic(position, rootPerspective)
		}

		return quiescence(position, alpha, beta, rootPerspective)
	}

	moves := orderMoves(position, position.LegalMoves())
	if len(moves) == 0 {
		return evaluateStatic(position, rootPerspective)
	}

	maximizing := position.SideToMove() == rootPerspective
	bestScore := 0

	for index, move := range moves {
		next, err := position.ApplyMove(move)
		if err != nil {
			panic(err)
		}

		score := alphaBeta(next, depth-1, alpha, beta, rootPerspective)
		if index == 0 || betterScore(score, bestScore, maximizing) {
			bestScore = score
		}

		if maximizing {
			if score > alpha {
				alpha = score
			}
		} else if score < beta {
			beta = score
		}

		if alpha >= beta {
			break
		}
	}

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

func quiescence(position chess.Position, alpha, beta int, rootPerspective chess.Side) int {
	standPat := evaluateStatic(position, rootPerspective)
	maximizing := position.SideToMove() == rootPerspective
	bestScore := standPat

	if maximizing {
		if standPat >= beta {
			return standPat
		}

		if standPat > alpha {
			alpha = standPat
		}
	} else {
		if standPat <= alpha {
			return standPat
		}

		if standPat < beta {
			beta = standPat
		}
	}

	moves := tacticalMoves(position)
	for _, move := range moves {
		next, err := position.ApplyMove(move)
		if err != nil {
			panic(err)
		}

		score := quiescence(next, alpha, beta, rootPerspective)
		if betterScore(score, bestScore, maximizing) {
			bestScore = score
		}

		if maximizing {
			if score > alpha {
				alpha = score
			}
		} else if score < beta {
			beta = score
		}

		if alpha >= beta {
			break
		}
	}

	return bestScore
}

func tacticalMoves(position chess.Position) []chess.Move {
	moves := position.LegalMoves()
	tactical := make([]chess.Move, 0, len(moves))

	for _, move := range moves {
		if move.Promotion != chess.NoPieceType || isCaptureMove(position, move) {
			tactical = append(tactical, move)
		}
	}

	return orderMoves(position, tactical)
}

func evaluateStatic(position chess.Position, rootPerspective chess.Side) int {
	return Evaluate(position, rootPerspective)
}

func orderMoves(position chess.Position, moves []chess.Move) []chess.Move {
	promotions := make([]chess.Move, 0, len(moves))
	captures := make([]chess.Move, 0, len(moves))
	quiets := make([]chess.Move, 0, len(moves))

	for _, move := range moves {
		switch movePriority(position, move) {
		case movePriorityPromotion:
			promotions = append(promotions, move)
		case movePriorityCapture:
			captures = append(captures, move)
		default:
			quiets = append(quiets, move)
		}
	}

	ordered := append(promotions, captures...)
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
