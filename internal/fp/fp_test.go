package fp

import (
	"reflect"
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

func TestMap(t *testing.T) {
	t.Run("maps values in order", func(t *testing.T) {
		as := []int{1, 2, 3}
		got := Map(as, func(a int) int { return a * 2 })
		want := []int{2, 4, 6}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("expected %v, got %v", want, got)
		}
	})

	t.Run("nil slice returns empty (non-nil) slice", func(t *testing.T) {
		var as []int
		got := Map(as, func(a int) int { return a * 2 })
		if got == nil || len(got) != 0 {
			t.Fatalf("expected empty non-nil slice, got %v", got)
		}
	})
}

func TestReduce(t *testing.T) {
	t.Run("reduce numbers sum", func(t *testing.T) {
		as := []int{10, 5, 5, 10, 20}
		got := Reduce(as, func(a int, b int) int { return a + b })
		want := 50
		if got != want {
			t.Fatalf("expected %v, got %v", want, got)
		}
	})
	t.Run("reduce string (mkString)", func(t *testing.T) {
		as := []string{"First", "Second", "Third"}
		want := "First\nSecond\nThird"
		got := Reduce(as, func(s string, s2 string) string { return s + "\n" + s2 })
		if want != got {
			t.Fatalf("expected %v, got %v", want, got)
		}
	})
}
