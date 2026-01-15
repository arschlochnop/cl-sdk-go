package crawlab

import (
	"context"
	"fmt"
	"math"
	"time"
)

// RetryFunc 可重试的函数类型
type RetryFunc func() error

// Retry 重试执行函数
//
// 艹！失败自动重试，支持固定延迟
// maxRetries: 最大重试次数（0表示不重试）
// delay: 每次重试之间的延迟
func Retry(fn RetryFunc, maxRetries int, delay time.Duration) error {
	return RetryWithContext(context.Background(), fn, maxRetries, delay)
}

// RetryWithContext 带Context的重试执行
//
// 艹！支持取消和超时
func RetryWithContext(ctx context.Context, fn RetryFunc, maxRetries int, delay time.Duration) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// 检查context是否取消
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled: %w", ctx.Err())
		default:
		}

		// 执行函数
		err := fn()
		if err == nil {
			// 成功
			if attempt > 0 {
				LogDebug("Retry succeeded on attempt %d", attempt+1)
			}
			return nil
		}

		lastErr = err

		// 已经是最后一次尝试
		if attempt >= maxRetries {
			break
		}

		// 记录重试日志
		LogWarn("Attempt %d failed: %v, retrying in %v...", attempt+1, err, delay)

		// 等待后重试
		select {
		case <-time.After(delay):
			// 继续重试
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled during delay: %w", ctx.Err())
		}
	}

	return fmt.Errorf("all %d attempts failed, last error: %w", maxRetries+1, lastErr)
}

// RetryWithBackoff 带指数退避的重试
//
// 艹！延迟时间指数增长：delay, delay*2, delay*4, delay*8...
// maxDelay: 最大延迟时间
func RetryWithBackoff(ctx context.Context, fn RetryFunc, maxRetries int, initialDelay, maxDelay time.Duration) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// 检查context是否取消
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled: %w", ctx.Err())
		default:
		}

		// 执行函数
		err := fn()
		if err == nil {
			if attempt > 0 {
				LogDebug("Retry with backoff succeeded on attempt %d", attempt+1)
			}
			return nil
		}

		lastErr = err

		// 已经是最后一次尝试
		if attempt >= maxRetries {
			break
		}

		// 计算退避延迟（指数增长）
		backoffDelay := time.Duration(float64(initialDelay) * math.Pow(2, float64(attempt)))
		if backoffDelay > maxDelay {
			backoffDelay = maxDelay
		}

		// 记录重试日志
		LogWarn("Attempt %d failed: %v, retrying in %v...", attempt+1, err, backoffDelay)

		// 等待后重试
		select {
		case <-time.After(backoffDelay):
			// 继续重试
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled during backoff: %w", ctx.Err())
		}
	}

	return fmt.Errorf("all %d attempts failed with backoff, last error: %w", maxRetries+1, lastErr)
}

// RetryIf 条件重试
//
// 艹！只有shouldRetry返回true时才重试
// 用于某些错误不需要重试的场景
func RetryIf(ctx context.Context, fn RetryFunc, shouldRetry func(error) bool, maxRetries int, delay time.Duration) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// 检查context是否取消
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled: %w", ctx.Err())
		default:
		}

		// 执行函数
		err := fn()
		if err == nil {
			if attempt > 0 {
				LogDebug("Retry succeeded on attempt %d", attempt+1)
			}
			return nil
		}

		lastErr = err

		// 检查是否应该重试
		if !shouldRetry(err) {
			LogWarn("Error is not retryable: %v", err)
			return err
		}

		// 已经是最后一次尝试
		if attempt >= maxRetries {
			break
		}

		// 记录重试日志
		LogWarn("Attempt %d failed: %v, retrying in %v...", attempt+1, err, delay)

		// 等待后重试
		select {
		case <-time.After(delay):
			// 继续重试
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled during delay: %w", ctx.Err())
		}
	}

	return fmt.Errorf("all %d conditional retries failed, last error: %w", maxRetries+1, lastErr)
}
