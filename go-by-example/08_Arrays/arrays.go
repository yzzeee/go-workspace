package main

import "fmt"

func main() {
	var a [5]int
	fmt.Println("emt: ", a) // [0 0 0 0 0]

	a[4] = 100
	fmt.Println("set: ", a) // [0 0 0 0 100]

	b := [5]int{1, 2, 3, 4, 5} // declare and initialize an array in one line
	fmt.Println("dcl: ", b)    // [1 2 3 4 5]

	c := [...]int{1, 2, 3, 4, 5}
	fmt.Println("dcl2: ", c)

	var d [2][3]int
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			d[i][j] = i + j
		}
	}

	fmt.Println("2d: ", d) // [[0 1 2][1 2 3]]
}
