# Crawlab Go SDK

Official Go SDK for Crawlab - Zero dependencies, stdlib only!

[中文文档](README_zh.md)

## Quick Start

```bash
go get github.com/arschlochnop/cl-sdk-go
```

```go
package main

import "github.com/arschlochnop/cl-sdk-go"

func main() {
    crawlab.SaveItem(map[string]interface{}{"title": "Hello Crawlab"})
}
```

## Core Modules

- **sdk.go** (191 lines) - IPC communication
- **spider.go** (180 lines) - Spider interface (optional)
- **config.go** (160 lines) - Configuration management
- **retry.go** (170 lines) - Retry mechanisms
- **httpclient.go** (181 lines) - HTTP client

**Total: 882 lines, zero dependencies!**

## API Reference

### Data Operations

```go
crawlab.SaveItem(item)          // Save single item
crawlab.SaveItems(item1, item2) // Save multiple items
crawlab.SaveBatch(items)         // Batch save (recommended)
```

### Logging

```go
crawlab.Log("message")
crawlab.LogInfo("info message")
crawlab.LogError("error message")
crawlab.LogWarn("warning message")
```

### Environment Variables

```go
taskID := crawlab.GetTaskID()
spiderID := crawlab.GetSpiderID()
param := crawlab.GetParam()

// Parse JSON param
var params MyParams
crawlab.ParseParamJSON(&params)
```

### Spider Interface (Optional)

```go
type MySpider struct {
    *crawlab.BaseSpider
}

func (s *MySpider) Run(ctx context.Context) error {
    s.Save(data)
    s.LogInfo("done")
    return nil
}

func main() {
    spider := &MySpider{
        BaseSpider: crawlab.NewSpider("MySpider"),
    }
    spider.Execute(spider)
}
```

### HTTP Client

```go
client := crawlab.NewHTTPClient(30 * time.Second)
client.SetHeader("User-Agent", "Crawlab/1.0")
client.SetRetry(3, 2*time.Second)

resp, err := client.Get(ctx, url)
```

## Best Practices

### Performance

```go
// Bad: 1000 IPC calls
for i := 0; i < 1000; i++ {
    crawlab.SaveItem(data)
}

// Good: 1 IPC call
items := make([]interface{}, 1000)
crawlab.SaveBatch(items)
```

### Error Handling

```go
// Recommended: Auto error handling
spider.Execute(spider)

// Not recommended: Manual handling
spider.Run(ctx)
```

## Environment Variables

### Crawlab Built-in

- `CRAWLAB_TASK_ID` - Task ID
- `CRAWLAB_SPIDER_ID` - Spider ID
- `CRAWLAB_NODE_ID` - Node ID
- `CRAWLAB_TASK_PARAM` - Task parameters (JSON)
- `CRAWLAB_SCHEDULE_ID` - Schedule ID

### SDK Configuration

- `CRAWLAB_MAX_RETRIES` (default: 3)
- `CRAWLAB_RETRY_DELAY` (default: 2s)
- `CRAWLAB_REQUEST_TIMEOUT` (default: 30s)
- `CRAWLAB_MAX_CONCURRENCY` (default: 10)
- `CRAWLAB_BATCH_SIZE` (default: 100)

## Examples

See [examples directory](./examples/) for complete examples:

- **[simple](./examples/simple/)** - Pure functional style (5 lines)
- **[spider](./examples/spider/)** - Spider interface with stats
- **[batch](./examples/batch/)** - Batch operations for performance
- **[http](./examples/http/)** - HTTP client with retry
- **[config](./examples/config/)** - Configuration management

## Links

- [GitHub Repository](https://github.com/arschlochnop/cl-sdk-go)
- [Crawlab Official](https://github.com/crawlab-team/crawlab)

## License

MIT
