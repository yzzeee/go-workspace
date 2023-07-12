package k

import (
	"fmt"
	"go-practice/z"
)

var A = map[int]int{
	1: 1,
	2: 2,
}

func init() {
	fmt.Println("k ->", A, z.Global)
}
