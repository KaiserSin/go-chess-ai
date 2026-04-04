package ai

import chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"

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
	root := BuildTree(position, depth)
	if len(root.Children) == 0 {
		return SearchResult{
			Score:   Evaluate(position, rootPerspective),
			HasMove: false,
		}
	}

	bestScore := minimax(&root, rootPerspective)
	for _, child := range root.Children {
		if child.Score == bestScore {
			return SearchResult{
				Move:    child.Move,
				Score:   child.Score,
				HasMove: true,
			}
		}
	}

	return SearchResult{
		Score:   bestScore,
		HasMove: false,
	}
}

func minimax(node *Node, rootPerspective chess.Side) int {
	if len(node.Children) == 0 {
		node.Score = Evaluate(node.Position, rootPerspective)
		return node.Score
	}

	maximizing := node.Position.SideToMove() == rootPerspective
	bestScore := 0

	for index := range node.Children {
		score := minimax(&node.Children[index], rootPerspective)
		if index == 0 || betterScore(score, bestScore, maximizing) {
			bestScore = score
		}
	}

	node.Score = bestScore
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
