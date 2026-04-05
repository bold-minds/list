package examples_test

import (
	"fmt"

	"github.com/bold-minds/list"
)

func ExampleUnique() {
	fmt.Println(list.Unique([]int{1, 2, 2, 3, 1, 4}))
	// Output: [1 2 3 4]
}

func ExampleUnique_strings() {
	fmt.Println(list.Unique([]string{"go", "web", "api", "go"}))
	// Output: [go web api]
}

func ExampleUnion() {
	fmt.Println(list.Union([]int{1, 2, 3}, []int{3, 4, 5}))
	// Output: [1 2 3 4 5]
}

func ExampleUnion_manySlices() {
	fmt.Println(list.Union(
		[]int{1, 2},
		[]int{2, 3},
		[]int{3, 4},
	))
	// Output: [1 2 3 4]
}

func ExampleIntersect() {
	fmt.Println(list.Intersect([]int{1, 2, 3, 4}, []int{2, 3, 5}))
	// Output: [2 3]
}

func ExampleIntersect_manySlices() {
	fmt.Println(list.Intersect(
		[]int{1, 2, 3},
		[]int{2, 3, 4},
		[]int{3, 4, 5},
	))
	// Output: [3]
}

func ExampleMinus() {
	admins := []int{1, 2, 3, 4}
	banned := []int{2, 4}
	fmt.Println(list.Minus(admins, banned))
	// Output: [1 3]
}

func ExampleWithout() {
	fmt.Println(list.Without([]int{1, 2, 3, 4, 5}, 2, 4))
	// Output: [1 3 5]
}

func ExampleWithout_preservesDuplicates() {
	// Unlike Unique/Union/Intersect/Minus, Without does NOT deduplicate
	fmt.Println(list.Without([]int{1, 2, 3, 2, 1}, 3))
	// Output: [1 2 2 1]
}

func Example_compose() {
	// Realistic pattern: active users = admins ∪ editors, minus banned
	admins := []int{1, 2, 3}
	editors := []int{3, 4, 5}
	banned := []int{2, 5}

	active := list.Minus(list.Union(admins, editors), banned)
	fmt.Println(active)
	// Output: [1 3 4]
}
