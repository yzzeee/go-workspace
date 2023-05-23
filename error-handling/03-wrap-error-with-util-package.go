package main

import (
	"fmt"
	"github.com/pkg/errors"
)

var (
	innerError3 = &innerErr3{Msg: "innerMsg3"}
	myError3    = &myErr3{Msg: "myMsg3"}
)

// Custom errors
type innerErr3 struct {
	Msg string
}

func (i *innerErr3) Error() string {
	return "innerError3"
}

type myErr3 struct {
	Msg string
}

func (i *myErr3) Error() string {
	return "myError3"
}

// Functions
func innerFunc3() error {
	return innerError3
}

func middleFunc3() error {
	// github.com/pkg/errors.Wrap() 함수를 통해서 Wrapping 을 수행한다. Wrapping 을 위해서 fmt Package 를 이용할 필요가 없다.
	if err := innerFunc3(); err != nil {
		return errors.Wrap(err, "middleError3")
	}
	return nil
}

func outerFunc3() error {
	if err := middleFunc3(); err != nil {
		return errors.Wrap(err, "outerError3")
	}
	return nil
}

func doWrapErrorWithUtilPackage() {
	// Get a wrapped error
	outerErr3 := outerFunc3()

	// Cause
	// github.com/pkg/errors Package 에서는 Wrapping 된 Error 를 하나씩 Unwrapping 하는 함수를 제공하지 않는다.
	// 대신에 가장 내부에 존재하는 Error를 반환하는 Cause() 함수를 제공한다.
	// 참고로 github.com/pkg/errors.Unwrap() 함수가 존재하는데,
	// github.com/pkg/errors.Wrap() 함수를 통해서 Wrapping 한 Error 가 아니라 fmt Package 를 활용하여 Wrapping 한 Error 를 Unwrapping 하는 함수이다.
	fmt.Printf("\n--- Cause ---\n")
	fmt.Printf("cause: %v\n", errors.Cause(outerErr3))

	// Stack
	fmt.Printf("\n--- Stack ---\n")
	fmt.Printf("%+v\n", outerErr3)

	// Is (Compare)
	fmt.Printf("\n--- Is ---\n")
	if errors.Is(outerErr3, innerError3) {
		fmt.Printf("innerError3 true\n") // Print
	} else {
		fmt.Printf("innerError3 false\n")
	}
	if errors.Is(outerErr3, myError3) {
		fmt.Printf("myError3 true\n")
	} else {
		fmt.Printf("myError3 false\n") // Print
	}

	// As (Assertion, Type Casting)
	fmt.Printf("\n--- As ---\n")
	var iErr3 *innerErr3
	if errors.As(outerErr3, &iErr3) {
		fmt.Printf("innerError3 true: %v\n", iErr3.Msg) // Print
	} else {
		fmt.Printf("innerError3 false\n")
	}

	// fmt.Printf() 함수와 함께 %+v 문법을 이용하여 Wrapping 된 Error 를 출력하면 Stack Trace 도 같이 출력된다.
	var mErr3 *myErr3
	if errors.As(outerErr3, &mErr3) {
		fmt.Printf("myError3 true: %v\n", mErr3.Msg)
	} else {
		fmt.Printf("myError3 false\n") // Print
	}
}
