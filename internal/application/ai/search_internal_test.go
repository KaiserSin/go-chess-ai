package ai

import (
	"testing"

	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

func TestBestMoveMatchesSequentialReference(t *testing.T) {
	tests := []struct {
		name     string
		position chess.Position
		depth    int
	}{
		{
			name:     "initial position",
			position: chess.NewInitialPosition(),
			depth:    2,
		},
		{
			name: "winning capture",
			position: mustBuildPosition(t,
				chess.NewPositionBuilder().
					WithSideToMove(chess.White).
					Place(mustParseSquare(t, "e1"), chess.White, chess.King).
					Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
					Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
					Place(mustParseSquare(t, "a8"), chess.Black, chess.Queen),
			),
			depth: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := bestMoveSequential(tt.position, tt.depth)
			got := BestMove(tt.position, tt.depth)

			if got.HasMove != want.HasMove {
				t.Fatalf("want HasMove %t, got %t", want.HasMove, got.HasMove)
			}

			if got.Score != want.Score {
				t.Fatalf("want score %d, got %d", want.Score, got.Score)
			}

			if got.Move != want.Move {
				t.Fatalf("want move %s, got %s", want.Move, got.Move)
			}
		})
	}
}

func TestBestMoveKeepsInitialTieBreak(t *testing.T) {
	result := BestMove(chess.NewInitialPosition(), 2)

	if !result.HasMove {
		t.Fatal("want best move")
	}

	if got := result.Move.String(); got != "a2a3" {
		t.Fatalf("want a2a3, got %s", got)
	}
}

func TestBestMoveMatchesExactDepthHelper(t *testing.T) {
	tests := []struct {
		name     string
		position chess.Position
		depth    int
	}{
		{
			name:     "initial position",
			position: chess.NewInitialPosition(),
			depth:    3,
		},
		{
			name: "winning capture",
			position: mustBuildPosition(t,
				chess.NewPositionBuilder().
					WithSideToMove(chess.White).
					Place(mustParseSquare(t, "e1"), chess.White, chess.King).
					Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
					Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
					Place(mustParseSquare(t, "a8"), chess.Black, chess.Queen),
			),
			depth: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := bestMoveAtDepth(tt.position, tt.depth, newTranspositionTable())
			got := BestMove(tt.position, tt.depth)

			if got.HasMove != want.HasMove {
				t.Fatalf("want HasMove %t, got %t", want.HasMove, got.HasMove)
			}

			if got.Score != want.Score {
				t.Fatalf("want score %d, got %d", want.Score, got.Score)
			}

			if got.Move != want.Move {
				t.Fatalf("want move %s, got %s", want.Move, got.Move)
			}
		})
	}
}

func TestAlphaBetaPrunesBranchesAndMatchesFullMinimax(t *testing.T) {
	position := chess.NewInitialPosition()
	rootPerspective := position.SideToMove()

	var fullLeafEvaluations int
	wantScore := fullMinimax(position, 2, rootPerspective, &fullLeafEvaluations)

	hooks := &searchHooks{}
	gotScore := alphaBeta(position, 2, -searchInfinity, searchInfinity, rootPerspective, nil, hooks)

	if gotScore != wantScore {
		t.Fatalf("want score %d, got %d", wantScore, gotScore)
	}

	if hooks.cutoffs == 0 {
		t.Fatal("want at least one alpha-beta cutoff")
	}

	if hooks.leafEvaluations >= fullLeafEvaluations {
		t.Fatalf("want fewer leaf evaluations than full minimax, got %d vs %d", hooks.leafEvaluations, fullLeafEvaluations)
	}
}

