package bitboard

import "testing"

func TestFirst(t *testing.T) {
	if _, ok := First(0); ok {
		t.Fatal("want no bit")
	}

	index, ok := First(0b1001000)
	if !ok || index != 3 {
		t.Fatalf("want 3, got %d, %t", index, ok)
	}
}

func TestForEach(t *testing.T) {
	var got []int

	ForEach(0b10101, func(index int) {
		got = append(got, index)
	})

	want := []int{0, 2, 4}
	if len(got) != len(want) {
		t.Fatalf("want %d bits, got %d", len(want), len(got))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("want %v, got %v", want, got)
		}
	}
}

func TestCount(t *testing.T) {
	if got := Count(0b101101); got != 4 {
		t.Fatalf("want 4, got %d", got)
	}
}
