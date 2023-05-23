package main

import (
	"errors"
	"fmt"
)

var (
	innerError2 = &innerErr2{Msg: "innerMsg2"}
	myError2    = &myErr2{Msg: "myMsg2"}
)

// Custom errors
type innerErr2 struct {
	Msg string
}

func (i *innerErr2) Error() string {
	return "innerError2"
}

type myErr2 struct {
	Msg string
}

func (i *myErr2) Error() string {
	return "myError2"
}

// Functions
func innerFunc2() error {
	// innerError2 Error 가 middleFunc2(), outerFunc2() 함수를 통해서 2번 Wrapping 되는 것을 확인할 수 있다.
	return innerError2
}

func middleFunc2() error {
	if err := innerFunc2(); err != nil {
		return fmt.Errorf("middleError2: %w", err)
	}
	return nil
}

func outerFunc2() error {
	if err := middleFunc2(); err != nil {
		// fmt.Errorf() 함수를 통해서 Error Wrapping 을 수행한다. 이 경우 반드시 %w 문법을 통해서 Error Wrapping 을 수행해야 한다.
		return fmt.Errorf("outerError2: %w", err)
	}
	return nil
}

func doWrapErrorWithStandardErrorPackage() {
	// Get a wrapped error
	outerErr2 := outerFunc2()

	// Unwrap
	// Wrapping 된 Error를 errors.Unwrap() 함수를 통해서 하나씩 Unwrapping 하며 출력한다.
	fmt.Printf("--- Unwrap ---\n")
	fmt.Printf("unwrap x 0: %v\n", outerErr2)
	fmt.Printf("unwrap x 1: %v\n", errors.Unwrap(outerErr2))
	fmt.Printf("unwrap x 2: %v\n", errors.Unwrap(errors.Unwrap(outerErr2)))

	// Is (Compare)
	// Wrapping 된 Error 를 errors.Is() 함수를 통해서 비교한다.
	fmt.Printf("\n--- Is ---\n")
	if errors.Is(outerErr2, innerError2) {
		fmt.Printf("innerError2 true\n") // Print
	} else {
		fmt.Printf("innerError2 false\n")
	}
	if errors.Is(outerErr2, myError2) {
		fmt.Printf("myError2 true\n")
	} else {
		fmt.Printf("myError2 false\n") // Print
	}

	// outerErr Error 를 errors.As() 함수를 통해서 Assertion 을 수행한다.
	// As (Assertion, Type Casting)
	fmt.Printf("\n--- As ---\n")
	var iErr2 *innerErr2
	if errors.As(outerErr2, &iErr2) {
		fmt.Printf("innerError2 true: %v\n", iErr2.Msg) // Print
	} else {
		fmt.Printf("innerError2 false\n")
	}
	var mErr2 *myErr2
	if errors.As(outerErr2, &mErr2) {
		fmt.Printf("myError2 true: %v\n", mErr2.Msg)
	} else {
		fmt.Printf("myError2 false\n") // Print
	}
}
