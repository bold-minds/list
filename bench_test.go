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

func BenchmarkUnique_Small(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = list.Unique(smallA)
	}
}

func BenchmarkUnique_WithDupes(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = list.Unique(dupes)
	}
}

func BenchmarkUnion_Two(b *testing.B) {
	b.ReportAllocs()
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

func BenchmarkIntersect_Two(b *testing.B) {
	b.ReportAllocs()
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

func BenchmarkSymmetricDifference(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = list.SymmetricDifference(smallA, smallB)
	}
}

func BenchmarkMinus_Basic(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = list.Minus(smallA, smallB)
	}
}

func BenchmarkWithout_SingleItem(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = list.Without(smallA, 5)
	}
}

func BenchmarkWithout_MultipleItems(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = list.Without(smallA, 1, 3, 5, 7, 9)
	}
}
