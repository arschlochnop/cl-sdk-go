# HTTP Client Example - HTTP客户端示例

艹！封装好的HTTP客户端，自动重试超方便！

## 功能

- 自动重试（5xx错误）
- 自定义请求头
- 支持GET/POST/PUT/DELETE
- Context支持（超时、取消）
- 批量爬取示例

## 运行

```bash
go run main.go
```

## 代码说明

```go
// 创建客户端
client := crawlab.NewHTTPClient(30 * time.Second)

// 设置Header
client.SetHeader("User-Agent", "Crawlab/1.0")

// 设置重试（最多3次，延迟2秒）
client.SetRetry(3, 2*time.Second)

// 发送请求
resp, err := client.Get(ctx, url)
if err != nil {
    // 处理错误
}

// 读取响应
body, _ := io.ReadAll(resp.Body)
resp.Body.Close()
```

## 特性

- ✅ 自动重试5xx错误
- ✅ 支持超时控制
- ✅ 统一Header管理
- ✅ Context取消支持

## 适用场景

- 网页爬取
- API数据采集
- 需要重试的场景
- 批量URL爬取
