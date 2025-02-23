package main

import (
	"fmt"
	"strings"
)

// 1. Ví dụ về Array
func arrayExample() {
	var arr [5]int = [5]int{1, 2, 3, 4, 5}
	fmt.Println("Array:", arr)

	arr[2] = 10
	fmt.Println("Modified Array:", arr)
}

// 2. Ví dụ về String
func stringExample() {
	str := "Hello, Golang!"
	fmt.Println("Original String:", str)
	fmt.Println("Length of string:", len(str))
	fmt.Println("Uppercase:", strings.ToUpper(str))
	fmt.Println("Contains 'Go':", strings.Contains(str, "Go"))
}

// 3. Ví dụ về Slices
func sliceExample() {
	slice := []int{1, 2, 3, 4, 5}
	fmt.Println("Slice:", slice)

	slice = append(slice, 6, 7)
	fmt.Println("Appended Slice:", slice)

	subSlice := slice[1:4]
	fmt.Println("SubSlice (1:4):", subSlice)
}

// 4. Ví dụ về Functions
func add(a int, b int) int {
	return a + b
}

// 5. Ví dụ về Methods trên Struct
type Rectangle struct {
	width, height float64
}

func (r Rectangle) Area() float64 {
	return r.width * r.height
}

// Hàm main duy nhất gọi tất cả các ví dụ trên
func main() {
	fmt.Println("\n--- Array Example ---")
	arrayExample()

	fmt.Println("\n--- String Example ---")
	stringExample()

	fmt.Println("\n--- Slice Example ---")
	sliceExample()

	fmt.Println("\n--- Function Example ---")
	result := add(5, 7)
	fmt.Println("Sum:", result)

	fmt.Println("\n--- Method Example ---")
	rect := Rectangle{width: 5, height: 10}
	fmt.Println("Rectangle Area:", rect.Area())
}
