## 🧠 Q1：无缓冲 channel 和有缓冲 channel 区别？

一句话：

无缓冲 channel 强制同步，有缓冲 channel 是队列，满/空才阻塞。

## 🧠 Q2：close channel 后会发生什么？

✔ 发送：panic
✔ 接收：返回零值 + ok=false
```
v, ok := <-ch
```

## 🧠 Q3：为什么通常只让发送方 close？

因为只有发送方知道不会再发送数据了。

接收方 close 可能导致已在发送中的 goroutine panic。

## 🧠 Q4：nil channel 会怎样？

👉 永远阻塞
```
var ch chan int
<-ch // 永远阻塞
```

## 🧠 Q5：channel 会导致 goroutine 泄漏吗？

会，比如：

```
func main() {
ch := make(chan int)
go func() {
ch <- 1 // 永远没人接收
}()
}
```
这个 goroutine 会一直卡住

## 🧠 Q6：select 如何避免 goroutine 泄漏？
可以用context来控制
```
select {
case ch <- v:
case <-ctx.Done():
return
}
```

## Q7：Mutex 和 Channel 的区别？

标准答：

| Mutex | Channel      |
| ----- | ------------ |
| 数据共享  | 数据传递         |
| 强调互斥  | 强调通信         |
| 只同步访问 | 可同步 + 阻塞变更流程 |

总结一句话：

如果是共享内存就用 Mutex，如果是通信就用 channel。


## 🧠 Q8：channel 底层是怎么实现的？

关键词必讲：

hchan 结构

环形队列

sendq / recvq 等待队列

goroutine park / wake

mutex 保护

“阻塞时 goroutine 会被 park 掉，由 scheduler 唤醒”

channel_internal里面就是底层原理


