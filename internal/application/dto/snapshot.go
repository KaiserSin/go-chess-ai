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
	Promotion      *PromotionSnapshot
}

type SquareSnapshot struct {
	File        int
	Rank        int
	Algebraic   string
	Occupied    bool
	PieceKey    string
	Selected    bool
	LegalTarget bool
}

type PromotionSnapshot struct {
	Visible      bool
	TargetSquare string
	Options      []PromotionOptionSnapshot
}

type PromotionOptionSnapshot struct {
	PieceType string
	PieceKey  string
}
