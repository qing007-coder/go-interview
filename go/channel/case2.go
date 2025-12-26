package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*  题目：任务调度模型
场景：
系统有多个任务要处理，比如“生成报告”、“发送邮件”、“数据备份”。
有两个任务生产者，每个生产者不断生成随机任务。
有三个任务处理者（消费者），从任务队列中获取任务并执行（打印执行信息即可）。
任务队列通过一个 channel 来表示，带缓冲容量。
当所有任务生产完成，确保消费者处理完任务后程序优雅退出
任务列表示例：
[]string{
    "生成日报", "发送邮件", "数据备份", "清理日志", "统计分析",
}
要求：
生产者不断生成任务，并发送到 channel（任务队列）。
消费者从 channel 获取任务并打印 任务处理者 X 执行任务 Y。
生产者生成完所有任务后通知消费者结束。
使用 sync.WaitGroup 或 done channel 来同步退出。
提示：
任务生产可以随机选择任务类型，也可以设定每个生产者生成 5~10 个任务。
消费者处理任务时可以 time.Sleep 模拟耗时操作。
channel 缓冲容量可以小一点，比如 3~5，模拟任务队列的“积压”。
*/

// 这个其实就是work pool 可以用ctx来退出

var wg sync.WaitGroup
var producerWg sync.WaitGroup

func main() {
	tasks := []string{"生成日报", "发送邮件", "数据备份", "清理日志", "统计分析"}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	taskChan := make(chan string, 5)

	for i := 0; i < 2; i++ {
		wg.Add(1)
		producerWg.Add(1)
		go producer(tasks, taskChan, 10)
	}

	go closeChannel(taskChan)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go consumer(ctx, i, taskChan)
	}

	wg.Wait()
	fmt.Println("done")
}

func producer(tasks []string, taskChan chan<- string, taskNum int) {
	defer wg.Done()
	defer producerWg.Done()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < taskNum; i++ {
		n := r.Intn(len(tasks))
		taskChan <- tasks[n]
	}
}

func closeChannel(taskChan chan<- string) {
	producerWg.Wait()
	close(taskChan)
}

func consumer(ctx context.Context, id int, taskChan <-chan string) {
	defer wg.Done()
	for {
		select {
		case task, ok := <-taskChan:
			if !ok {
				return
			}
			fmt.Printf("消费者%d完成了%s\n", id, task)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // 模拟完成任务的时间
		case <-ctx.Done():
			fmt.Println(ctx.Err())
			return
		}
	}
}
