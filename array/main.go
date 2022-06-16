package main

import (
	"fmt"
	"strconv"
	"time"
)

// MetricKey 단위 정의 키
type MetricKey string

const (
	NodeCpu    = MetricKey("node_cpu")
	NodeMemory = MetricKey("node_memory")
)

var (
	MetricDefinitions = map[MetricKey]string{
		NodeCpu:    "1111",
		NodeMemory: "2222",
	}
)

func main() {
	//fmt.Println(time.Now().Unix())
	//timestamp := time.Unix(time.Now().Unix(), 0)

	// specific date setting
	//timestamp := time.Date(2022, 6, 26, 12, 10, 30, 0, time.Local)
	//fmt.Printf("%v\n", timestamp)
	//fmt.Println(timestamp.Unix())
	//
	//timestamp2 := time.Unix(1656245430, 0)
	//fmt.Printf("%v\n", timestamp2.UTC())

	//loc, err := time.LoadLocation("Asia/Seoul")
	//if err != nil {
	//	panic(err)
	//}
	//now := time.Unix(1656288524, 0) // Go Playground 에서는 항상 시각은 2009-11-10 23:00:00 +0000 UTC 에서 시작한다.
	//t := now.In(loc)
	//fmt.Println("now=", now)
	//fmt.Println("kst=", t)

	now := time.Now()
	tt := now.Add(-time.Minute * 5)
	//fmt.Println(time.Unix(1656244439, 0), time.Now())
	fmt.Println("===>", tt.Unix(), strconv.Itoa(int(now.Unix())))
	//var queries = []string{"1: %s", "2: %s"}
	//fmt.Println(queries)
	//
	//requests := make([]interface{}, len(queries))
	//for i, query := range queries {
	//	requests[i] = query
	//}
	//
	//fmt.Println(MetricDefinitions[MetricKey("node_cpu")])
	//
	//var maxFloat64 interface{}
	//var values interface{} = "15131231"
	//if str, ok := values.(float64); ok {
	//	maxFloat64 = values.(float64)
	//	fmt.Println(ok, str)
	//} else {
	//	fmt.Println("float 64 not ok", str)
	//}
	//
	//fmt.Println(maxFloat64)

	var t time.Time

	// 2022-06-22 10:50:11
	printTime(time.Date(2022, 6, 22, 10, 50, 11, 0, time.Local))

	// now
	fmt.Println("Now")
	now2 := time.Now()
	printTime(now2)

	// +5min
	fmt.Println("+5min")
	t = time.Time.Add(now2, time.Minute*5)
	printTime(t)

	// +30min
	fmt.Println("+30min")
	t = time.Time.Add(now2, time.Minute*30)
	printTime(t)

	// +60min
	fmt.Println("+60min")
	t = time.Time.Add(now2, time.Minute*60)
	printTime(t)
}

func printTime(t time.Time) {
	fmt.Println("unix: ", t.Unix())

	// Javascript - new Date('2022-06-22 10:50:11').getTime() 와 동일한 형태
	fmt.Println("ms: ", t.UnixNano())

	// Reference time (포맷팅은 이 시간을 기준으로 한다.)
	// Mon Jan 2 15:04:05 -0700 MST 2006
	fmt.Println(t.Format("2006-01-02 15:04:05"))
	fmt.Println(t.Format("2006-01-02 PM 03:04:05"))
	fmt.Println()
}
