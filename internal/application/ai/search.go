package ai

import (
	"time"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

const (
	DefaultSearchTimeBudget = 3 * time.Second
	searchInfinity          = mateScore + 1
	aspirationWindow        = 50
	fullWindowThreshold     = searchInfinity
)

const (
	movePriorityPromotion = iota
	movePriorityCapture
	movePriorityQuiet
)

// SearchResult contains the selected move and score returned by the AI search.
type SearchResult struct {
	Move    chess.Move
	Score   int
	HasMove bool
}

// rootMoveResult keeps a searched root move together with its original order.
type rootMoveResult struct {
	move          chess.Move
	score         int
	originalIndex int
}

// searchOptions controls deadlines and optional search optimizations.
type searchOptions struct {
	deadline      time.Time
	useAspiration bool
}

// scoreResult carries a searched score and whether the search finished fully.
type scoreResult struct {
	score     int
	completed bool
}

// BestMove searches the current position and returns the move chosen for the side to move.
func BestMove(position chess.Position) SearchResult {
	return BestMoveWithin(position, DefaultSearchTimeBudget)
}

// BestMoveWithin searches the current position until the time budget expires.
func BestMoveWithin(position chess.Position, budget time.Duration) SearchResult {
	return bestMoveWithOptions(position, searchOptions{
		deadline:      time.Now().Add(budget),
		useAspiration: true,
	})
}

// bestMoveWithOptions runs iterative deepening and returns the last completed result.
func bestMoveWithOptions(position chess.Position, options searchOptions) SearchResult {
	if options.deadline.IsZero() {
		options.deadline = time.Now().Add(DefaultSearchTimeBudget)
	}

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

	for currentDepth := 1; ; currentDepth++ {
		result, completed := searchAtDepth(position, currentDepth, previousScore, options)
		if !completed {
			break
		}

		bestResult = result
		previousScore = result.Score
	}

	return bestResult
}

// fallbackSearchResult gives the AI a legal move if a deeper search cannot finish.
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

// searchAtDepth searches one depth, first with aspiration windows when enabled.
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

// searchAtDepthWindow searches root moves inside one alpha beta window.
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

// noMoveSearchResult returns a static score for terminal or move-less positions.
func noMoveSearchResult(position chess.Position, rootPerspective chess.Side) SearchResult {
	return SearchResult{
		Score:   Evaluate(position, rootPerspective),
		HasMove: false,
	}
}

// searchRootMove applies one root move and searches the resulting position.
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

// collectRootResults evaluates ordered root moves and tracks their original order for tie breaking.
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

// findMoveIndex returns the move position in the original legal move list.
func findMoveIndex(moves []chess.Move, target chess.Move) int {
	for index, move := range moves {
		if move == target {
			return index
		}
	}

	panic("move not found in original move list")
}

// pickBestRootResult selects the highest scoring root move with deterministic tie breaking.
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

// alphaBeta searches legal moves recursively and cuts off branches when the alpha beta window closes.
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

// betterScore compares scores according to the current maximizing or minimizing side.
func betterScore(candidate, current int, maximizing bool) bool {
	if maximizing {
		return candidate > current
	}

	return candidate < current
}

// isTerminalPosition checks game-ending positions that should not be searched deeper.
func isTerminalPosition(position chess.Position) bool {
	status := position.Status()
	if status == chess.Checkmate || status == chess.Stalemate {
		return true
	}

	return chess.HasInsufficientMaterial(position) || chess.IsFiftyMoveDraw(position)
}

// quiescence extends static evaluation through tactical moves so the search does not stop on unstable captures or promotions.
func quiescence(position chess.Position, alpha, beta int, rootPerspective chess.Side, options searchOptions) scoreResult {
	if deadlineExceeded(options) {
		return scoreResult{completed: false}
	}

	if position.IsInCheck(position.SideToMove()) {
		return quiescenceEvasions(position, alpha, beta, rootPerspective, options)
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

// quiescenceEvasions searches all legal check evasions because stand-pat is invalid while in check.
func quiescenceEvasions(position chess.Position, alpha, beta int, rootPerspective chess.Side, options searchOptions) scoreResult {
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

		result := quiescence(next, alpha, beta, rootPerspective, options)
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

// tacticalMoves returns captures and promotions for quiescence search.
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

// evaluateStatic wraps the evaluation function used at search leaves.
func evaluateStatic(position chess.Position, rootPerspective chess.Side) int {
	return Evaluate(position, rootPerspective)
}

// orderMoves groups moves so tactically important moves are searched first.
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

// movePriority assigns a simple ordering bucket to one move.
func movePriority(position chess.Position, move chess.Move) int {
	if move.Promotion != chess.NoPieceType {
		return movePriorityPromotion
	}

	if isCaptureMove(position, move) {
		return movePriorityCapture
	}

	return movePriorityQuiet
}

// isCaptureMove detects normal captures and en passant-style pawn captures.
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

// deadlineExceeded reports whether the configured search deadline has passed.
func deadlineExceeded(options searchOptions) bool {
	if options.deadline.IsZero() {
		return false
	}

	return !time.Now().Before(options.deadline)
}
