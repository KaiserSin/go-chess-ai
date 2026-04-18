package ai

import (
	"time"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

const (
	FixedSearchDepth    = 3
	searchInfinity      = mateScore + 1
	aspirationWindow    = 50
	fullWindowThreshold = searchInfinity
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

type searchOptions struct {
	deadline      time.Time
	useAspiration bool
}

type scoreResult struct {
	score     int
	completed bool
}

func BestMove(position chess.Position) SearchResult {
	return bestMoveWithOptions(position, searchOptions{
		deadline:      time.Now().Add(fixedSearchTimeBudget()),
		useAspiration: true,
	})
}

func bestMoveWithOptions(position chess.Position, options searchOptions) SearchResult {
	rootPerspective := position.SideToMove()
	if isTerminalPosition(position) {
		return noMoveSearchResult(position, rootPerspective)
	}

	originalMoves := position.LegalMoves()
	if len(originalMoves) == 0 {
		return noMoveSearchResult(position, rootPerspective)
	}

	orderedMoves := orderMoves(position, originalMoves)
	bestResult := fallbackSearchResult(position, orderedMoves[0], rootPerspective)
	previousScore := 0

	for currentDepth := 1; currentDepth <= FixedSearchDepth; currentDepth++ {
		result, completed := searchAtDepth(position, currentDepth, previousScore, options)
		if !completed {
			break
		}

		bestResult = result
		previousScore = result.Score
	}

	return bestResult
}

func fixedSearchTimeBudget() time.Duration {
	return time.Second * time.Duration(FixedSearchDepth)
}

func fallbackSearchResult(position chess.Position, move chess.Move, rootPerspective chess.Side) SearchResult {
	next, err := position.ApplyMove(move)
	if err != nil {
		panic(err)
	}

	return SearchResult{
		Move:    move,
		Score:   Evaluate(next, rootPerspective),
		HasMove: true,
	}
}

func searchAtDepth(position chess.Position, depth int, previousScore int, options searchOptions) (SearchResult, bool) {
	rootPerspective := position.SideToMove()
	if depth <= 0 || isTerminalPosition(position) {
		return noMoveSearchResult(position, rootPerspective), true
	}

	originalMoves := position.LegalMoves()
	if len(originalMoves) == 0 {
		return noMoveSearchResult(position, rootPerspective), true
	}

	if deadlineExceeded(options) {
		return SearchResult{}, false
	}

	if !options.useAspiration || depth == 1 {
		return searchAtDepthWindow(position, depth, -searchInfinity, searchInfinity, options)
	}

	window := aspirationWindow
	for {
		alpha := maxInt(-searchInfinity, previousScore-window)
		beta := minInt(searchInfinity, previousScore+window)

		result, completed := searchAtDepthWindow(position, depth, alpha, beta, options)
		if !completed {
			return SearchResult{}, false
		}

		if result.Score <= alpha {
			window *= 2
			if window >= fullWindowThreshold {
				return searchAtDepthWindow(position, depth, -searchInfinity, searchInfinity, options)
			}
			continue
		}

		if result.Score >= beta {
			window *= 2
			if window >= fullWindowThreshold {
				return searchAtDepthWindow(position, depth, -searchInfinity, searchInfinity, options)
			}
			continue
		}

		return result, true
	}
}

func searchAtDepthWindow(position chess.Position, depth int, alpha, beta int, options searchOptions) (SearchResult, bool) {
	rootPerspective := position.SideToMove()
	if depth <= 0 || isTerminalPosition(position) {
		return noMoveSearchResult(position, rootPerspective), true
	}

	originalMoves := position.LegalMoves()
	if len(originalMoves) == 0 {
		return noMoveSearchResult(position, rootPerspective), true
	}

	orderedMoves := orderMoves(position, originalMoves)
	results, completed := collectRootResults(position, originalMoves, orderedMoves, depth-1, rootPerspective, alpha, beta, options)
	if !completed {
		return SearchResult{}, false
	}

	return pickBestRootResult(results), true
}

func noMoveSearchResult(position chess.Position, rootPerspective chess.Side) SearchResult {
	return SearchResult{
		Score:   Evaluate(position, rootPerspective),
		HasMove: false,
	}
}

func searchRootMove(position chess.Position, move chess.Move, depth int, alpha, beta int, rootPerspective chess.Side, options searchOptions) scoreResult {
	if deadlineExceeded(options) {
		return scoreResult{completed: false}
	}

	next, err := position.ApplyMove(move)
	if err != nil {
		panic(err)
	}

	return alphaBeta(next, depth, alpha, beta, rootPerspective, options)
}

func collectRootResults(position chess.Position, originalMoves []chess.Move, orderedMoves []chess.Move, depth int, rootPerspective chess.Side, alpha, beta int, options searchOptions) ([]rootMoveResult, bool) {
	results := make([]rootMoveResult, 0, len(orderedMoves))

	for _, move := range orderedMoves {
		scored := searchRootMove(position, move, depth, alpha, beta, rootPerspective, options)
		if !scored.completed {
			return nil, false
		}

		result := rootMoveResult{
			move:          move,
			score:         scored.score,
			originalIndex: findMoveIndex(originalMoves, move),
		}
		results = append(results, result)

		if rootPerspective == position.SideToMove() {
			if scored.score > alpha {
				alpha = scored.score
			}
		} else if scored.score < beta {
			beta = scored.score
		}

		if alpha >= beta {
			break
		}
	}

	return results, true
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

func alphaBeta(position chess.Position, depth int, alpha, beta int, rootPerspective chess.Side, options searchOptions) scoreResult {
	if deadlineExceeded(options) {
		return scoreResult{completed: false}
	}

	if depth <= 0 || isTerminalPosition(position) {
		if isTerminalPosition(position) {
			return scoreResult{
				score:     evaluateStatic(position, rootPerspective),
				completed: true,
			}
		}

		return quiescence(position, alpha, beta, rootPerspective, options)
	}

	moves := orderMoves(position, position.LegalMoves())
	if len(moves) == 0 {
		return scoreResult{
			score:     evaluateStatic(position, rootPerspective),
			completed: true,
		}
	}

	maximizing := position.SideToMove() == rootPerspective
	bestScore := 0

	for index, move := range moves {
		if deadlineExceeded(options) {
			return scoreResult{completed: false}
		}

		next, err := position.ApplyMove(move)
		if err != nil {
			panic(err)
		}

		result := alphaBeta(next, depth-1, alpha, beta, rootPerspective, options)
		if !result.completed {
			return scoreResult{completed: false}
		}

		score := result.score
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

	return scoreResult{
		score:     bestScore,
		completed: true,
	}
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

func quiescence(position chess.Position, alpha, beta int, rootPerspective chess.Side, options searchOptions) scoreResult {
	if deadlineExceeded(options) {
		return scoreResult{completed: false}
	}

	standPat := evaluateStatic(position, rootPerspective)
	maximizing := position.SideToMove() == rootPerspective
	bestScore := standPat

	if maximizing {
		if standPat >= beta {
			return scoreResult{score: standPat, completed: true}
		}

		if standPat > alpha {
			alpha = standPat
		}
	} else {
		if standPat <= alpha {
			return scoreResult{score: standPat, completed: true}
		}

		if standPat < beta {
			beta = standPat
		}
	}

	moves := tacticalMoves(position)
	for _, move := range moves {
		if deadlineExceeded(options) {
			return scoreResult{completed: false}
		}

		next, err := position.ApplyMove(move)
		if err != nil {
			panic(err)
		}

		result := quiescence(next, alpha, beta, rootPerspective, options)
		if !result.completed {
			return scoreResult{completed: false}
		}

		score := result.score
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

	return scoreResult{
		score:     bestScore,
		completed: true,
	}
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

func deadlineExceeded(options searchOptions) bool {
	if options.deadline.IsZero() {
		return false
	}

	return !time.Now().Before(options.deadline)
}
