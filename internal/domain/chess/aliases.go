package chess

import (
	chessgame "github.com/KaiserSin/go-chess-ai/internal/domain/chess/game"
	chessmodel "github.com/KaiserSin/go-chess-ai/internal/domain/chess/model"
	chessposition "github.com/KaiserSin/go-chess-ai/internal/domain/chess/position"
)

type Side = chessmodel.Side
type PieceType = chessmodel.PieceType
type Piece = chessmodel.Piece
type Square = chessmodel.Square
type Move = chessmodel.Move
type CastlingRights = chessmodel.CastlingRights
type Status = chessmodel.Status
type OutcomeReason = chessmodel.OutcomeReason
type Outcome = chessmodel.Outcome
type Position = chessposition.Position
type PositionBuilder = chessposition.PositionBuilder
type RepetitionKey = chessposition.RepetitionKey
type Game = chessgame.Game

const (
	White Side = chessmodel.White
	Black Side = chessmodel.Black
)

const (
	NoPieceType PieceType = chessmodel.NoPieceType
	Pawn        PieceType = chessmodel.Pawn
	Knight      PieceType = chessmodel.Knight
	Bishop      PieceType = chessmodel.Bishop
	Rook        PieceType = chessmodel.Rook
	Queen       PieceType = chessmodel.Queen
	King        PieceType = chessmodel.King
)

const (
	WhiteKingSide  CastlingRights = chessmodel.WhiteKingSide
	WhiteQueenSide CastlingRights = chessmodel.WhiteQueenSide
	BlackKingSide  CastlingRights = chessmodel.BlackKingSide
	BlackQueenSide CastlingRights = chessmodel.BlackQueenSide
)

const (
	Ongoing   Status = chessmodel.Ongoing
	Check     Status = chessmodel.Check
	Checkmate Status = chessmodel.Checkmate
	Stalemate Status = chessmodel.Stalemate
)

const (
	NoOutcomeReason               OutcomeReason = chessmodel.NoOutcomeReason
	OutcomeByCheckmate            OutcomeReason = chessmodel.OutcomeByCheckmate
	OutcomeByStalemate            OutcomeReason = chessmodel.OutcomeByStalemate
	OutcomeByThreefoldRepetition  OutcomeReason = chessmodel.OutcomeByThreefoldRepetition
	OutcomeByFiftyMoveRule        OutcomeReason = chessmodel.OutcomeByFiftyMoveRule
	OutcomeByInsufficientMaterial OutcomeReason = chessmodel.OutcomeByInsufficientMaterial
)
