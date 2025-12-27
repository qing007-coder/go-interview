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

	// 传递数据
	ctx = context.WithValue(context.Background(), "id", "123456789")
	processRequest(ctx)
}

func processRequest(ctx context.Context) {
	id := ctx.Value("id")
	fmt.Println("id:", id)
}
