package main

import (
	"context"
	"fmt"
	"time"
)

// 两个协程交替打印 0~99
// 偶数：0 2 4 ...
// 奇数：1 3 5 ...

//var wg sync.WaitGroup

func main() {
	evenChan := make(chan struct{})
	oddChan := make(chan struct{})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	wg.Add(2)
	go even(ctx, evenChan, oddChan, 100)
	go odd(ctx, evenChan, oddChan, 100)

	// 启动
	evenChan <- struct{}{}

	wg.Wait()
	fmt.Println("done")
}

func even(ctx context.Context, evenCh, oddCh chan struct{}, num int) {
	defer wg.Done()

	for i := 0; i < num; i += 2 {
		select {
		case <-ctx.Done():
			return
		case <-evenCh:
			fmt.Println(i)
			// 最后一次不再发送
			if i+1 < num {
				oddCh <- struct{}{}
			}
		}
	}
}

func odd(ctx context.Context, evenCh, oddCh chan struct{}, num int) {
	defer wg.Done()

	for i := 1; i < num; i += 2 {
		select {
		case <-ctx.Done():
			return
		case <-oddCh:
			fmt.Println(i)
			// 最后一次不再发送
			if i+1 < num {
				evenCh <- struct{}{}
			}
		}
	}
}
