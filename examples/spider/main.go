package main

import (
	"context"
	"fmt"

	"github.com/arschlochnop/cl-sdk-go"
)

// MySpider 自定义爬虫
//
// 艹！嵌入BaseSpider，自动获得统计、日志、保存等功能
type MySpider struct {
	*crawlab.BaseSpider
}

// Run 实现Spider接口
//
// 艹！所有爬取逻辑写在这里
func (s *MySpider) Run(ctx context.Context) error {
	s.LogInfo("开始爬取数据")

	// 模拟爬取100条数据
	for i := 1; i <= 100; i++ {
		// 检查是否取消
		select {
		case <-ctx.Done():
			s.LogWarn("任务被取消")
			return ctx.Err()
		default:
		}

		// 爬取数据
		data := map[string]interface{}{
			"id":     i,
			"title":  fmt.Sprintf("Item %d", i),
			"url":    fmt.Sprintf("https://example.com/item/%d", i),
			"status": "active",
		}

		// 保存数据（自动更新统计）
		if err := s.Save(data); err != nil {
			s.LogError("保存失败: %v", err)
			return err
		}

		// 每10条输出一次进度
		if i%10 == 0 {
			s.LogInfo("已爬取 %d 条数据", i)
		}

		// 模拟请求计数
		s.IncRequests()
	}

	s.LogInfo("爬取完成")
	return nil
}

func main() {
	// 创建爬虫实例
	spider := &MySpider{
		BaseSpider: crawlab.NewSpider("MySpider"),
	}

	// 执行爬虫（自动处理panic、打印统计）
	if err := spider.Execute(spider); err != nil {
		crawlab.LogError("执行失败: %v", err)
	}
}
