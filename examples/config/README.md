# Config Example - 配置管理示例

艹！统一管理爬虫配置，环境变量一把梭！

## 功能

- 从环境变量加载配置
- 配置验证
- 默认值设置
- 配置打印

## 运行

```bash
# 使用默认配置
go run main.go

# 自定义配置
CRAWLAB_MAX_RETRIES=5 \
CRAWLAB_BATCH_SIZE=200 \
CRAWLAB_REQUEST_TIMEOUT=60s \
go run main.go
```

## 代码说明

```go
// 加载配置
config := crawlab.LoadConfig()

// 验证配置
config.Validate()

// 使用配置
spider := &ConfigSpider{
    config: config,
}

// 在爬虫中使用
batchSize := s.config.BatchSize
maxRetries := s.config.MaxRetries
```

## 环境变量

| 变量名 | 类型 | 默认值 | 说明 |
|--------|------|--------|------|
| `CRAWLAB_MAX_RETRIES` | int | 3 | 最大重试次数 |
| `CRAWLAB_RETRY_DELAY` | duration | 2s | 重试延迟 |
| `CRAWLAB_REQUEST_TIMEOUT` | duration | 30s | 请求超时 |
| `CRAWLAB_MAX_CONCURRENCY` | int | 10 | 最大并发数 |
| `CRAWLAB_BATCH_SIZE` | int | 100 | 批量大小 |

## 适用场景

- 需要灵活配置的爬虫
- 多环境部署（开发/测试/生产）
- 性能调优
- 团队协作（统一配置规范）