func TestAlphaBetaOrderingDoesNotIncreaseLeafEvaluations(t *testing.T) {
	position := mustBuildPosition(t,
		chess.NewPositionBuilder().
			WithSideToMove(chess.White).
			Place(mustParseSquare(t, "e1"), chess.White, chess.King).
			Place(mustParseSquare(t, "a1"), chess.White, chess.Rook).
			Place(mustParseSquare(t, "e8"), chess.Black, chess.King).
			Place(mustParseSquare(t, "a8"), chess.Black, chess.Queen),
	)
	rootPerspective := position.SideToMove()

	orderedHooks := &searchHooks{}
	orderedScore := alphaBeta(position, 2, -searchInfinity, searchInfinity, rootPerspective, nil, orderedHooks)

	unorderedHooks := &searchHooks{}
	unorderedScore := alphaBetaUnordered(position, 2, -searchInfinity, searchInfinity, rootPerspective, unorderedHooks)

	if orderedScore != unorderedScore {
		t.Fatalf("want matching score %d, got %d", unorderedScore, orderedScore)
	}

	if orderedHooks.leafEvaluations > unorderedHooks.leafEvaluations {
		t.Fatalf("want ordered leaf evaluations <= unordered, got %d vs %d", orderedHooks.leafEvaluations, unorderedHooks.leafEvaluations)
	}
}

func TestAlphaBetaWithTTMatchesReference(t *testing.T) {
	position := chess.NewInitialPosition()
	rootPerspective := position.SideToMove()

	var fullLeafEvaluations int
	wantScore := fullMinimax(position, 4, rootPerspective, &fullLeafEvaluations)

	tt := newTranspositionTable()
	hooks := &searchHooks{}
	gotScore := alphaBeta(position, 4, -searchInfinity, searchInfinity, rootPerspective, tt, hooks)

	if gotScore != wantScore {
		t.Fatalf("want score %d, got %d", wantScore, gotScore)
	}
}

func TestTranspositionTableReducesLeafEvaluations(t *testing.T) {
	position := chess.NewInitialPosition()
	rootPerspective := position.SideToMove()

	noTTHooks := &searchHooks{}
	noTTScore := alphaBeta(position, 4, -searchInfinity, searchInfinity, rootPerspective, nil, noTTHooks)

	withTTHooks := &searchHooks{}
	withTTScore := alphaBeta(position, 4, -searchInfinity, searchInfinity, rootPerspective, newTranspositionTable(), withTTHooks)

	if withTTScore != noTTScore {
		t.Fatalf("want matching score %d, got %d", noTTScore, withTTScore)
	}

	if withTTHooks.leafEvaluations >= noTTHooks.leafEvaluations {
		t.Fatalf("want TT leaf evaluations < no TT, got %d vs %d", withTTHooks.leafEvaluations, noTTHooks.leafEvaluations)
	}
}

func TestTranspositionTableKeepsDeeperEntry(t *testing.T) {
	tt := newTranspositionTable()
	position := chess.NewInitialPosition()

	tt.store(position, ttEntry{
		depth: 4,
		score: 123,
		bound: ttExact,
	})
	tt.store(position, ttEntry{
		depth: 2,
		score: 456,
		bound: ttLower,
	})

	entry, ok := tt.probe(position)
	if !ok {
		t.Fatal("want TT entry")
	}

	if entry.depth != 4 {
		t.Fatalf("want depth 4, got %d", entry.depth)
	}

	if entry.score != 123 {
		t.Fatalf("want score 123, got %d", entry.score)
	}

	if entry.bound != ttExact {
		t.Fatalf("want exact bound, got %d", entry.bound)
	}
}

func TestOrderedMovesPutsHashMoveFirst(t *testing.T) {
	position := chess.NewInitialPosition()
	hashMove := mustMove(t, "b1", "c3")

	moves := orderedMoves(position, hashMove, true)
	if len(moves) == 0 {
		t.Fatal("want ordered moves")
	}

	if moves[0] != hashMove {
		t.Fatalf("want hash move %s first, got %s", hashMove, moves[0])
	}
}

func TestIterativeDeepeningReusesTTBetweenDepths(t *testing.T) {
	position := chess.NewInitialPosition()
	table := newTranspositionTable()

	depthOne := bestMoveAtDepth(position, 1, table)
	if !depthOne.HasMove {
		t.Fatal("want best move at depth one")
	}

	entry, ok := table.probe(position)
	if !ok {
		t.Fatal("want root TT entry after depth one search")
	}

	if entry.bestMove != depthOne.Move {
		t.Fatalf("want root best move %s, got %s", depthOne.Move, entry.bestMove)
	}

	hooks := &searchHooks{}
	depthTwo := bestMoveAtDepthWithHooks(position, 2, table, hooks)
	if !depthTwo.HasMove {
		t.Fatal("want best move at depth two")
	}

	if hooks.ttHits == 0 {
		t.Fatal("want TT hits on deeper iteration")
	}
}

