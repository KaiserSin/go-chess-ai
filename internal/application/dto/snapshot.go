package dto

type GameSnapshot struct {
	Squares        []SquareSnapshot
	SideToMove     string
	Status         string
	OutcomeReason  string
	Winner         string
	HasWinner      bool
	HalfmoveClock  int
	FullmoveNumber int
}

type SquareSnapshot struct {
	File      int
	Rank      int
	Algebraic string
	Occupied  bool
	PieceKey  string
}
