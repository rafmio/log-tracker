package main

import (
	"fmt"
)

func main() {
	slise := make([]int, 0)
	fillSlice(&slise)

	fmt.Println(slise)
}

func fillSlice(s *[]int) {
	for i := 0; i <= 10; i++ {
		*s = append(*s, i*i)
	}

	fmt.Println(s)
}
