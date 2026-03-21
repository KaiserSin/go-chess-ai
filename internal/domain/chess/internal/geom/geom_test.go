package geom

import "testing"

func TestRanks(t *testing.T) {
	if BackRank(true) != 0 {
		t.Fatal("want white back rank 0")
	}

	if BackRank(false) != 7 {
		t.Fatal("want black back rank 7")
	}

	if PawnStartRank(true) != 1 {
		t.Fatal("want white pawn start 1")
	}

	if PawnStartRank(false) != 6 {
		t.Fatal("want black pawn start 6")
	}

	if PawnPromotionRank(true) != 7 {
		t.Fatal("want white promo 7")
	}

	if PawnPromotionRank(false) != 0 {
		t.Fatal("want black promo 0")
	}
}

func TestPawnDirection(t *testing.T) {
	if PawnDirection(true) != 1 {
		t.Fatal("want white dir 1")
	}

	if PawnDirection(false) != -1 {
		t.Fatal("want black dir -1")
	}
}

func TestOffsets(t *testing.T) {
	if len(PawnCaptureFiles) != 2 {
		t.Fatal("want 2 pawn files")
	}

	if len(KnightOffsets) != 8 {
		t.Fatal("want 8 knight moves")
	}

	if len(KingOffsets) != 8 {
		t.Fatal("want 8 king moves")
	}

	if len(BishopDirections) != 4 {
		t.Fatal("want 4 bishop dirs")
	}

	if len(RookDirections) != 4 {
		t.Fatal("want 4 rook dirs")
	}

	if len(QueenDirections) != 8 {
		t.Fatal("want 8 queen dirs")
	}
}
