# Spider Interface Example - Spider接口示例

艹！标准化的爬虫开发方式！

## 功能

- 使用Spider接口
- 自动统计（保存数、请求数、错误数）
- 自动异常捕获
- 支持取消信号

## 运行

```bash
go run main.go
```

## 代码说明

```go
// 1. 嵌入BaseSpider
type MySpider struct {
    *crawlab.BaseSpider
}

// 2. 实现Run方法
func (s *MySpider) Run(ctx context.Context) error {
    // 爬取逻辑
    s.Save(data)      // 保存数据
    s.LogInfo("...")  // 输出日志
    s.IncRequests()   // 增加请求计数
    return nil
}

// 3. 执行爬虫
spider := &MySpider{
    BaseSpider: crawlab.NewSpider("MySpider"),
}
spider.Execute(spider)
```

## 适用场景

- 标准爬虫开发
- 需要统计信息
- 长时间运行的任务
- 团队协作项目
