package crawlab

import (
	"os"
	"strconv"
	"time"
)

// Config 爬虫配置
//
// 艹！从环境变量加载配置，设置合理的默认值
type Config struct {
	// Crawlab环境信息
	TaskID     string // 任务ID
	SpiderID   string // 爬虫ID
	NodeID     string // 节点ID
	Param      string // 任务参数
	ScheduleID string // 调度ID

	// 性能配置
	MaxRetries     int           // 最大重试次数（默认3）
	RetryDelay     time.Duration // 重试延迟（默认2秒）
	RequestTimeout time.Duration // 请求超时（默认30秒）
	MaxConcurrency int           // 最大并发数（默认10）
	BatchSize      int           // 批量保存大小（默认100）
}

// LoadConfig 从环境变量加载配置
//
// 艹！自动设置默认值，开箱即用
func LoadConfig() *Config {
	cfg := &Config{
		// 加载Crawlab环境变量
		TaskID:     GetTaskID(),
		SpiderID:   GetSpiderID(),
		NodeID:     GetNodeID(),
		Param:      GetParam(),
		ScheduleID: GetScheduleID(),

		// 默认性能配置
		MaxRetries:     3,
		RetryDelay:     2 * time.Second,
		RequestTimeout: 30 * time.Second,
		MaxConcurrency: 10,
		BatchSize:      100,
	}

	// 从环境变量覆盖配置
	cfg.MaxRetries = cfg.GetEnvInt("CRAWLAB_MAX_RETRIES", cfg.MaxRetries)
	cfg.RetryDelay = cfg.GetEnvDuration("CRAWLAB_RETRY_DELAY", cfg.RetryDelay)
	cfg.RequestTimeout = cfg.GetEnvDuration("CRAWLAB_REQUEST_TIMEOUT", cfg.RequestTimeout)
	cfg.MaxConcurrency = cfg.GetEnvInt("CRAWLAB_MAX_CONCURRENCY", cfg.MaxConcurrency)
	cfg.BatchSize = cfg.GetEnvInt("CRAWLAB_BATCH_SIZE", cfg.BatchSize)

	return cfg
}

// GetEnvInt 从环境变量获取整数值
//
// 艹！获取失败或解析失败返回默认值
func (c *Config) GetEnvInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		LogWarn("Failed to parse %s as int: %v, using default %d", key, err, defaultVal)
		return defaultVal
	}

	return val
}

// GetEnvDuration 从环境变量获取时间间隔
//
// 艹！支持格式：1s, 2m, 3h 等
// 解析失败返回默认值
func (c *Config) GetEnvDuration(key string, defaultVal time.Duration) time.Duration {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}

	val, err := time.ParseDuration(valStr)
	if err != nil {
		LogWarn("Failed to parse %s as duration: %v, using default %v", key, err, defaultVal)
		return defaultVal
	}

	return val
}

// GetEnvBool 从环境变量获取布尔值
//
// 艹！true/1/yes/on 都算true，其他算false
func (c *Config) GetEnvBool(key string, defaultVal bool) bool {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}

	switch valStr {
	case "true", "1", "yes", "on", "TRUE", "YES", "ON":
		return true
	case "false", "0", "no", "off", "FALSE", "NO", "OFF":
		return false
	default:
		LogWarn("Unknown boolean value for %s: %s, using default %v", key, valStr, defaultVal)
		return defaultVal
	}
}

// Validate 验证配置
//
// 艹！检查配置是否合理
func (c *Config) Validate() error {
	if c.MaxRetries < 0 {
		LogWarn("MaxRetries is negative, setting to 0")
		c.MaxRetries = 0
	}

	if c.RetryDelay < 0 {
		LogWarn("RetryDelay is negative, setting to 0")
		c.RetryDelay = 0
	}

	if c.RequestTimeout < 0 {
		LogWarn("RequestTimeout is negative, setting to 30s")
		c.RequestTimeout = 30 * time.Second
	}

	if c.MaxConcurrency < 1 {
		LogWarn("MaxConcurrency is less than 1, setting to 1")
		c.MaxConcurrency = 1
	}

	if c.BatchSize < 1 {
		LogWarn("BatchSize is less than 1, setting to 1")
		c.BatchSize = 1
	}

	return nil
}

// Print 打印配置信息
func (c *Config) Print() {
	LogInfo("========== 配置信息 ==========")
	LogInfo("TaskID: %s", c.TaskID)
	LogInfo("SpiderID: %s", c.SpiderID)
	LogInfo("NodeID: %s", c.NodeID)
	LogInfo("ScheduleID: %s", c.ScheduleID)
	LogInfo("MaxRetries: %d", c.MaxRetries)
	LogInfo("RetryDelay: %v", c.RetryDelay)
	LogInfo("RequestTimeout: %v", c.RequestTimeout)
	LogInfo("MaxConcurrency: %d", c.MaxConcurrency)
	LogInfo("BatchSize: %d", c.BatchSize)
	LogInfo("=============================")
}
