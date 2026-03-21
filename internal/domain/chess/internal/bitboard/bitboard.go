package bitboard

import "math/bits"

func First(value uint64) (int, bool) {
	if value == 0 {
		return 0, false
	}

	return bits.TrailingZeros64(value), true
}

func ForEach(value uint64, visit func(int)) {
	for value != 0 {
		index := bits.TrailingZeros64(value)
		visit(index)
		value &= value - 1
	}
}

func Count(value uint64) int {
	return bits.OnesCount64(value)
}
