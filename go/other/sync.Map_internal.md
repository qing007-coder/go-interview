# Go sync.Map 核心原理整理

## 1️⃣ 数据结构概览

sync.Map 内部主要由三部分组成：

- **read map**
    - 只读
    - 支持无锁读取（atomic.Value）
    - 命中概率非常高

- **dirty map**
    - 存储最新写入的数据
    - 写操作主要发生在这里
    - 操作时需要加锁

- **entry**
    - 真正存储值
    - 值通过 atomic 操作保证一致性

```go
type Map struct {
    mu     Mutex
    read   atomic.Value // readOnly {m map[any]*entry}
    dirty  map[any]*entry // dirty 只是“缓冲区”，不是“存储区”
    misses int // 表示read map是不是太久了 也就是通过这个指标来判断dirty map是否要写回 read map
}
```

- **流程**
    - 先访问 read（无锁快速路径） 命中 → 直接返回  未命中但 read.amended = true（说明 dirty 里可能有新值）
      → 进入 慢路径，需要加锁
    - 加锁后去 dirty查 m.misses++
    - 如果 read 和 dirty 都没有 就会在 dirty 里 创建 entry