func fullMinimax(position chess.Position, depth int, rootPerspective chess.Side, leafEvaluations *int) int {
	if depth <= 0 || isTerminalPosition(position) {
		*leafEvaluations = *leafEvaluations + 1
		return Evaluate(position, rootPerspective)
	}

	moves := position.LegalMoves()
	if len(moves) == 0 {
		*leafEvaluations = *leafEvaluations + 1
		return Evaluate(position, rootPerspective)
	}

	maximizing := position.SideToMove() == rootPerspective
	bestScore := 0

	for index, move := range moves {
		next, err := position.ApplyMove(move)
		if err != nil {
			panic(err)
		}

		score := fullMinimax(next, depth-1, rootPerspective, leafEvaluations)
		if index == 0 || betterScore(score, bestScore, maximizing) {
			bestScore = score
		}
	}

	return bestScore
}

func bestMoveSequential(position chess.Position, depth int) SearchResult {
	rootPerspective := position.SideToMove()
	if depth <= 0 || isTerminalPosition(position) {
		return noMoveSearchResult(position, rootPerspective)
	}

	moves := position.LegalMoves()
	if len(moves) == 0 {
		return noMoveSearchResult(position, rootPerspective)
	}

	ordered := reorderMoves(position, moves, chess.Move{}, false)
	moveIndex := indexMoves(moves)
	results := make([]SearchResult, len(moves))

	firstMove := ordered[0]
	firstScore := searchRootMove(position, firstMove, depth-1, rootPerspective, -searchInfinity, nil, nil)
	results[moveIndex[firstMove]] = SearchResult{
		Move:    firstMove,
		Score:   firstScore,
		HasMove: true,
	}

	for _, move := range ordered[1:] {
		score := searchRootMove(position, move, depth-1, rootPerspective, firstScore, nil, nil)
		results[moveIndex[move]] = SearchResult{
			Move:    move,
			Score:   score,
			HasMove: true,
		}
	}

	return pickBestResult(results)
}

func alphaBetaUnordered(position chess.Position, depth int, alpha, beta int, rootPerspective chess.Side, hooks *searchHooks) int {
	if depth <= 0 || isTerminalPosition(position) {
		if hooks != nil {
			hooks.leafEvaluations++
		}

		return Evaluate(position, rootPerspective)
	}

	moves := position.LegalMoves()
	if len(moves) == 0 {
		if hooks != nil {
			hooks.leafEvaluations++
		}

		return Evaluate(position, rootPerspective)
	}

	maximizing := position.SideToMove() == rootPerspective
	bestScore := 0

	for index, move := range moves {
		next, err := position.ApplyMove(move)
		if err != nil {
			panic(err)
		}

		score := alphaBetaUnordered(next, depth-1, alpha, beta, rootPerspective, hooks)
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
			if hooks != nil {
				hooks.cutoffs++
			}

			break
		}
	}

	return bestScore
}

func mustBuildPosition(t *testing.T, builder *chess.PositionBuilder) chess.Position {
	t.Helper()

	position, err := builder.Build()
	if err != nil {
		t.Fatalf("Build error: %v", err)
	}

	return position
}

func mustParseSquare(t *testing.T, raw string) chess.Square {
	t.Helper()

	square, err := chess.ParseSquare(raw)
	if err != nil {
		t.Fatalf("ParseSquare(%q) error: %v", raw, err)
	}

	return square
}

func mustMove(t *testing.T, from, to string, promotion ...chess.PieceType) chess.Move {
	t.Helper()

	move := chess.Move{
		From: mustParseSquare(t, from),
		To:   mustParseSquare(t, to),
	}

	if len(promotion) > 0 {
		move.Promotion = promotion[0]
	}

	return move
}
