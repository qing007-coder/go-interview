package main

import (
	"context"
	"fmt"
	"time"
)

func main() {

	// 控制协程

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	dataCh := make(chan string, 1)
	go func() {
		time.Sleep(time.Second * 10)
		dataCh <- "hello"
	}()

	select {
	case data := <-dataCh:
		fmt.Println(data)
	case <-ctx.Done():
		fmt.Println("time out")
	}
}
