# Simple Example - 纯函数式爬虫

艹！最简单的爬虫示例，5分钟上手！

## 功能

- 纯函数式，无需继承
- 保存10条示例数据
- 基础日志输出

## 运行

```bash
go run main.go
```

## 代码说明

```go
// 保存数据
crawlab.SaveItem(data)

// 输出日志
crawlab.Log("消息")
crawlab.LogInfo("信息")
crawlab.LogError("错误")
```

## 适用场景

- 快速原型开发
- 简单的单次爬取任务
- 学习SDK基础用法
