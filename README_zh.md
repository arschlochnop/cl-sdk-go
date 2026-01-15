# Crawlab Go SDK

艹！Crawlab官方Go SDK，零依赖、纯标准库实现！

## 🚀 快速开始

### 安装

```bash
go get github.com/arschlochnop/cl-sdk-go
```

### 5行代码上手

```go
package main

import "github.com/arschlochnop/cl-sdk-go"

func main() {
    crawlab.SaveItem(map[string]interface{}{"title": "Hello"})
}
```

## 📦 核心模块

| 模块 | 文件 | 说明 | 代码量 |
|-----|------|------|--------|
| 基础IPC | sdk.go | 核心通信函数 | 191行 |
| Spider接口 | spider.go | 爬虫基类（可选） | 180行 |
| 配置管理 | config.go | 环境变量配置 | 160行 |
| 重试机制 | retry.go | 自动重试 | 170行 |
| HTTP客户端 | httpclient.go | 网络请求 | 181行 |

**总计：882行，零外部依赖！**

## 📖 API参考

### 1. 数据保存

```go
// 保存单条数据
func SaveItem(item interface{}) error

// 保存多条数据
func SaveItems(items ...interface{}) error

// 批量保存（推荐）
func SaveBatch(items []interface{}) error
```

### 2. 日志输出

```go
// 基础日志
func Log(format string, args ...interface{})

// 分级日志
func LogInfo(format string, args ...interface{})
func LogError(format string, args ...interface{})
func LogWarn(format string, args ...interface{})
func LogDebug(format string, args ...interface{})
```

### 3. 环境变量

```go
// 获取Crawlab环境变量
func GetTaskID() string      // 任务ID
func GetSpiderID() string    // 爬虫ID
func GetNodeID() string      // 节点ID
func GetParam() string       // 任务参数
func GetScheduleID() string  // 调度ID

// 工具函数
func MustGetEnv(key string) string  // 必须存在
func GetEnv(key, defaultValue string) string
func ParseParamJSON(v interface{}) error  // 解析JSON参数
```

### 4. Spider接口（可选）

```go
// Spider接口
type Spider interface {
    Run(ctx context.Context) error
}

// 创建爬虫
spider := crawlab.NewSpider("MySpider")

// 保存数据（自动统计）
spider.Save(data)
spider.SaveBatch(items)

// 输出日志
spider.LogInfo("消息")
spider.LogError("错误")

// 统计信息
spider.Stats.ItemsSaved  // 保存数
spider.Stats.Requests    // 请求数
spider.Stats.Errors      // 错误数

// 执行爬虫（自动处理异常）
spider.Execute(spider)
```

### 5. 配置管理

```go
// 加载配置
config := crawlab.LoadConfig()

// 配置字段
config.MaxRetries      // 最大重试次数（默认3）
config.RetryDelay      // 重试延迟（默认2秒）
config.RequestTimeout  // 请求超时（默认30秒）
config.MaxConcurrency  // 最大并发（默认10）
config.BatchSize       // 批量大小（默认100）

// 环境变量覆盖
// CRAWLAB_MAX_RETRIES=5
// CRAWLAB_BATCH_SIZE=200
```

### 6. 重试机制

```go
// 基础重试
err := crawlab.Retry(fn, maxRetries, delay)

// Context支持
err := crawlab.RetryWithContext(ctx, fn, maxRetries, delay)

// 指数退避
err := crawlab.RetryWithBackoff(ctx, fn, maxRetries, initialDelay, maxDelay)

// 条件重试
err := crawlab.RetryIf(ctx, fn, shouldRetry, maxRetries, delay)
```

### 7. HTTP客户端

```go
// 创建客户端
client := crawlab.NewHTTPClient(30 * time.Second)

// 设置Header
client.SetHeader("User-Agent", "Crawlab/1.0")

// 设置重试
client.SetRetry(3, 2*time.Second)

// 发送请求
resp, err := client.Get(ctx, url)
resp, err := client.Post(ctx, url, body)
resp, err := client.Put(ctx, url, body)
resp, err := client.Delete(ctx, url)
```

## 💡 使用示例

### 纯函数式

```go
func main() {
    // 保存数据
    data := map[string]interface{}{"title": "Example"}
    crawlab.SaveItem(data)

    // 输出日志
    crawlab.LogInfo("完成")
}
```

### Spider接口

```go
type MySpider struct {
    *crawlab.BaseSpider
}

func (s *MySpider) Run(ctx context.Context) error {
    s.Save(map[string]interface{}{"title": "Example"})
    s.LogInfo("完成")
    return nil
}

func main() {
    spider := &MySpider{
        BaseSpider: crawlab.NewSpider("MySpider"),
    }
    spider.Execute(spider)
}
```

### HTTP爬虫

```go
type WebSpider struct {
    *crawlab.BaseSpider
    client *crawlab.HTTPClient
}

func (s *WebSpider) Run(ctx context.Context) error {
    // 创建HTTP客户端
    s.client = crawlab.NewHTTPClient(30 * time.Second)
    s.client.SetRetry(3, 2*time.Second)

    // 发送请求
    resp, err := s.client.Get(ctx, "https://example.com")
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // 读取数据
    body, _ := io.ReadAll(resp.Body)

    // 保存
    s.Save(map[string]interface{}{
        "url":  resp.Request.URL.String(),
        "size": len(body),
    })

    return nil
}
```

## 🎯 最佳实践

### 1. 性能优化

```go
// ❌ 低效：1000次IPC调用
for i := 0; i < 1000; i++ {
    crawlab.SaveItem(data)
}

// ✅ 高效：1次IPC调用
items := make([]interface{}, 1000)
crawlab.SaveBatch(items)
```

### 2. 错误处理

```go
// ✅ 推荐：Execute自动处理异常
spider.Execute(spider)

// ❌ 不推荐：手动处理
if err := spider.Run(ctx); err != nil {
    panic(err)
}
```

### 3. 数据大小检查

```go
// SDK自动检查数据大小
// 超过5MB会输出警告
crawlab.SaveItem(largeData)  // 自动警告

// 建议：大数据使用外部存储（OSS、S3等）
```

## 📝 环境变量

### Crawlab内置

| 变量名 | 说明 |
|--------|------|
| `CRAWLAB_TASK_ID` | 任务ID |
| `CRAWLAB_SPIDER_ID` | 爬虫ID |
| `CRAWLAB_NODE_ID` | 节点ID |
| `CRAWLAB_TASK_PARAM` | 任务参数（JSON） |
| `CRAWLAB_SCHEDULE_ID` | 调度ID |

### SDK配置

| 变量名 | 类型 | 默认值 |
|--------|------|--------|
| `CRAWLAB_MAX_RETRIES` | int | 3 |
| `CRAWLAB_RETRY_DELAY` | duration | 2s |
| `CRAWLAB_REQUEST_TIMEOUT` | duration | 30s |
| `CRAWLAB_MAX_CONCURRENCY` | int | 10 |
| `CRAWLAB_BATCH_SIZE` | int | 100 |

## 🔗 相关链接

- [示例代码](../examples/)
- [快速开始](../examples/QUICKSTART.md)
- [Crawlab文档](https://docs.crawlab.cn/)

---

艹！轻量高效的SDK，老王精心打造！
