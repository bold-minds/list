package list_test

import (
	"testing"

	"github.com/bold-minds/list"
)

var (
	smallA = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	smallB = []int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	dupes  = []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 5}
)

// makeRange returns a slice of n sequential integers starting at start.
func makeRange(start, n int) []int {
	out := make([]int, n)
	for i := 0; i < n; i++ {
		out[i] = start + i
	}
	return out
}

// makeRangeDupes returns a slice of length n where each value is repeated
// `repeat` times. Used to stress the dedup path on large inputs.
func makeRangeDupes(n, repeat int) []int {
	out := make([]int, 0, n)
	for i := 0; len(out) < n; i++ {
		for j := 0; j < repeat && len(out) < n; j++ {
			out = append(out, i)
		}
	}
	return out
}

var (
	large1K   = makeRange(0, 1000)
	large1Kb  = makeRange(500, 1000) // 50% overlap with large1K
	large10K  = makeRange(0, 10000)
	large10Kb = makeRange(5000, 10000)
	dupes1K   = makeRangeDupes(1000, 4) // ~250 unique values
)

func BenchmarkUnique_Small(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Unique(smallA)
	}
}

func BenchmarkUnique_WithDupes(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Unique(dupes)
	}
}

func BenchmarkUnique_1K(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Unique(large1K)
	}
}

func BenchmarkUnique_1K_WithDupes(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Unique(dupes1K)
	}
}

func BenchmarkUnique_10K(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Unique(large10K)
	}
}

func BenchmarkUnion_Two(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Union(smallA, smallB)
	}
}

func BenchmarkUnion_Three(b *testing.B) {
	third := []int{15, 16, 17, 18}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Union(smallA, smallB, third)
	}
}

func BenchmarkUnion_WithDupes(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Union(dupes, dupes)
	}
}

func BenchmarkUnion_1K(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Union(large1K, large1Kb)
	}
}

func BenchmarkUnion_10K(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Union(large10K, large10Kb)
	}
}

func BenchmarkIntersect_Two(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Intersect(smallA, smallB)
	}
}

func BenchmarkIntersect_Three(b *testing.B) {
	third := []int{3, 4, 5, 6, 7}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Intersect(smallA, smallB, third)
	}
}

func BenchmarkIntersect_WithDupes(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Intersect(dupes, dupes)
	}
}

func BenchmarkIntersect_1K(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Intersect(large1K, large1Kb)
	}
}

func BenchmarkIntersect_10K(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Intersect(large10K, large10Kb)
	}
}

func BenchmarkSymmetricDifference(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.SymmetricDifference(smallA, smallB)
	}
}

func BenchmarkSymmetricDifference_1K(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.SymmetricDifference(large1K, large1Kb)
	}
}

func BenchmarkMinus_Basic(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Minus(smallA, smallB)
	}
}

func BenchmarkMinus_WithDupes(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Minus(dupes, smallB)
	}
}

func BenchmarkMinus_1K(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Minus(large1K, large1Kb)
	}
}

func BenchmarkWithout_SingleItem(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Without(smallA, 5)
	}
}

func BenchmarkWithout_MultipleItems(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Without(smallA, 1, 3, 5, 7, 9)
	}
}

func BenchmarkWithout_1K(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = list.Without(large1K, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19)
	}
}
