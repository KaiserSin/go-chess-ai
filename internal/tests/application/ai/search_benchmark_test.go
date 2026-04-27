package ai_test

import (
	"testing"
	"time"

	"github.com/KaiserSin/go-chess-ai/internal/application/ai"
	chess "github.com/KaiserSin/go-chess-ai/internal/domain/chess"
)

var bestMoveBenchmarkResult ai.SearchResult

func BenchmarkBestMove(b *testing.B) {
	cases := []struct {
		name     string
		position chess.Position
	}{
		{name: "InitialPosition", position: chess.NewInitialPosition()},
		{
			name: "OpeningFourMoves",
			position: mustPositionFromMoves(b,
				mustMove(b, "e2", "e4"),
				mustMove(b, "e7", "e5"),
				mustMove(b, "g1", "f3"),
				mustMove(b, "b8", "c6"),
			),
		},
		{name: "PoisonedCapture", position: poisonedCapturePosition(b)},
		{name: "ForcedMateDepthThree", position: forcedMateDepthThreePosition(b)},
		{name: "PromotionReady", position: promotionReadyPosition(b)},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bestMoveBenchmarkResult = ai.BestMoveWithin(tc.position, 100*time.Millisecond)
			}
		})
	}
}
