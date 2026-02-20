package fp

import (
	"testing"
)

func TestFoldLeft(t *testing.T) {
	t.Run("nil slice returns acc", func(t *testing.T) {
		var as []int
		got := FoldLeft(as, 10, func(a, acc int) int { return acc + a })
		if got != 10 {
			t.Fatalf("expected acc 10, got %d", got)
		}
	})

	t.Run("empty slice returns acc", func(t *testing.T) {
		as := []int{}
		got := FoldLeft(as, 7, func(a, acc int) int { return acc * a })
		if got != 7 {
			t.Fatalf("expected acc 7, got %d", got)
		}
	})

	t.Run("applies left to right", func(t *testing.T) {
		as := []int{1, 2, 3}
		got := FoldLeft(as, 0, func(a, acc int) int { return acc*10 + a })
		if got != 123 {
			t.Fatalf("expected 123, got %d", got)
		}
	})
}

func TestForAll(t *testing.T) {
	t.Run("all match returns true", func(t *testing.T) {
		as := []int{2, 4, 6}
		got := ForAll(as, func(a int) bool { return a%2 == 0 })
		if !got {
			t.Fatal("expected true for all even numbers")
		}
	})

	t.Run("one mismatch returns false", func(t *testing.T) {
		as := []int{2, 3, 6}
		got := ForAll(as, func(a int) bool { return a%2 == 0 })
		if got {
			t.Fatal("expected false when one element is odd")
		}
	})

	t.Run("empty slice returns true", func(t *testing.T) {
		var as []int
		got := ForAll(as, func(a int) bool { return false })
		if !got {
			t.Fatal("expected true for empty slice")
		}
	})
}
