package main

// https://ssup2.github.io/programming/Golang_Error_Wrapping/
func main() {
	/* Golang 에서 Error Wrapping 없이 Error 를 처리하는 일반적인 방법
	   문제는 main() 함수에서는 outerFunc() 함수가 반환하는 "outerErr"
	   Error 만 확인이 가능할 뿐, middleFunc() 또는 innerFunc() 함수가 반환하는 Error 의 내용을 확인할 수가 없다. */
	doError()

	/* Golang 1.13 이후 Version 부터 fmt.Errorf() 함수를 통해서 Error Wrapping 이 가능하며, Wrapping 된 Error 는 errors.Unwrap() 함수를 통해서 다시 얻을 수 있다.
	   또한 Wrapping 된 Error 의 비교는 errors.Is() 함수를 통해서 가능하며, Wrapping 된 Error 의 Assertion 은 errors.As() 함수를 통해서 가능하다. */
	doWrapErrorWithStandardErrorPackage()

	/* Golang 의 Standard Package 를 활용하여 Error Wrapping 을 수행할 경우 단점은 Error 가 Code 어디서 발생하는 파악이 어렵다.
	   github.com/pkg/errors Package 를 이용할 경우 Stack Trace 출력이 가능하기 때문에 Error 가 Code 어디서 발생하였는지 쉽게 파악이 가능하다.	*/
	doWrapErrorWithUtilPackage()
}
