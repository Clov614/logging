package logging

import (
	"testing"
	"time"
)

func TestInitLoggerAndUsage(t *testing.T) {
	// 定义一个简单的配置
	config := Config{
		LogPath:             "./test.log",
		ProjectKey:          "project_key",
		ProjectName:         "testProject",
		MaxLogSize:          1024 * 1024, // 1MB
		MonitorInterval:     5 * time.Second,
		EnableConsoleOutput: true,
		EnableFileOutput:    true,
	}

	// 初始化日志记录器
	InitLogger(config)

	// 使用日志记录器记录一些信息
	Info("This is an info message.")
	Error("This is an error message.")
	Debug("This is a debug message.")
	Warn("This is a warning message.")

	// 使用 NewLogger 创建一个新的 Logger 实例并记录日志
	fields := map[string]interface{}{
		"user": "testuser",
		"id":   123,
	}
	logger := NewLogger(fields)
	logger.Info().Msg("This is a message from a new logger.")

	// 关闭日志记录器
	Close()
}
