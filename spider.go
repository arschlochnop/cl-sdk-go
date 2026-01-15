package crawlab

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// Spider 爬虫接口
//
// 艹！实现这个接口就能用BaseSpider的Execute方法
type Spider interface {
	Run(ctx context.Context) error
}

// Stats 爬虫统计信息
type Stats struct {
	ItemsSaved int64     // 保存的数据条数
	Requests   int64     // 请求次数
	Errors     int64     // 错误次数
	StartTime  time.Time // 开始时间
}

// SpiderContext 爬虫上下文
type SpiderContext struct {
	TaskID     string                // 任务ID
	SpiderID   string                // 爬虫ID
	NodeID     string                // 节点ID
	Param      string                // 任务参数
	ScheduleID string                // 调度ID
	CancelFunc context.CancelFunc    // 取消函数
}

// BaseSpider 基础爬虫实现
//
// 艹！嵌入到你的Spider里，自动获得统计、日志、保存等功能
type BaseSpider struct {
	Name    string         // 爬虫名称
	Stats   *Stats         // 统计信息
	Context *SpiderContext // 爬虫上下文
	mu      sync.Mutex     // 保护Stats的并发访问
}

// NewSpider 创建一个新的BaseSpider
//
// 艹！自动获取环境变量初始化上下文
func NewSpider(name string) *BaseSpider {
	return &BaseSpider{
		Name: name,
		Stats: &Stats{
			StartTime: time.Now(),
		},
		Context: &SpiderContext{
			TaskID:     GetTaskID(),
			SpiderID:   GetSpiderID(),
			NodeID:     GetNodeID(),
			Param:      GetParam(),
			ScheduleID: GetScheduleID(),
		},
	}
}

// Save 保存单条数据
//
// 艹！自动更新统计信息
func (s *BaseSpider) Save(item interface{}) error {
	if err := SaveItem(item); err != nil {
		atomic.AddInt64(&s.Stats.Errors, 1)
		return err
	}
	atomic.AddInt64(&s.Stats.ItemsSaved, 1)
	return nil
}

// SaveBatch 批量保存数据
//
// 艹！一次保存多条，自动更新统计
func (s *BaseSpider) SaveBatch(items []interface{}) error {
	if err := SaveBatch(items); err != nil {
		atomic.AddInt64(&s.Stats.Errors, 1)
		return err
	}
	atomic.AddInt64(&s.Stats.ItemsSaved, int64(len(items)))
	return nil
}

// LogInfo 输出INFO日志
func (s *BaseSpider) LogInfo(format string, args ...interface{}) {
	LogInfo("[%s] "+format, append([]interface{}{s.Name}, args...)...)
}

// LogError 输出ERROR日志
//
// 艹！自动增加错误计数
func (s *BaseSpider) LogError(format string, args ...interface{}) {
	atomic.AddInt64(&s.Stats.Errors, 1)
	LogError("[%s] "+format, append([]interface{}{s.Name}, args...)...)
}

// LogWarn 输出WARNING日志
func (s *BaseSpider) LogWarn(format string, args ...interface{}) {
	LogWarn("[%s] "+format, append([]interface{}{s.Name}, args...)...)
}

// LogDebug 输出DEBUG日志
func (s *BaseSpider) LogDebug(format string, args ...interface{}) {
	LogDebug("[%s] "+format, append([]interface{}{s.Name}, args...)...)
}

// IncRequests 增加请求计数
//
// 艹！爬取网页后记得调用
func (s *BaseSpider) IncRequests() {
	atomic.AddInt64(&s.Stats.Requests, 1)
}

// IncErrors 增加错误计数
func (s *BaseSpider) IncErrors() {
	atomic.AddInt64(&s.Stats.Errors, 1)
}

// PrintStats 打印统计信息
//
// 艹！任务结束时自动调用
func (s *BaseSpider) PrintStats() {
	elapsed := time.Since(s.Stats.StartTime)
	s.LogInfo("========== 统计信息 ==========")
	s.LogInfo("运行时间: %v", elapsed)
	s.LogInfo("保存数据: %d 条", s.Stats.ItemsSaved)
	s.LogInfo("请求次数: %d 次", s.Stats.Requests)
	s.LogInfo("错误次数: %d 次", s.Stats.Errors)
	s.LogInfo("=============================")
}

// Execute 执行爬虫
//
// 艹！自动处理panic、打印统计、取消信号
// 用法：spider.Execute(spider)  // 把自己传进去
func (s *BaseSpider) Execute(spider Spider) error {
	// 创建可取消的context
	ctx, cancel := context.WithCancel(context.Background())
	s.Context.CancelFunc = cancel
	defer cancel()

	// 捕获panic
	defer func() {
		if r := recover(); r != nil {
			s.LogError("爬虫发生panic: %v", r)
		}
		// 打印统计信息
		s.PrintStats()
	}()

	s.LogInfo("开始执行爬虫: %s", s.Name)
	s.LogInfo("任务ID: %s", s.Context.TaskID)
	s.LogInfo("爬虫ID: %s", s.Context.SpiderID)

	// 运行爬虫
	if err := spider.Run(ctx); err != nil {
		s.LogError("爬虫执行失败: %v", err)
		return err
	}

	s.LogInfo("爬虫执行完成")
	return nil
}

// GetDuration 获取运行时长
func (s *BaseSpider) GetDuration() time.Duration {
	return time.Since(s.Stats.StartTime)
}

// ParseParam 解析任务参数到结构体
//
// 艹！方便调用，等价于ParseParamJSON
func (s *BaseSpider) ParseParam(v interface{}) error {
	return ParseParamJSON(v)
}
