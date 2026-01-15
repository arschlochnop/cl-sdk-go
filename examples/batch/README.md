# Batch Save Example - 批量保存示例

艹！大量数据高效保存的正确姿势！

## 功能

- 批量保存数据（减少IPC次数）
- 分批处理大量数据
- 性能优化示例

## 运行

```bash
go run main.go
```

## 代码说明

```go
// 方式1: 一次保存100条
items := make([]interface{}, 100)
for i := 0; i < 100; i++ {
    items[i] = map[string]interface{}{"id": i}
}
crawlab.SaveBatch(items)  // 一次IPC调用

// 方式2: 分批保存1000条（每批100条）
for start := 0; start < 1000; start += 100 {
    batch := items[start:start+100]
    crawlab.SaveBatch(batch)
}
```

## 性能对比

| 方式 | 1000条数据 | IPC次数 |
|-----|-----------|---------|
| SaveItem | `for i in 1000: SaveItem()` | 1000次 |
| SaveBatch | `SaveBatch(1000 items)` | 1次 |
| 分批 | `10 batches × 100 items` | 10次 |

## 适用场景

- 大量数据爬取（>100条）
- 性能敏感的任务
- 高并发爬取
- 数据ETL任务
