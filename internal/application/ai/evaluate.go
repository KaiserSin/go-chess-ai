package ai

import chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"

const mateScore = 100000
const endgameNonPawnMaterialLimit = 2600

const doubledPawnPenalty = 12
const isolatedPawnPenalty = 10
const passedPawnBonus = 18
const kingMovedPenalty = 20
const kingCenterPenalty = 10
const kingShieldPawnBonus = 8
const endgameKingActivityBonus = 3
const endgamePassedPawnStepBonus = 4

var pawnPieceSquareTable = [8][8]int{
	{0, 0, 0, 0, 0, 0, 0, 0},
	{5, 5, 5, 5, 5, 5, 5, 5},
	{4, 4, 6, 8, 8, 6, 4, 4},
	{2, 3, 6, 10, 10, 6, 3, 2},
	{2, 3, 6, 10, 10, 6, 3, 2},
	{3, 4, 8, 12, 12, 8, 4, 3},
	{6, 6, 10, 14, 14, 10, 6, 6},
	{0, 0, 0, 0, 0, 0, 0, 0},
}

var knightPieceSquareTable = [8][8]int{
	{-20, -10, -10, -10, -10, -10, -10, -20},
	{-10, 0, 0, 5, 5, 0, 0, -10},
	{-10, 5, 10, 10, 10, 10, 5, -10},
	{-10, 0, 10, 15, 15, 10, 0, -10},
	{-10, 5, 10, 15, 15, 10, 5, -10},
	{-10, 0, 5, 10, 10, 5, 0, -10},
	{-10, 0, 0, 0, 0, 0, 0, -10},
	{-20, -10, -10, -10, -10, -10, -10, -20},
}

var bishopPieceSquareTable = [8][8]int{
	{-10, -5, -5, -5, -5, -5, -5, -10},
	{-5, 5, 0, 0, 0, 0, 5, -5},
	{-5, 10, 10, 10, 10, 10, 10, -5},
	{-5, 0, 10, 10, 10, 10, 0, -5},
	{-5, 5, 5, 10, 10, 5, 5, -5},
	{-5, 0, 5, 10, 10, 5, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-10, -5, -5, -5, -5, -5, -5, -10},
}

var rookPieceSquareTable = [8][8]int{
	{0, 0, 0, 5, 5, 0, 0, 0},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-5, 0, 5, 5, 5, 5, 0, -5},
	{-5, 0, 5, 5, 5, 5, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{5, 10, 10, 10, 10, 10, 10, 5},
	{0, 0, 0, 0, 0, 0, 0, 0},
}

var queenPieceSquareTable = [8][8]int{
	{-10, -5, -5, 0, 0, -5, -5, -10},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-5, 0, 5, 5, 5, 5, 0, -5},
	{0, 0, 5, 5, 5, 5, 0, -5},
	{-5, 0, 5, 5, 5, 5, 0, -5},
	{-5, 0, 5, 5, 5, 5, 0, -5},
	{-5, 0, 0, 0, 0, 0, 0, -5},
	{-10, -5, -5, 0, 0, -5, -5, -10},
}

var kingMiddlegameTable = [8][8]int{
	{20, 25, 10, 0, 0, 10, 25, 20},
	{20, 20, 0, 0, 0, 0, 20, 20},
	{-10, -15, -20, -20, -20, -20, -15, -10},
	{-20, -25, -30, -35, -35, -30, -25, -20},
	{-25, -30, -35, -40, -40, -35, -30, -25},
	{-30, -35, -40, -45, -45, -40, -35, -30},
	{-35, -40, -45, -50, -50, -45, -40, -35},
	{-35, -40, -45, -50, -50, -45, -40, -35},
}

var kingEndgameTable = [8][8]int{
	{-20, -10, -10, -10, -10, -10, -10, -20},
	{-10, 0, 0, 0, 0, 0, 0, -10},
	{-10, 0, 10, 15, 15, 10, 0, -10},
	{-10, 0, 15, 20, 20, 15, 0, -10},
	{-10, 0, 15, 20, 20, 15, 0, -10},
	{-10, 0, 10, 15, 15, 10, 0, -10},
	{-10, 0, 0, 0, 0, 0, 0, -10},
	{-20, -10, -10, -10, -10, -10, -10, -20},
}

