package main

import (
	"context"
	"io"
	"time"

	"github.com/arschlochnop/cl-sdk-go"
)

// 艹！HTTP客户端示例 - 自动重试、自定义Header

func main() {
	crawlab.Log("开始HTTP爬取")

	// 创建HTTP客户端（30秒超时）
	client := crawlab.NewHTTPClient(30 * time.Second)

	// 设置User-Agent
	client.SetHeader("User-Agent", "Crawlab/1.0 (Spider Bot)")
	client.SetHeader("Accept", "text/html,application/json")

	// 设置重试参数（最多重试3次，延迟2秒）
	client.SetRetry(3, 2*time.Second)

	ctx := context.Background()

	// 示例1: 爬取单个页面
	crawlab.LogInfo("爬取示例网站...")
	resp, err := client.Get(ctx, "https://httpbin.org/get")
	if err != nil {
		crawlab.LogError("请求失败: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	data := map[string]interface{}{
		"url":         resp.Request.URL.String(),
		"status_code": resp.StatusCode,
		"content_len": len(body),
		"headers":     resp.Header,
	}

	if err := crawlab.SaveItem(data); err != nil {
		crawlab.LogError("保存失败: %v", err)
		return
	}

	crawlab.LogInfo("✅ 页面爬取完成: %d bytes", len(body))

	// 示例2: 批量爬取多个URL
	urls := []string{
		"https://httpbin.org/status/200",
		"https://httpbin.org/status/404",
		"https://httpbin.org/delay/1",
	}

	crawlab.LogInfo("开始批量爬取 %d 个URL...", len(urls))

	items := make([]interface{}, 0, len(urls))

	for i, url := range urls {
		crawlab.LogInfo("爬取 %d/%d: %s", i+1, len(urls), url)

		resp, err := client.Get(ctx, url)
		if err != nil {
			crawlab.LogWarn("请求失败: %s - %v", url, err)
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		items = append(items, map[string]interface{}{
			"url":    url,
			"status": resp.StatusCode,
			"size":   len(body),
			"time":   time.Now().Unix(),
		})

		crawlab.LogInfo("✅ 完成: %s (状态码: %d)", url, resp.StatusCode)
	}

	// 批量保存
	if len(items) > 0 {
		if err := crawlab.SaveBatch(items); err != nil {
			crawlab.LogError("批量保存失败: %v", err)
		} else {
			crawlab.LogInfo("✅ 批量保存完成: %d 条数据", len(items))
		}
	}

	crawlab.Log("HTTP爬取完成")
}
