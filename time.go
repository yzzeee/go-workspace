package main

import (
	"fmt"
	"strconv"
	"time"
)

func MsToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(0, msInt*int64(time.Millisecond)), nil
}

func main() {
	now := time.Now()
	fmt.Println("현재 시간: ", now)

	add1min := now.Add(time.Duration(1) * time.Minute)
	fmt.Println("현재 시간에 1분 더하기: ", add1min)

	add1minTimestamp := add1min.UnixNano() / int64(time.Millisecond)
	fmt.Println("현재 시간 1분을 더한 타임스탬프 값 확인: ", add1minTimestamp)

	add2minTimestamp := add1minTimestamp + (2 * time.Minute).Milliseconds()

	fmt.Println(add2minTimestamp, 1668652754892)

	readableTime1, _ := MsToTime(fmt.Sprintf("%d", add1minTimestamp))
	readableTime2, _ := MsToTime(fmt.Sprintf("%d", add2minTimestamp))
	fmt.Println("타임스탬프 값을 다시 보기 좋게 변환: ", readableTime1)
	fmt.Println("타임스탬프 값을 다시 보기 좋게 변환: ", readableTime2)

	fmt.Println("다시 1분을 빼면?!: ", readableTime1.Sub(now))
	fmt.Println("또 1분을 빼면?!: ", readableTime2.Sub(now))

	fmt.Println(time.Now().Local())

}