func Evaluate(position chess.Position, perspective chess.Side) int {
	if perspective != chess.White && perspective != chess.Black {
		return 0
	}

	switch position.Status() {
	case chess.Checkmate:
		if perspective == position.SideToMove() {
			return -mateScore
		}

		return mateScore
	case chess.Stalemate:
		return 0
	}

	if chess.HasInsufficientMaterial(position) || chess.IsFiftyMoveDraw(position) {
		return 0
	}

	info := collectEvalInfo(position)
	score := 0
	score += materialScore(info)
	score += pieceSquareScore(info)
	score += pawnStructureScore(info)
	score += kingSafetyScore(position, info)
	score += endgameBonusScore(info)

	if perspective == chess.Black {
		return -score
	}

	return score
}

func materialScore(info evalInfo) int {
	return info.material
}

func pieceSquareScore(info evalInfo) int {
	return info.pieceSquares
}

func pawnStructureScore(info evalInfo) int {
	score := 0

	score += signedScore(chess.White, sidePawnStructureScore(chess.White, info.whitePawns, info.whitePawnFiles, info.blackPawns))
	score += signedScore(chess.Black, sidePawnStructureScore(chess.Black, info.blackPawns, info.blackPawnFiles, info.whitePawns))

	return score
}

func kingSafetyScore(position chess.Position, info evalInfo) int {
	if info.isEndgame() {
		return 0
	}

	score := 0

	score += signedScore(chess.White, sideKingSafetyScore(position, chess.White, info.whiteKing))
	score += signedScore(chess.Black, sideKingSafetyScore(position, chess.Black, info.blackKing))

	return score
}

func endgameBonusScore(info evalInfo) int {
	if !info.isEndgame() {
		return 0
	}

	score := 0

	score += signedScore(chess.White, kingActivityBonus(info.whiteKing))
	score += signedScore(chess.Black, kingActivityBonus(info.blackKing))

	for _, pawn := range info.whitePawns {
		if isPassedPawn(pawn, chess.White, info.blackPawns) {
			score += endgamePassedPawnBonus(chess.White, pawn)
		}
	}

	for _, pawn := range info.blackPawns {
		if isPassedPawn(pawn, chess.Black, info.whitePawns) {
			score -= endgamePassedPawnBonus(chess.Black, pawn)
		}
	}

	return score
}

func pieceValue(pieceType chess.PieceType) int {
	switch pieceType {
	case chess.Pawn:
		return 100
	case chess.Knight:
		return 320
	case chess.Bishop:
		return 330
	case chess.Rook:
		return 500
	case chess.Queen:
		return 900
	default:
		return 0
	}
}

func pieceSquareValue(pieceType chess.PieceType, side chess.Side, file, rank int, endgame bool) int {
	lookupRank := rank
	if side == chess.Black {
		lookupRank = 7 - rank
	}

	switch pieceType {
	case chess.Pawn:
		return pawnPieceSquareTable[lookupRank][file]
	case chess.Knight:
		return knightPieceSquareTable[lookupRank][file]
	case chess.Bishop:
		return bishopPieceSquareTable[lookupRank][file]
	case chess.Rook:
		return rookPieceSquareTable[lookupRank][file]
	case chess.Queen:
		return queenPieceSquareTable[lookupRank][file]
	case chess.King:
		if endgame {
			return kingEndgameTable[lookupRank][file]
		}

		return kingMiddlegameTable[lookupRank][file]
	default:
		return 0
	}
}

type evalInfo struct {
	whitePawns      []chess.Square
	blackPawns      []chess.Square
	whitePawnFiles  [8]int
	blackPawnFiles  [8]int
	whiteKing       chess.Square
	blackKing       chess.Square
	material        int
	pieceSquares    int
	nonPawnMaterial int
}

func collectEvalInfo(position chess.Position) evalInfo {
	info := evalInfo{}

	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			square := squareMust(file, rank)
			piece, ok := position.PieceAt(square)
			if !ok {
				continue
			}

			info.material += signedScore(piece.Side(), pieceValue(piece.Type()))

			switch piece.Type() {
			case chess.Pawn:
				if piece.Side() == chess.White {
					info.whitePawns = append(info.whitePawns, square)
					info.whitePawnFiles[file]++
				} else {
					info.blackPawns = append(info.blackPawns, square)
					info.blackPawnFiles[file]++
				}
				info.pieceSquares += signedScore(piece.Side(), pieceSquareValue(piece.Type(), piece.Side(), file, rank, false))
			case chess.King:
				if piece.Side() == chess.White {
					info.whiteKing = square
				} else {
					info.blackKing = square
				}
			default:
				info.nonPawnMaterial += pieceValue(piece.Type())
				info.pieceSquares += signedScore(piece.Side(), pieceSquareValue(piece.Type(), piece.Side(), file, rank, false))
			}
		}
	}

	endgame := info.isEndgame()
	info.pieceSquares += signedScore(chess.White, pieceSquareValue(chess.King, chess.White, info.whiteKing.File(), info.whiteKing.Rank(), endgame))
	info.pieceSquares += signedScore(chess.Black, pieceSquareValue(chess.King, chess.Black, info.blackKing.File(), info.blackKing.Rank(), endgame))

	return info
}

