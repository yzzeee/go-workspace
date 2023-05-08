package main

import (
	"fmt"
	"time"
)

func main() {
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	// 3개의 워커가 실행
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// 5개의 작업을 보냄
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	// 작업을 다 보냈음을 알리기 위해 채널을 close
	close(jobs)

	// 모든 작업의 결과값들을 가져옴
	for a := 1; a <= 5; a++ {
		<-results
	}
}

// worker 여러개의 인스턴스를 동시에 실행할 워커
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}
