package main

import (
	"context"

	"github.com/arschlochnop/cl-sdk-go"
)

// ConfigSpider 使用配置的爬虫
//
// 艹！展示如何使用Config管理爬虫配置
type ConfigSpider struct {
	*crawlab.BaseSpider
	config *crawlab.Config
}

// Run 实现Spider接口
func (s *ConfigSpider) Run(ctx context.Context) error {
	// 打印配置信息
	s.config.Print()

	s.LogInfo("开始使用配置执行爬虫")

	// 使用配置中的参数
	s.LogInfo("最大重试次数: %d", s.config.MaxRetries)
	s.LogInfo("重试延迟: %v", s.config.RetryDelay)
	s.LogInfo("请求超时: %v", s.config.RequestTimeout)
	s.LogInfo("最大并发数: %d", s.config.MaxConcurrency)
	s.LogInfo("批量大小: %d", s.config.BatchSize)

	// 模拟按批次保存数据
	totalItems := 500
	items := make([]interface{}, 0, s.config.BatchSize)

	for i := 1; i <= totalItems; i++ {
		items = append(items, map[string]interface{}{
			"id":    i,
			"value": i * 100,
		})

		// 达到批量大小时保存
		if len(items) >= s.config.BatchSize {
			s.LogInfo("保存批次数据: %d 条", len(items))
			if err := s.SaveBatch(items); err != nil {
				return err
			}
			items = items[:0] // 清空
		}
	}

	// 保存剩余数据
	if len(items) > 0 {
		s.LogInfo("保存最后批次: %d 条", len(items))
		if err := s.SaveBatch(items); err != nil {
			return err
		}
	}

	s.LogInfo("配置爬虫执行完成")
	return nil
}

func main() {
	// 加载配置（从环境变量）
	config := crawlab.LoadConfig()

	// 验证配置
	if err := config.Validate(); err != nil {
		crawlab.LogError("配置验证失败: %v", err)
		return
	}

	// 创建爬虫
	spider := &ConfigSpider{
		BaseSpider: crawlab.NewSpider("ConfigSpider"),
		config:     config,
	}

	// 执行爬虫
	if err := spider.Execute(spider); err != nil {
		crawlab.LogError("执行失败: %v", err)
	}
}