func (info evalInfo) isEndgame() bool {
	return info.nonPawnMaterial <= endgameNonPawnMaterialLimit
}

func sidePawnStructureScore(side chess.Side, pawns []chess.Square, fileCounts [8]int, opponentPawns []chess.Square) int {
	score := 0

	for file := 0; file < 8; file++ {
		if fileCounts[file] > 1 {
			score -= doubledPawnPenalty * (fileCounts[file] - 1)
		}
	}

	for _, pawn := range pawns {
		file := pawn.File()
		if !hasAdjacentPawnFile(fileCounts, file) {
			score -= isolatedPawnPenalty
		}

		if isPassedPawn(pawn, side, opponentPawns) {
			score += passedPawnBonus
		}
	}

	return score
}

func hasAdjacentPawnFile(fileCounts [8]int, file int) bool {
	if file > 0 && fileCounts[file-1] > 0 {
		return true
	}

	if file < 7 && fileCounts[file+1] > 0 {
		return true
	}

	return false
}

func isPassedPawn(pawn chess.Square, side chess.Side, opponentPawns []chess.Square) bool {
	for _, opponentPawn := range opponentPawns {
		fileDiff := absInt(opponentPawn.File() - pawn.File())
		if fileDiff > 1 {
			continue
		}

		if side == chess.White && opponentPawn.Rank() > pawn.Rank() {
			return false
		}

		if side == chess.Black && opponentPawn.Rank() < pawn.Rank() {
			return false
		}
	}

	return true
}

func sideKingSafetyScore(position chess.Position, side chess.Side, king chess.Square) int {
	score := pawnShieldBonus(position, side, king)

	if kingMovedAwayFromHome(side, king) {
		score -= kingMovedPenalty
	}

	if kingInCenter(king) {
		score -= kingCenterPenalty
	}

	return score
}

func pawnShieldBonus(position chess.Position, side chess.Side, king chess.Square) int {
	nextRank := king.Rank() + pawnDirection(side)
	if nextRank < 0 || nextRank > 7 {
		return 0
	}

	score := 0
	for file := maxInt(0, king.File()-1); file <= minInt(7, king.File()+1); file++ {
		piece, ok := position.PieceAt(squareMust(file, nextRank))
		if ok && piece.Side() == side && piece.Type() == chess.Pawn {
			score += kingShieldPawnBonus
		}
	}

	return score
}

func kingMovedAwayFromHome(side chess.Side, king chess.Square) bool {
	if side == chess.White {
		return king.Rank() > 1
	}

	return king.Rank() < 6
}

func kingInCenter(king chess.Square) bool {
	return king.File() >= 2 && king.File() <= 5 && king.Rank() >= 2 && king.Rank() <= 5
}

func kingActivityBonus(king chess.Square) int {
	fileDistance := minInt(absInt(king.File()-3), absInt(king.File()-4))
	rankDistance := minInt(absInt(king.Rank()-3), absInt(king.Rank()-4))
	score := 12 - endgameKingActivityBonus*(fileDistance+rankDistance)
	if score < 0 {
		return 0
	}

	return score
}

func endgamePassedPawnBonus(side chess.Side, pawn chess.Square) int {
	return pawnAdvance(side, pawn.Rank()) * endgamePassedPawnStepBonus
}

func pawnAdvance(side chess.Side, rank int) int {
	if side == chess.White {
		return rank
	}

	return 7 - rank
}

func signedScore(side chess.Side, value int) int {
	if side == chess.White {
		return value
	}

	return -value
}

func squareMust(file, rank int) chess.Square {
	square, err := chess.NewSquare(file, rank)
	if err != nil {
		panic(err)
	}

	return square
}

func minInt(left, right int) int {
	if left < right {
		return left
	}

	return right
}

func maxInt(left, right int) int {
	if left > right {
		return left
	}

	return right
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}

	return value
}

func pawnDirection(side chess.Side) int {
	if side == chess.White {
		return 1
	}

	return -1
}
