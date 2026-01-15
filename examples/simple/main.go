package main

import (
	"github.com/arschlochnop/cl-sdk-go"
)

// 艹！最简单的爬虫示例 - 纯函数式
// 不需要继承任何基类，直接用SDK函数即可

func main() {
	crawlab.Log("开始执行爬虫")

	// 模拟爬取10条数据
	for i := 1; i <= 10; i++ {
		data := map[string]interface{}{
			"id":    i,
			"title": "Example " + string(rune(i+'A'-1)),
			"url":   "https://example.com/" + string(rune(i+'0')),
		}

		// 保存数据
		if err := crawlab.SaveItem(data); err != nil {
			crawlab.LogError("保存数据失败: %v", err)
			continue
		}

		crawlab.LogInfo("已保存数据 #%d: %s", i, data["title"])
	}

	crawlab.Log("爬虫执行完成")
}
