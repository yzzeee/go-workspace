package main

import (
	"fmt"
)

var (
	a = c + b // == 9
	b = f()   // == 4
	c = f()   // == 5
	d = 3     // == 5 after initialization has finished
)

func f() int {
	d++
	return d
}

func main() {
	fmt.Println("hello")
	fmt.Println(a, b, c, d)
}
