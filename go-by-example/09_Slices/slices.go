package main

import "fmt"

func main() {

	// 생성 방법 1
	// Slice 선언은 배열을 선언하듯이 var v []T 처럼 하는데, 배열과 달리 크기 지정을 하지 않음
	var a []int
	fmt.Println("a: ", a) // []
	a = []int{1, 2, 3}
	fmt.Println("a: ", a) // [1, 2, 3]

	// 생성 방법 2
	// Go의 내장함수 make 사용
	// var v []T 형태로 선언한 것과는 다르게 모든 요소가 Zero value 인 슬라이스를 만듦
	// 세 번째 파라미터를 생략하였을 때 Capacity 와 Length 는 동일한 값을 가진다.
	b := make([]int, 3)
	fmt.Println("b: ", b, len(b), cap(b)) // [0 0 0]

}
