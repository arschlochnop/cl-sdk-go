package crawlab

import (
	"encoding/json"
	"fmt"
	"os"
)

// IPCMessage IPC消息结构
type IPCMessage struct {
	IPC     bool        `json:"ipc"`
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

const (
	// MaxIPCMessageSize 单条IPC消息最大大小（5MB），考虑到JSON序列化开销
	MaxIPCMessageSize = 5 * 1024 * 1024

	// 环境变量名称
	EnvTaskID     = "CRAWLAB_TASK_ID"
	EnvSpiderID   = "CRAWLAB_SPIDER_ID"
	EnvNodeID     = "CRAWLAB_NODE_ID"
	EnvParam      = "CRAWLAB_TASK_PARAM"
	EnvScheduleID = "CRAWLAB_SCHEDULE_ID"
)

// SaveItem 保存单条数据到Crawlab
//
// 艹！自动检查数据大小，超过5MB会警告
func SaveItem(item interface{}) error {
	return SaveItems(item)
}

// SaveItems 保存多条数据到Crawlab
//
// 艹！每条数据单独发送，适合少量数据
// 如果数据量大，用SaveBatch批量发送
func SaveItems(items ...interface{}) error {
	for _, item := range items {
		// 检查数据大小
		data, err := json.Marshal(item)
		if err != nil {
			return fmt.Errorf("failed to marshal item: %w", err)
		}

		// 如果数据太大，输出警告
		if len(data) > MaxIPCMessageSize {
			LogWarn("Item size (%d bytes) exceeds recommended limit (%d bytes)", len(data), MaxIPCMessageSize)
			LogWarn("Consider splitting large data or using external storage")
		}

		if err := sendIPCMessage("data", item); err != nil {
			return fmt.Errorf("failed to save item: %w", err)
		}
	}
	return nil
}

// SaveBatch 批量保存数据（发送数组）
//
// 艹！一次发送整个数组，减少IPC次数，性能更好
// 适合大量数据，但注意总大小不要超过5MB
func SaveBatch(items []interface{}) error {
	if len(items) == 0 {
		return nil
	}

	// 检查批次总大小
	data, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("failed to marshal batch: %w", err)
	}

	if len(data) > MaxIPCMessageSize {
		LogWarn("Batch size (%d bytes) exceeds recommended limit (%d bytes)", len(data), MaxIPCMessageSize)
		LogWarn("Consider splitting into smaller batches")
	}

	return sendIPCMessage("data", items)
}

// Log 输出日志到stderr（会被Runner捕获为任务日志）
//
// 艹！基础日志输出
func Log(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[Crawlab] "+format+"\n", args...)
}

// LogInfo 输出INFO级别日志
func LogInfo(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[Crawlab] [INFO] "+format+"\n", args...)
}

// LogError 输出ERROR级别日志
//
// 艹！出错了就用这个
func LogError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[Crawlab] [ERROR] "+format+"\n", args...)
}

// LogWarn 输出WARNING级别日志
func LogWarn(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[Crawlab] [WARN] "+format+"\n", args...)
}

// LogDebug 输出DEBUG级别日志
func LogDebug(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "[Crawlab] [DEBUG] "+format+"\n", args...)
}

// GetTaskID 从环境变量获取当前任务ID
func GetTaskID() string {
	return os.Getenv(EnvTaskID)
}

// GetSpiderID 从环境变量获取爬虫ID
func GetSpiderID() string {
	return os.Getenv(EnvSpiderID)
}

// GetNodeID 从环境变量获取节点ID
func GetNodeID() string {
	return os.Getenv(EnvNodeID)
}

// GetParam 从环境变量获取任务参数（原始字符串）
func GetParam() string {
	return os.Getenv(EnvParam)
}

// GetScheduleID 从环境变量获取调度ID
func GetScheduleID() string {
	return os.Getenv(EnvScheduleID)
}

// MustGetEnv 获取环境变量，不存在则panic
//
// 艹！必须有的环境变量用这个
func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return value
}

// GetEnv 获取环境变量，不存在则返回默认值
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// ParseParamJSON 解析任务参数JSON到结构体
//
// 艹！自动获取CRAWLAB_TASK_PARAM并解析
// 用法：var params MyParams; ParseParamJSON(&params)
func ParseParamJSON(v interface{}) error {
	param := GetParam()
	if param == "" {
		return fmt.Errorf("task param is empty")
	}

	if err := json.Unmarshal([]byte(param), v); err != nil {
		return fmt.Errorf("failed to parse param JSON: %w", err)
	}

	return nil
}

// sendIPCMessage 发送IPC消息到stdout
//
// 艹！内部函数，别直接用
func sendIPCMessage(msgType string, payload interface{}) error {
	ipcMsg := IPCMessage{
		IPC:     true,
		Type:    msgType,
		Payload: payload,
	}

	data, err := json.Marshal(ipcMsg)
	if err != nil {
		return fmt.Errorf("failed to marshal IPC message: %w", err)
	}

	fmt.Println(string(data))
	return nil
}
