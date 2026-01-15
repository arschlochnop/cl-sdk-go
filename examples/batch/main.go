package main

import (
	"fmt"

	"github.com/arschlochnop/cl-sdk-go"
)

// 艹！批量保存示例 - 减少IPC次数，性能更好

func main() {
	crawlab.Log("开始批量爬取")

	// 方式1: 使用SaveBatch批量保存（推荐）
	batchSize := 100
	items := make([]interface{}, batchSize)

	for i := 0; i < batchSize; i++ {
		items[i] = map[string]interface{}{
			"id":       i + 1,
			"title":    fmt.Sprintf("Batch Item %d", i+1),
			"url":      fmt.Sprintf("https://example.com/batch/%d", i+1),
			"category": "batch",
		}
	}

	// 一次发送100条数据（减少IPC次数）
	crawlab.LogInfo("批量保存 %d 条数据...", len(items))
	if err := crawlab.SaveBatch(items); err != nil {
		crawlab.LogError("批量保存失败: %v", err)
		return
	}

	crawlab.LogInfo("✅ 批量保存完成")

	// 方式2: 分批保存大量数据
	totalItems := 1000
	batchSize = 100

	crawlab.LogInfo("开始分批保存 %d 条数据，每批 %d 条", totalItems, batchSize)

	for start := 0; start < totalItems; start += batchSize {
		end := start + batchSize
		if end > totalItems {
			end = totalItems
		}

		// 创建当前批次数据
		batch := make([]interface{}, end-start)
		for i := 0; i < len(batch); i++ {
			batch[i] = map[string]interface{}{
				"id":    start + i + 1,
				"title": fmt.Sprintf("Item %d", start+i+1),
				"batch": start / batchSize,
			}
		}

		// 保存当前批次
		if err := crawlab.SaveBatch(batch); err != nil {
			crawlab.LogError("批次 %d 保存失败: %v", start/batchSize, err)
			continue
		}

		crawlab.LogInfo("✅ 批次 %d 完成（%d-%d/%d）", start/batchSize, start+1, end, totalItems)
	}

	crawlab.Log("所有数据保存完成")
}
