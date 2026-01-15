# Crawlab Go SDK Examples

艹！5个基础示例，从简单到复杂，5分钟上手！

## 📁 示例列表

| 示例 | 说明 | 代码量 | 难度 |
|-----|------|--------|------|
| [simple](./simple/) | 纯函数式，最简单 | 5行 | ⭐ |
| [spider](./spider/) | Spider接口，带统计 | 30行 | ⭐⭐ |
| [batch](./batch/) | 批量保存，性能优化 | 40行 | ⭐⭐ |
| [http](./http/) | HTTP客户端，自动重试 | 50行 | ⭐⭐⭐ |
| [config](./config/) | 配置管理，环境变量 | 35行 | ⭐⭐ |

## 🚀 快速运行

```bash
# 进入任意示例目录
cd simple

# 运行
go run main.go
```

## 📖 学习路径

### 初学者
1. **simple** - 了解基础用法
2. **spider** - 学习Spider接口
3. **batch** - 掌握性能优化

### 进阶开发
4. **http** - 网页爬虫开发
5. **config** - 灵活配置管理

## 💡 示例说明

### 1. simple - 纯函数式

最简单的方式，直接调用SDK函数：

```go
import "github.com/arschlochnop/cl-sdk-go"

func main() {
    crawlab.SaveItem(data)
    crawlab.Log("完成")
}
```

**适用场景：** 简单脚本、快速原型

### 2. spider - Spider接口

标准开发方式，自动统计和异常处理：

```go
type MySpider struct {
    *crawlab.BaseSpider
}

func (s *MySpider) Run(ctx context.Context) error {
    s.Save(data)
    s.LogInfo("完成")
    return nil
}
```

**适用场景：** 标准项目、需要统计

### 3. batch - 批量保存

性能优化，减少IPC调用次数：

```go
items := make([]interface{}, 1000)
crawlab.SaveBatch(items)  // 1次IPC vs 1000次
```

**适用场景：** 大量数据、性能敏感

### 4. http - HTTP客户端

网页爬虫，自动重试和超时：

```go
client := crawlab.NewHTTPClient(30 * time.Second)
client.SetRetry(3, 2*time.Second)
resp, _ := client.Get(ctx, url)
```

**适用场景：** 网页爬取、API调用

### 5. config - 配置管理

环境变量配置，灵活可配置：

```go
config := crawlab.LoadConfig()
config.MaxRetries      // 默认3
config.BatchSize       // 默认100
```

**适用场景：** 需要灵活配置的项目

## 🔧 本地开发

所有示例使用相对路径引用SDK：

```go
// go.mod
replace github.com/arschlochnop/cl-sdk-go => ../../
```

## 📦 生产使用

移除replace，直接使用GitHub版本：

```go
// go.mod
require github.com/arschlochnop/cl-sdk-go v0.1.0
```

---

艹！选个示例跑跑看吧，老王保证简单易懂！
