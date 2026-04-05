package list_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/bold-minds/list"
)

// =============================================================================
// Unique
// =============================================================================

func TestUnique_WithDuplicates(t *testing.T) {
	got := list.Unique([]int{1, 2, 2, 3, 1, 4, 2})
	want := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestUnique_NoDuplicates(t *testing.T) {
	got := list.Unique([]int{1, 2, 3, 4})
	want := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestUnique_Strings(t *testing.T) {
	got := list.Unique([]string{"go", "web", "api", "go", "web"})
	want := []string{"go", "web", "api"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestUnique_Empty(t *testing.T) {
	got := list.Unique([]int{})
	if got == nil {
		t.Error("expected non-nil empty slice")
	}
	if len(got) != 0 {
		t.Errorf("expected length 0, got %d", len(got))
	}
}

func TestUnique_Nil(t *testing.T) {
	got := list.Unique[int](nil)
	if got == nil {
		t.Error("expected non-nil empty slice")
	}
	if len(got) != 0 {
		t.Errorf("expected length 0, got %d", len(got))
	}
}

func TestUnique_SingleElement(t *testing.T) {
	got := list.Unique([]int{42})
	if !reflect.DeepEqual(got, []int{42}) {
		t.Errorf("got %v, want [42]", got)
	}
}

func TestUnique_AllSame(t *testing.T) {
	got := list.Unique([]int{5, 5, 5, 5})
	if !reflect.DeepEqual(got, []int{5}) {
		t.Errorf("got %v, want [5]", got)
	}
}

func TestUnique_PreservesOrder(t *testing.T) {
	// Dedup must preserve order of first occurrence
	got := list.Unique([]int{3, 1, 2, 1, 3, 2, 4})
	want := []int{3, 1, 2, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v (order-of-first-occurrence)", got, want)
	}
}

// =============================================================================
// Union
// =============================================================================

func TestUnion_TwoSlices(t *testing.T) {
	got := list.Union([]int{1, 2, 3}, []int{3, 4, 5})
	want := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestUnion_ThreeSlices(t *testing.T) {
	got := list.Union([]int{1, 2}, []int{2, 3}, []int{3, 4})
	want := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestUnion_ZeroSlices(t *testing.T) {
	got := list.Union[int]()
	if got == nil {
		t.Error("expected non-nil empty slice")
	}
	if len(got) != 0 {
		t.Errorf("expected length 0, got %d", len(got))
	}
}

func TestUnion_OneSlice(t *testing.T) {
	got := list.Union([]int{1, 2, 2, 3})
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestUnion_DuplicatesWithinSlices(t *testing.T) {
	got := list.Union([]int{1, 1, 2}, []int{2, 2, 3})
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestUnion_NilSlices(t *testing.T) {
	got := list.Union[int](nil, nil)
	if got == nil {
		t.Error("expected non-nil empty slice")
	}
	if len(got) != 0 {
		t.Errorf("expected length 0, got %d", len(got))
	}
}

func TestUnion_NilWithNonNil(t *testing.T) {
	got := list.Union[int](nil, []int{1, 2, 3})
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestUnion_PreservesOrder(t *testing.T) {
	got := list.Union([]int{3, 1}, []int{2, 1, 4}, []int{5, 3})
	want := []int{3, 1, 2, 4, 5}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// =============================================================================
// Intersect
// =============================================================================

func TestIntersect_TwoSlices(t *testing.T) {
	got := list.Intersect([]int{1, 2, 3, 4}, []int{2, 3, 5})
	want := []int{2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestIntersect_ThreeSlices(t *testing.T) {
	got := list.Intersect([]int{1, 2, 3}, []int{2, 3, 4}, []int{3, 4, 5})
	want := []int{3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestIntersect_NoOverlap(t *testing.T) {
	got := list.Intersect([]int{1, 2}, []int{3, 4})
	if len(got) != 0 {
		t.Errorf("got %v, want empty slice", got)
	}
}

func TestIntersect_ZeroSlices(t *testing.T) {
	got := list.Intersect[int]()
	if len(got) != 0 {
		t.Errorf("got %v, want empty slice", got)
	}
}

func TestIntersect_OneSlice(t *testing.T) {
	// Single-slice call is equivalent to Unique
	got := list.Intersect([]int{1, 2, 2, 3})
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestIntersect_DuplicatesWithin(t *testing.T) {
	// Duplicates within a single slice should not affect the count
	got := list.Intersect([]int{1, 1, 2, 2, 3}, []int{1, 2, 4})
	want := []int{1, 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestIntersect_OrderFromFirstSlice(t *testing.T) {
	got := list.Intersect([]int{4, 3, 2, 1}, []int{1, 2, 3, 4})
	want := []int{4, 3, 2, 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestIntersect_EmptySlice(t *testing.T) {
	got := list.Intersect([]int{1, 2, 3}, []int{})
	if len(got) != 0 {
		t.Errorf("got %v, want empty slice", got)
	}
}

// =============================================================================
// Minus
// =============================================================================

func TestMinus_Basic(t *testing.T) {
	got := list.Minus([]int{1, 2, 3, 4}, []int{2, 4})
	want := []int{1, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestMinus_Strings(t *testing.T) {
	got := list.Minus([]string{"alice", "bob", "carol", "dave"}, []string{"bob", "dave"})
	want := []string{"alice", "carol"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestMinus_NothingRemoved(t *testing.T) {
	got := list.Minus([]int{1, 2, 3}, []int{})
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestMinus_AllRemoved(t *testing.T) {
	got := list.Minus([]int{1, 2, 3}, []int{1, 2, 3})
	if len(got) != 0 {
		t.Errorf("got %v, want empty slice", got)
	}
}

func TestMinus_EmptyA(t *testing.T) {
	got := list.Minus([]int{}, []int{1, 2})
	if len(got) != 0 {
		t.Errorf("got %v, want empty slice", got)
	}
}

func TestMinus_NilA(t *testing.T) {
	got := list.Minus[int](nil, []int{1, 2})
	if len(got) != 0 {
		t.Errorf("got %v, want empty slice", got)
	}
}

func TestMinus_NilB(t *testing.T) {
	got := list.Minus([]int{1, 2, 3}, nil)
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestMinus_DeduplicatesA(t *testing.T) {
	// Minus deduplicates a while removing b's elements
	got := list.Minus([]int{1, 2, 2, 3, 1, 4}, []int{3})
	want := []int{1, 2, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// =============================================================================
// Without
// =============================================================================

func TestWithout_Basic(t *testing.T) {
	got := list.Without([]int{1, 2, 3, 4}, 2, 4)
	want := []int{1, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestWithout_PreservesDuplicates(t *testing.T) {
	// Unlike other list functions, Without does NOT deduplicate
	got := list.Without([]int{1, 2, 3, 2, 1}, 3)
	want := []int{1, 2, 2, 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v (Without preserves remaining duplicates)", got, want)
	}
}

func TestWithout_NothingToRemove(t *testing.T) {
	got := list.Without([]int{1, 2, 3}, 99)
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestWithout_NoItems(t *testing.T) {
	// Zero items to remove returns a copy of the input
	got := list.Without([]int{1, 2, 3})
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestWithout_AllRemoved(t *testing.T) {
	got := list.Without([]int{1, 2, 3}, 1, 2, 3)
	if len(got) != 0 {
		t.Errorf("got %v, want empty slice", got)
	}
}

func TestWithout_EmptySlice(t *testing.T) {
	got := list.Without([]int{}, 1, 2)
	if got == nil {
		t.Error("expected non-nil empty slice")
	}
	if len(got) != 0 {
		t.Errorf("got length %d, want 0", len(got))
	}
}

func TestWithout_NilSlice(t *testing.T) {
	got := list.Without[int](nil, 1, 2)
	if got == nil {
		t.Error("expected non-nil empty slice")
	}
	if len(got) != 0 {
		t.Errorf("got length %d, want 0", len(got))
	}
}

func TestWithout_Strings(t *testing.T) {
	got := list.Without([]string{"alice", "bob", "carol"}, "bob")
	want := []string{"alice", "carol"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// =============================================================================
// Adversarial — edge cases a deep code review should catch
// =============================================================================

// TestNilReturnGuarantees verifies every function returns non-nil slices
// on every input variant that's documented as returning "empty slice".
// A regression here would break "safe to range without nil check".
func TestNilReturnGuarantees(t *testing.T) {
	// Unique
	if got := list.Unique[int](nil); got == nil {
		t.Error("Unique(nil) returned nil")
	}
	if got := list.Unique([]int{}); got == nil {
		t.Error("Unique([]) returned nil")
	}

	// Union
	if got := list.Union[int](); got == nil {
		t.Error("Union() returned nil")
	}
	if got := list.Union[int](nil); got == nil {
		t.Error("Union(nil) returned nil")
	}
	if got := list.Union[int](nil, nil); got == nil {
		t.Error("Union(nil, nil) returned nil")
	}

	// Intersect
	if got := list.Intersect[int](); got == nil {
		t.Error("Intersect() returned nil")
	}
	if got := list.Intersect[int](nil); got == nil {
		t.Error("Intersect(nil) returned nil")
	}
	if got := list.Intersect([]int{1, 2}, []int{3, 4}); got == nil {
		t.Error("Intersect with no overlap returned nil")
	}

	// Minus
	if got := list.Minus[int](nil, nil); got == nil {
		t.Error("Minus(nil, nil) returned nil")
	}
	if got := list.Minus([]int{1, 2}, []int{1, 2}); got == nil {
		t.Error("Minus with all removed returned nil")
	}

	// Without
	if got := list.Without[int](nil, 1); got == nil {
		t.Error("Without(nil, ...) returned nil")
	}
	if got := list.Without([]int{}, 1); got == nil {
		t.Error("Without([], ...) returned nil")
	}
	if got := list.Without([]int{1, 2}, 1, 2); got == nil {
		t.Error("Without with all removed returned nil")
	}
}

// TestImmutability verifies that functions never mutate their inputs.
// This is a core safety guarantee — a regression would cause spooky
// action at a distance in caller code.
func TestImmutability(t *testing.T) {
	// Unique
	in := []int{1, 2, 2, 3, 1}
	snapshot := append([]int{}, in...)
	_ = list.Unique(in)
	if !reflect.DeepEqual(in, snapshot) {
		t.Errorf("Unique mutated input: %v vs %v", in, snapshot)
	}

	// Union
	a := []int{1, 2, 3}
	b := []int{3, 4, 5}
	aSnap := append([]int{}, a...)
	bSnap := append([]int{}, b...)
	_ = list.Union(a, b)
	if !reflect.DeepEqual(a, aSnap) || !reflect.DeepEqual(b, bSnap) {
		t.Errorf("Union mutated input")
	}

	// Intersect
	a2 := []int{1, 2, 3, 4}
	b2 := []int{2, 3, 5}
	aSnap2 := append([]int{}, a2...)
	bSnap2 := append([]int{}, b2...)
	_ = list.Intersect(a2, b2)
	if !reflect.DeepEqual(a2, aSnap2) || !reflect.DeepEqual(b2, bSnap2) {
		t.Errorf("Intersect mutated input")
	}

	// Minus
	a3 := []int{1, 2, 3, 4}
	b3 := []int{2, 4}
	aSnap3 := append([]int{}, a3...)
	bSnap3 := append([]int{}, b3...)
	_ = list.Minus(a3, b3)
	if !reflect.DeepEqual(a3, aSnap3) || !reflect.DeepEqual(b3, bSnap3) {
		t.Errorf("Minus mutated input")
	}

	// Without
	in4 := []int{1, 2, 3, 2, 1}
	snap4 := append([]int{}, in4...)
	_ = list.Without(in4, 2)
	if !reflect.DeepEqual(in4, snap4) {
		t.Errorf("Without mutated input")
	}
}

// TestResultIsNotAliased verifies that mutating the returned slice does
// not affect the input slice. A regression here would be a subtle bug
// where the result shares backing storage with the input.
func TestResultIsNotAliased(t *testing.T) {
	// Without with zero items returns a copy — verify it's not aliased
	in := []int{1, 2, 3}
	out := list.Without(in)
	out[0] = 999
	if in[0] == 999 {
		t.Error("Without() returned a slice aliased to input")
	}

	// Unique on an already-unique slice — verify not aliased
	in2 := []int{1, 2, 3}
	out2 := list.Unique(in2)
	if len(out2) > 0 {
		out2[0] = 999
		if in2[0] == 999 {
			t.Error("Unique returned a slice aliased to input")
		}
	}

	// Union — mutating output must not touch any input slice
	ua := []int{1, 2, 3}
	ub := []int{4, 5, 6}
	uOut := list.Union(ua, ub)
	uOut[0] = 999
	if ua[0] == 999 || ub[0] == 999 {
		t.Error("Union returned a slice aliased to an input")
	}

	// Intersect — mutating output must not touch any input slice
	ia := []int{1, 2, 3, 4}
	ib := []int{2, 3, 5}
	iOut := list.Intersect(ia, ib)
	if len(iOut) > 0 {
		iOut[0] = 999
		if ia[1] == 999 || ib[0] == 999 {
			t.Error("Intersect returned a slice aliased to an input")
		}
	}

	// Minus — mutating output must not touch input a
	ma := []int{1, 2, 3, 4}
	mb := []int{2, 4}
	mOut := list.Minus(ma, mb)
	if len(mOut) > 0 {
		mOut[0] = 999
		if ma[0] == 999 {
			t.Error("Minus returned a slice aliased to input a")
		}
	}

	// SymmetricDifference — mutating output must not touch either input
	sa := []int{1, 2, 3}
	sb := []int{3, 4, 5}
	sOut := list.SymmetricDifference(sa, sb)
	if len(sOut) > 0 {
		sOut[0] = 999
		if sa[0] == 999 || sb[0] == 999 {
			t.Error("SymmetricDifference returned a slice aliased to an input")
		}
	}
}

// TestSymmetricDifference_OrderingCrossDuplicates pins the emission order
// when both slices contain duplicates that cross the a/b boundary. The
// contract is: emit a's unique-to-a elements in a's first-occurrence
// order, then b's unique-to-b elements in b's first-occurrence order.
func TestSymmetricDifference_OrderingCrossDuplicates(t *testing.T) {
	got := list.SymmetricDifference([]int{1, 2, 3, 2}, []int{3, 2, 4, 4, 5})
	// From a: 1 (unique to a), 2 and 3 are in b, dup 2 skipped
	// From b: 3 and 2 are in a, 4 unique to b, dup 4 skipped, 5 unique to b
	want := []int{1, 4, 5}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestIntersect_NilMiddleSlice makes the nil-safety contract explicit:
// a nil slice in the middle of the variadic list must be treated as
// empty, yielding an empty intersection (not a panic, not a stale result).
func TestIntersect_NilMiddleSlice(t *testing.T) {
	got := list.Intersect(
		[]int{1, 2, 3},
		nil,
		[]int{1, 2, 3},
	)
	if got == nil {
		t.Error("expected non-nil empty slice")
	}
	if len(got) != 0 {
		t.Errorf("got %v, expected empty (nil middle slice)", got)
	}
}

// TestNaNSemantics documents and verifies how floating-point NaN values
// behave in list operations. NaN follows Go's map-key semantics: NaN !=
// NaN, so each NaN is a distinct map key. This means NaNs cannot be
// deduplicated, intersected, or excluded by any list function.
//
// This is a deliberate documented limitation, not a bug. If a user needs
// NaN-aware set operations, they must pre-process the slice.
func TestNaNSemantics(t *testing.T) {
	nan := math.NaN()

	// Unique cannot dedupe NaNs because NaN != NaN in map keys
	uniq := list.Unique([]float64{nan, nan, nan})
	if len(uniq) != 3 {
		t.Errorf("Unique with NaNs: got len %d, expected 3 (NaN semantics)", len(uniq))
	}

	// Intersect cannot match NaNs across slices
	inter := list.Intersect([]float64{nan, 1.0}, []float64{nan, 1.0})
	if len(inter) != 1 || inter[0] != 1.0 {
		t.Errorf("Intersect with NaNs: got %v, expected [1.0] (NaN can't intersect with itself)", inter)
	}

	// Minus cannot remove NaN from a slice because NaN != NaN
	minus := list.Minus([]float64{nan, 1.0, 2.0}, []float64{nan, 1.0})
	// NaN remains in result because it can't be matched; 1.0 is removed; 2.0 remains
	if len(minus) != 2 {
		t.Errorf("Minus with NaN in b: got %v, expected 2 elements", minus)
	}

	// Without cannot remove NaN from a slice
	without := list.Without([]float64{nan, 1.0, 2.0}, nan)
	// NaN remains because it can't be matched; 1.0 and 2.0 remain
	if len(without) != 3 {
		t.Errorf("Without NaN: got %v, expected all 3 elements (NaN unmatchable)", without)
	}
}

// TestIntersect_CountingAlgorithm is a white-box test targeting the
// internal counting logic. The counts[v] == i guard ensures an element
// is only counted if it appeared in EVERY prior slice. A bug here would
// let elements missing from some slices slip through.
func TestIntersect_CountingAlgorithm(t *testing.T) {
	// Element "2" appears in slices 0 and 2 but not 1 — should NOT be in result
	got := list.Intersect(
		[]int{1, 2, 3},
		[]int{1, 3},
		[]int{1, 2, 3},
	)
	want := []int{1, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v (element 2 missing from slice 1, must not appear)", got, want)
	}

	// Element in every slice but at different positions
	got2 := list.Intersect(
		[]int{1, 2, 3, 4, 5},
		[]int{5, 4, 3, 2, 1},
		[]int{3, 1, 5, 2, 4},
	)
	// Order from first slice, all 5 elements present in all
	want2 := []int{1, 2, 3, 4, 5}
	if !reflect.DeepEqual(got2, want2) {
		t.Errorf("got %v, want %v", got2, want2)
	}

	// Empty intermediate slice
	got3 := list.Intersect(
		[]int{1, 2, 3},
		[]int{},
		[]int{1, 2, 3},
	)
	if len(got3) != 0 {
		t.Errorf("got %v, expected empty (middle slice is empty)", got3)
	}
}

// TestUnion_OrderAcrossDuplicates verifies that Union preserves
// first-occurrence order across slice boundaries. If slice A has element
// X and slice B also has X, X should appear in its slice-A position.
func TestUnion_OrderAcrossDuplicates(t *testing.T) {
	got := list.Union(
		[]int{3, 1, 4},
		[]int{1, 5, 9, 2, 6, 5, 3},
	)
	// Order: 3,1,4 from first; then 5,9,2,6 from second (1 and 3 already seen)
	want := []int{3, 1, 4, 5, 9, 2, 6}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestCustomComparableTypes verifies that named types based on
// comparable primitives work correctly with all functions. This is the
// common case for domain models (type UserID string, type Port uint16).
func TestCustomComparableTypes(t *testing.T) {
	type UserID string
	ids := []UserID{"alice", "bob", "alice", "carol"}
	got := list.Unique(ids)
	want := []UserID{"alice", "bob", "carol"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// TestStructKeys verifies that struct types with comparable fields
// work correctly. Go makes structs comparable if all their fields are
// comparable, so this is a legitimate use case.
func TestStructKeys(t *testing.T) {
	type Point struct {
		X, Y int
	}
	points := []Point{{1, 2}, {3, 4}, {1, 2}, {5, 6}}
	got := list.Unique(points)
	want := []Point{{1, 2}, {3, 4}, {5, 6}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}

	// Verify Union/Intersect/Minus/Without also work with structs
	a := []Point{{1, 1}, {2, 2}, {3, 3}}
	b := []Point{{2, 2}, {3, 3}, {4, 4}}

	gotUnion := list.Union(a, b)
	if len(gotUnion) != 4 {
		t.Errorf("Union: got %v, expected 4 unique points", gotUnion)
	}

	gotIntersect := list.Intersect(a, b)
	if len(gotIntersect) != 2 {
		t.Errorf("Intersect: got %v, expected 2 common points", gotIntersect)
	}

	gotMinus := list.Minus(a, b)
	if !reflect.DeepEqual(gotMinus, []Point{{1, 1}}) {
		t.Errorf("Minus: got %v, expected [{1 1}]", gotMinus)
	}
}

// TestMinus_DuplicatesInB verifies that duplicates in the exclusion slice
// are handled correctly (the exclude set collapses them, so only the first
// occurrence matters for correctness).
func TestMinus_DuplicatesInB(t *testing.T) {
	got := list.Minus([]int{1, 2, 3, 4, 5}, []int{2, 2, 2, 4, 4})
	want := []int{1, 3, 5}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// =============================================================================
// Composition — realistic usage patterns
// =============================================================================

func TestComposition_ActiveAdmins(t *testing.T) {
	admins := []int{1, 2, 3, 4}
	editors := []int{3, 4, 5, 6}
	banned := []int{2, 5}

	// Active admins or editors, minus the banned ones
	active := list.Minus(list.Union(admins, editors), banned)
	want := []int{1, 3, 4, 6}
	if !reflect.DeepEqual(active, want) {
		t.Errorf("got %v, want %v", active, want)
	}
}

func TestComposition_CommonTags(t *testing.T) {
	post1Tags := []string{"go", "web", "api", "backend"}
	post2Tags := []string{"web", "api", "frontend"}
	post3Tags := []string{"api", "web", "testing"}

	common := list.Intersect(post1Tags, post2Tags, post3Tags)
	want := []string{"web", "api"}
	if !reflect.DeepEqual(common, want) {
		t.Errorf("got %v, want %v", common, want)
	}
}

// =============================================================================
// Pointer identity semantics
// =============================================================================

// TestPointerIdentity documents that pointer types compare by address, not by
// pointed-at value. Two *int pointing at distinct ints with the same value
// are distinct keys.
func TestPointerIdentity(t *testing.T) {
	a, b := new(int), new(int)
	*a, *b = 5, 5
	got := list.Unique([]*int{a, b, a})
	if len(got) != 2 {
		t.Errorf("expected 2 distinct pointers (identity, not value), got %d", len(got))
	}
}

// =============================================================================
// Runtime panic on non-comparable interface values
// =============================================================================

// TestNonComparableInterfacePanics verifies the package-doc claim that
// interface slices carrying non-comparable dynamic values (e.g. slices)
// panic at runtime rather than silently succeeding. If this test stops
// panicking, the package docs are lying.
func TestNonComparableInterfacePanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected runtime panic comparing non-comparable interface values")
		}
	}()
	// []int is not comparable; wrapping it in any defers the check to runtime.
	_ = list.Unique([]any{[]int{1}, []int{1}})
}

// =============================================================================
// SymmetricDifference
// =============================================================================

func TestSymmetricDifference_Basic(t *testing.T) {
	got := list.SymmetricDifference([]int{1, 2, 3, 4}, []int{3, 4, 5, 6})
	want := []int{1, 2, 5, 6}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSymmetricDifference_Disjoint(t *testing.T) {
	got := list.SymmetricDifference([]int{1, 2}, []int{3, 4})
	want := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSymmetricDifference_Identical(t *testing.T) {
	got := list.SymmetricDifference([]int{1, 2, 3}, []int{1, 2, 3})
	want := []int{}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSymmetricDifference_EmptyA(t *testing.T) {
	got := list.SymmetricDifference([]int{}, []int{1, 2, 2, 3})
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSymmetricDifference_EmptyB(t *testing.T) {
	got := list.SymmetricDifference([]int{1, 2, 2, 3}, []int{})
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSymmetricDifference_NilInputs(t *testing.T) {
	got := list.SymmetricDifference[int](nil, nil)
	if got == nil {
		t.Error("expected non-nil empty slice")
	}
	if len(got) != 0 {
		t.Errorf("expected empty result, got %v", got)
	}
}

func TestSymmetricDifference_DeduplicatesWithinInputs(t *testing.T) {
	got := list.SymmetricDifference([]int{1, 1, 2, 2}, []int{2, 3, 3})
	want := []int{1, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSymmetricDifference_Immutability(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{3, 4, 5}
	aOrig := append([]int(nil), a...)
	bOrig := append([]int(nil), b...)
	_ = list.SymmetricDifference(a, b)
	if !reflect.DeepEqual(a, aOrig) {
		t.Errorf("input a mutated: got %v, want %v", a, aOrig)
	}
	if !reflect.DeepEqual(b, bOrig) {
		t.Errorf("input b mutated: got %v, want %v", b, bOrig)
	}
}

// =============================================================================
// Fuzz
// =============================================================================

// FuzzUnique validates Unique against a simple reference implementation
// and the documented invariants: length non-increasing, every result
// element comes from the input, no duplicates in the result, and the
// order matches first-occurrence order in the input.
func FuzzUnique(f *testing.F) {
	f.Add([]byte{1, 2, 2, 3, 1, 4})
	f.Add([]byte{})
	f.Add([]byte{5, 5, 5, 5, 5})
	f.Add([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

	f.Fuzz(func(t *testing.T, data []byte) {
		in := make([]int, len(data))
		for i, b := range data {
			in[i] = int(b)
		}

		got := list.Unique(in)

		if got == nil {
			t.Fatal("Unique returned nil; contract is non-nil empty")
		}

		// Reference: preserve first-occurrence order, drop later repeats.
		seen := make(map[int]struct{}, len(in))
		want := make([]int, 0, len(in))
		for _, v := range in {
			if _, ok := seen[v]; ok {
				continue
			}
			seen[v] = struct{}{}
			want = append(want, v)
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Unique(%v) = %v, reference = %v", in, got, want)
		}

		// Invariant: length never exceeds input.
		if len(got) > len(in) {
			t.Fatalf("Unique grew the slice: len(got)=%d len(in)=%d", len(got), len(in))
		}

		// Invariant: result has no duplicates.
		resSeen := make(map[int]struct{}, len(got))
		for _, v := range got {
			if _, dup := resSeen[v]; dup {
				t.Fatalf("Unique returned duplicate %v: %v", v, got)
			}
			resSeen[v] = struct{}{}
		}

		// Invariant: every result element appears in the input.
		inSet := make(map[int]struct{}, len(in))
		for _, v := range in {
			inSet[v] = struct{}{}
		}
		for _, v := range got {
			if _, ok := inSet[v]; !ok {
				t.Fatalf("Unique fabricated element %v not in input %v", v, in)
			}
		}
	})
}

func TestComposition_UniqueExcludingStopwords(t *testing.T) {
	words := []string{"the", "quick", "brown", "fox", "the", "lazy", "dog", "the"}
	stopwords := []string{"the", "a", "an"}

	significant := list.Minus(list.Unique(words), stopwords)
	want := []string{"quick", "brown", "fox", "lazy", "dog"}
	if !reflect.DeepEqual(significant, want) {
		t.Errorf("got %v, want %v", significant, want)
	}
}

// TestIntersect_FirstSliceEmpty covers the edge of the counting loop: when
// slice 0 is empty or nil, counts is never populated, so no element can
// reach the "present in every slice" threshold. The result must be an
// empty (non-nil) slice, not a panic and not stray entries from later
// slices.
func TestIntersect_FirstSliceEmpty(t *testing.T) {
	t.Run("empty first slice", func(t *testing.T) {
		got := list.Intersect([]int{}, []int{1, 2, 3}, []int{1, 2, 3})
		if got == nil {
			t.Fatal("expected non-nil empty slice")
		}
		if len(got) != 0 {
			t.Errorf("got %v, want empty", got)
		}
	})
	t.Run("nil first slice", func(t *testing.T) {
		got := list.Intersect[int](nil, []int{1, 2, 3}, []int{1, 2, 3})
		if got == nil {
			t.Fatal("expected non-nil empty slice")
		}
		if len(got) != 0 {
			t.Errorf("got %v, want empty", got)
		}
	})
	t.Run("single empty slice", func(t *testing.T) {
		// One-slice call is equivalent to Unique: empty in → empty out.
		got := list.Intersect([]int{})
		if got == nil {
			t.Fatal("expected non-nil empty slice")
		}
		if len(got) != 0 {
			t.Errorf("got %v, want empty", got)
		}
	})
}

// FuzzIntersect verifies Intersect against a reference implementation
// across arbitrary byte-derived int inputs split into two slices. Covers
// the trickier counting-and-emit algorithm that plain example tests can't
// exhaust.
func FuzzIntersect(f *testing.F) {
	f.Add([]byte{1, 2, 3, 4}, []byte{3, 4, 5, 6})
	f.Add([]byte{}, []byte{1, 2, 3})
	f.Add([]byte{1, 1, 1}, []byte{1})
	f.Add([]byte{1, 2, 3}, []byte{})
	f.Add([]byte{1, 2, 3, 4, 5}, []byte{5, 4, 3, 2, 1})

	f.Fuzz(func(t *testing.T, a, b []byte) {
		inA := make([]int, len(a))
		for i, v := range a {
			inA[i] = int(v)
		}
		inB := make([]int, len(b))
		for i, v := range b {
			inB[i] = int(v)
		}

		got := list.Intersect(inA, inB)

		if got == nil {
			t.Fatal("Intersect returned nil; contract is non-nil empty")
		}

		// Reference: elements in first-occurrence order from inA that also
		// appear in inB, deduped.
		bSet := make(map[int]struct{}, len(inB))
		for _, v := range inB {
			bSet[v] = struct{}{}
		}
		seen := make(map[int]struct{}, len(inA))
		want := make([]int, 0)
		for _, v := range inA {
			if _, ok := bSet[v]; !ok {
				continue
			}
			if _, dup := seen[v]; dup {
				continue
			}
			seen[v] = struct{}{}
			want = append(want, v)
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Intersect(%v, %v) = %v, reference = %v", inA, inB, got, want)
		}

		// Invariants: no dupes, every element present in both inputs.
		resSeen := make(map[int]struct{}, len(got))
		for _, v := range got {
			if _, dup := resSeen[v]; dup {
				t.Fatalf("Intersect returned duplicate %v: %v", v, got)
			}
			resSeen[v] = struct{}{}
			if _, ok := bSet[v]; !ok {
				t.Fatalf("Intersect fabricated element %v not in b: %v", v, got)
			}
		}
	})
}

// FuzzSymmetricDifference verifies SymmetricDifference against a reference
// implementation. Exercises the two-map emission logic (in a but not b,
// then in b but not a) with deduplication.
func FuzzSymmetricDifference(f *testing.F) {
	f.Add([]byte{1, 2, 3}, []byte{3, 4, 5})
	f.Add([]byte{}, []byte{1, 2, 3})
	f.Add([]byte{1, 2, 3}, []byte{})
	f.Add([]byte{1, 1, 2, 2}, []byte{2, 2, 3, 3})
	f.Add([]byte{1, 2, 3}, []byte{1, 2, 3})

	f.Fuzz(func(t *testing.T, a, b []byte) {
		inA := make([]int, len(a))
		for i, v := range a {
			inA[i] = int(v)
		}
		inB := make([]int, len(b))
		for i, v := range b {
			inB[i] = int(v)
		}

		got := list.SymmetricDifference(inA, inB)

		if got == nil {
			t.Fatal("SymmetricDifference returned nil; contract is non-nil empty")
		}

		// Reference: (a \ b) in a's first-occurrence order, then (b \ a)
		// in b's first-occurrence order, deduped.
		aSet := make(map[int]struct{}, len(inA))
		for _, v := range inA {
			aSet[v] = struct{}{}
		}
		bSet := make(map[int]struct{}, len(inB))
		for _, v := range inB {
			bSet[v] = struct{}{}
		}
		seen := make(map[int]struct{})
		want := make([]int, 0)
		for _, v := range inA {
			if _, inB := bSet[v]; inB {
				continue
			}
			if _, dup := seen[v]; dup {
				continue
			}
			seen[v] = struct{}{}
			want = append(want, v)
		}
		for _, v := range inB {
			if _, inA := aSet[v]; inA {
				continue
			}
			if _, dup := seen[v]; dup {
				continue
			}
			seen[v] = struct{}{}
			want = append(want, v)
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("SymmetricDifference(%v, %v) = %v, reference = %v", inA, inB, got, want)
		}

		// Invariant: no element of the result appears in BOTH inputs.
		for _, v := range got {
			_, inA := aSet[v]
			_, inB := bSet[v]
			if inA && inB {
				t.Fatalf("SymmetricDifference emitted %v which is in both inputs", v)
			}
		}
	})
}

// =============================================================================
// Positional extraction — FirstN / LastN / Between / At / First / Last
// =============================================================================

func TestFirstN(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	if got := list.FirstN(s, 3); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Errorf("FirstN(3) = %v, want [1 2 3]", got)
	}
	if got := list.FirstN(s, 0); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("FirstN(0) = %v, want []", got)
	}
	if got := list.FirstN(s, -1); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("FirstN(-1) = %v, want []", got)
	}
	if got := list.FirstN(s, 100); !reflect.DeepEqual(got, []int{1, 2, 3, 4, 5}) {
		t.Errorf("FirstN(100) = %v, want full slice", got)
	}
	if got := list.FirstN[int](nil, 3); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("FirstN(nil) = %v, want []", got)
	}
	// Result must not alias input.
	out := list.FirstN(s, 3)
	out[0] = 999
	if s[0] == 999 {
		t.Error("FirstN result aliases input")
	}
}

func TestLastN(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	if got := list.LastN(s, 3); !reflect.DeepEqual(got, []int{3, 4, 5}) {
		t.Errorf("LastN(3) = %v, want [3 4 5]", got)
	}
	if got := list.LastN(s, 0); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("LastN(0) = %v, want []", got)
	}
	if got := list.LastN(s, -5); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("LastN(-5) = %v, want []", got)
	}
	if got := list.LastN(s, 100); !reflect.DeepEqual(got, []int{1, 2, 3, 4, 5}) {
		t.Errorf("LastN(100) = %v, want full slice", got)
	}
	if got := list.LastN[int](nil, 2); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("LastN(nil) = %v, want []", got)
	}
}

func TestBetween(t *testing.T) {
	s := []int{10, 20, 30, 40, 50}
	if got := list.Between(s, 1, 4); !reflect.DeepEqual(got, []int{20, 30, 40}) {
		t.Errorf("Between(1,4) = %v, want [20 30 40]", got)
	}
	// Negative start clamps to 0.
	if got := list.Between(s, -10, 3); !reflect.DeepEqual(got, []int{10, 20, 30}) {
		t.Errorf("Between(-10,3) = %v, want [10 20 30]", got)
	}
	// End beyond length clamps to len.
	if got := list.Between(s, 2, 100); !reflect.DeepEqual(got, []int{30, 40, 50}) {
		t.Errorf("Between(2,100) = %v, want [30 40 50]", got)
	}
	// Inverted range → empty.
	if got := list.Between(s, 4, 2); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("Between(4,2) = %v, want []", got)
	}
	// start == end → empty.
	if got := list.Between(s, 2, 2); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("Between(2,2) = %v, want []", got)
	}
	// Empty input.
	if got := list.Between([]int{}, 0, 3); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("Between(empty) = %v, want []", got)
	}
	// Result must not alias input.
	out := list.Between(s, 1, 4)
	out[0] = 999
	if s[1] == 999 {
		t.Error("Between result aliases input")
	}
}

func TestAt(t *testing.T) {
	s := []string{"a", "b", "c", "d"}

	v, ok := list.At(s, 0)
	if !ok || v != "a" {
		t.Errorf("At(0) = (%q, %v), want (a, true)", v, ok)
	}
	v, ok = list.At(s, 3)
	if !ok || v != "d" {
		t.Errorf("At(3) = (%q, %v), want (d, true)", v, ok)
	}
	// Negative index counts from the end.
	v, ok = list.At(s, -1)
	if !ok || v != "d" {
		t.Errorf("At(-1) = (%q, %v), want (d, true)", v, ok)
	}
	v, ok = list.At(s, -4)
	if !ok || v != "a" {
		t.Errorf("At(-4) = (%q, %v), want (a, true)", v, ok)
	}
	// Out of range → (zero, false).
	if _, ok := list.At(s, 4); ok {
		t.Error("At(4) on len 4 slice should be out of range")
	}
	if _, ok := list.At(s, -5); ok {
		t.Error("At(-5) on len 4 slice should be out of range")
	}
	// Empty slice.
	if _, ok := list.At([]string{}, 0); ok {
		t.Error("At(empty, 0) should be out of range")
	}
	// Nil slice.
	if _, ok := list.At[string](nil, 0); ok {
		t.Error("At(nil, 0) should be out of range")
	}
}

func TestFirstAndLast(t *testing.T) {
	s := []int{10, 20, 30}

	if v, ok := list.First(s); !ok || v != 10 {
		t.Errorf("First = (%v, %v), want (10, true)", v, ok)
	}
	if v, ok := list.Last(s); !ok || v != 30 {
		t.Errorf("Last = (%v, %v), want (30, true)", v, ok)
	}

	// Empty inputs.
	if _, ok := list.First([]int{}); ok {
		t.Error("First(empty) should be false")
	}
	if _, ok := list.Last([]int{}); ok {
		t.Error("Last(empty) should be false")
	}
}

// =============================================================================
// Sampling — Sample / SampleN
// =============================================================================

func TestSample(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	// Sample should always return an element that exists in s.
	for range 50 {
		v, ok := list.Sample(s)
		if !ok {
			t.Fatal("Sample returned ok=false on non-empty slice")
		}
		found := false
		for _, x := range s {
			if x == v {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Sample returned %v, which is not in the input", v)
		}
	}

	// Empty slice.
	if _, ok := list.Sample([]int{}); ok {
		t.Error("Sample(empty) should be false")
	}
	if _, ok := list.Sample[int](nil); ok {
		t.Error("Sample(nil) should be false")
	}
}

func TestSampleN(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}

	// n=3 returns exactly 3 distinct elements all drawn from s.
	got := list.SampleN(s, 3)
	if len(got) != 3 {
		t.Errorf("SampleN(3) len = %d, want 3", len(got))
	}
	seen := make(map[int]bool)
	for _, v := range got {
		if seen[v] {
			t.Errorf("SampleN(3) returned duplicate %v", v)
		}
		seen[v] = true
		found := false
		for _, x := range s {
			if x == v {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("SampleN(3) returned %v not in input", v)
		}
	}

	// n=0 / n<0 / empty / nil.
	if got := list.SampleN(s, 0); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("SampleN(0) = %v, want []", got)
	}
	if got := list.SampleN(s, -1); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("SampleN(-1) = %v, want []", got)
	}
	if got := list.SampleN([]int{}, 3); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("SampleN(empty) = %v, want []", got)
	}

	// n >= len: returns a shuffled copy, same length as input, same multiset.
	got = list.SampleN(s, 10)
	if len(got) != len(s) {
		t.Errorf("SampleN(n >= len) len = %d, want %d", len(got), len(s))
	}
	gotSorted := list.Sort(got)
	inSorted := list.Sort(s)
	if !reflect.DeepEqual(gotSorted, inSorted) {
		t.Errorf("SampleN(n >= len) multiset = %v, want %v", gotSorted, inSorted)
	}
}

// =============================================================================
// Reordering — Reverse / Shuffle
// =============================================================================

func TestReverse(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	if got := list.Reverse(s); !reflect.DeepEqual(got, []int{5, 4, 3, 2, 1}) {
		t.Errorf("Reverse = %v, want [5 4 3 2 1]", got)
	}
	// Single element.
	if got := list.Reverse([]int{42}); !reflect.DeepEqual(got, []int{42}) {
		t.Errorf("Reverse(single) = %v, want [42]", got)
	}
	// Empty.
	if got := list.Reverse([]int{}); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("Reverse(empty) = %v, want []", got)
	}
	if got := list.Reverse[int](nil); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("Reverse(nil) = %v, want []", got)
	}
	// Input not mutated.
	in := []int{1, 2, 3}
	_ = list.Reverse(in)
	if !reflect.DeepEqual(in, []int{1, 2, 3}) {
		t.Error("Reverse mutated input")
	}
}

func TestShuffle_PreservesMultiset(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	got := list.Shuffle(s)
	if len(got) != len(s) {
		t.Errorf("Shuffle len = %d, want %d", len(got), len(s))
	}
	// Same multiset.
	gotSorted := list.Sort(got)
	inSorted := list.Sort(s)
	if !reflect.DeepEqual(gotSorted, inSorted) {
		t.Errorf("Shuffle multiset = %v, want %v", gotSorted, inSorted)
	}
	// Input not mutated.
	if !reflect.DeepEqual(s, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
		t.Error("Shuffle mutated input")
	}
	// Empty / nil.
	if got := list.Shuffle([]int{}); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("Shuffle(empty) = %v, want []", got)
	}
	if got := list.Shuffle[int](nil); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("Shuffle(nil) = %v, want []", got)
	}
}

// =============================================================================
// Sorting — Sort / SortDesc / SortBy
// =============================================================================

func TestSort(t *testing.T) {
	// Classic numeric-vs-lexicographic check: correct order is
	// 1, 2, 10, 20, 100, not "1", "10", "100", "2", "20".
	in := []int{100, 1, 20, 2, 10}
	want := []int{1, 2, 10, 20, 100}
	if got := list.Sort(in); !reflect.DeepEqual(got, want) {
		t.Errorf("Sort(ints) = %v, want %v", got, want)
	}

	// Strings.
	if got := list.Sort([]string{"c", "a", "b"}); !reflect.DeepEqual(got, []string{"a", "b", "c"}) {
		t.Errorf("Sort(strings) = %v, want [a b c]", got)
	}

	// Floats with NaN: NaN sorts before everything per cmp.Compare.
	floats := []float64{3.0, math.NaN(), 1.0, 2.0}
	got := list.Sort(floats)
	if !math.IsNaN(got[0]) {
		t.Errorf("Sort(floats) first element = %v, want NaN", got[0])
	}
	if !reflect.DeepEqual(got[1:], []float64{1.0, 2.0, 3.0}) {
		t.Errorf("Sort(floats) tail = %v, want [1 2 3]", got[1:])
	}

	// Empty / nil.
	if got := list.Sort([]int{}); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("Sort(empty) = %v, want []", got)
	}
	if got := list.Sort[int](nil); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("Sort(nil) = %v, want []", got)
	}

	// Input not mutated.
	in2 := []int{3, 1, 2}
	_ = list.Sort(in2)
	if !reflect.DeepEqual(in2, []int{3, 1, 2}) {
		t.Error("Sort mutated input")
	}
}

func TestSortDesc(t *testing.T) {
	if got := list.SortDesc([]int{1, 2, 10, 20, 100}); !reflect.DeepEqual(got, []int{100, 20, 10, 2, 1}) {
		t.Errorf("SortDesc = %v, want [100 20 10 2 1]", got)
	}
	if got := list.SortDesc([]string{"a", "c", "b"}); !reflect.DeepEqual(got, []string{"c", "b", "a"}) {
		t.Errorf("SortDesc(strings) = %v, want [c b a]", got)
	}
	if got := list.SortDesc([]int{}); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("SortDesc(empty) = %v, want []", got)
	}
}

func TestSortBy(t *testing.T) {
	type person struct {
		name string
		age  int
	}
	people := []person{
		{"Alice", 30},
		{"Bob", 25},
		{"Carol", 35},
	}
	byAge := list.SortBy(people, func(a, b person) int { return a.age - b.age })
	wantAges := []int{25, 30, 35}
	for i, p := range byAge {
		if p.age != wantAges[i] {
			t.Errorf("SortBy[%d].age = %d, want %d", i, p.age, wantAges[i])
		}
	}

	// Input not mutated.
	if people[0].name != "Alice" || people[1].name != "Bob" || people[2].name != "Carol" {
		t.Error("SortBy mutated input")
	}

	// Empty.
	if got := list.SortBy([]int{}, func(a, b int) int { return a - b }); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("SortBy(empty) = %v, want []", got)
	}
}

// =============================================================================
// Compact — strip zero values
// =============================================================================

func TestCompact(t *testing.T) {
	// Strings: "" is the zero value.
	if got := list.Compact([]string{"a", "", "b", "", "c"}); !reflect.DeepEqual(got, []string{"a", "b", "c"}) {
		t.Errorf("Compact(strings) = %v, want [a b c]", got)
	}
	// Ints: 0 is the zero value.
	if got := list.Compact([]int{1, 0, 2, 0, 3}); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Errorf("Compact(ints) = %v, want [1 2 3]", got)
	}
	// All zeros.
	if got := list.Compact([]int{0, 0, 0}); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("Compact(all zero) = %v, want []", got)
	}
	// No zeros.
	if got := list.Compact([]int{1, 2, 3}); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Errorf("Compact(no zeros) = %v, want [1 2 3]", got)
	}
	// Empty / nil.
	if got := list.Compact([]int{}); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("Compact(empty) = %v, want []", got)
	}
	if got := list.Compact[int](nil); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("Compact(nil) = %v, want []", got)
	}
}

// =============================================================================
// Substitution — Replace / ReplaceFirst
// =============================================================================

func TestReplace(t *testing.T) {
	if got := list.Replace([]int{1, 2, 3, 2, 1}, 2, 99); !reflect.DeepEqual(got, []int{1, 99, 3, 99, 1}) {
		t.Errorf("Replace = %v, want [1 99 3 99 1]", got)
	}
	// No matches.
	if got := list.Replace([]int{1, 2, 3}, 99, 0); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Errorf("Replace(no match) = %v, want [1 2 3]", got)
	}
	// Empty / nil.
	if got := list.Replace([]int{}, 1, 2); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("Replace(empty) = %v, want []", got)
	}

	// Input not mutated.
	in := []int{1, 2, 3}
	_ = list.Replace(in, 2, 99)
	if !reflect.DeepEqual(in, []int{1, 2, 3}) {
		t.Error("Replace mutated input")
	}
}

func TestReplaceFirst(t *testing.T) {
	// Only the first matching element is replaced.
	if got := list.ReplaceFirst([]int{1, 2, 3, 2, 1}, 2, 99); !reflect.DeepEqual(got, []int{1, 99, 3, 2, 1}) {
		t.Errorf("ReplaceFirst = %v, want [1 99 3 2 1]", got)
	}
	// No matches → unchanged copy.
	if got := list.ReplaceFirst([]int{1, 2, 3}, 99, 0); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Errorf("ReplaceFirst(no match) = %v, want [1 2 3]", got)
	}
	// Empty.
	if got := list.ReplaceFirst([]int{}, 1, 2); !reflect.DeepEqual(got, []int{}) {
		t.Errorf("ReplaceFirst(empty) = %v, want []", got)
	}
}
