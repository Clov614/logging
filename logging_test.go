package logging

import (
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
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

	// 使用 SetField 设置全局字段并记录日志
	fields := map[string]interface{}{
		"user": "testuser",
		"id":   123,
	}
	SetField(fields)
	Info("This is a message with global fields.")

	// 关闭日志记录器
	Close()

	// 验证日志文件是否存在
	if _, err := os.Stat("./test.log"); os.IsNotExist(err) {
		t.Errorf("Log file was not created: %v", err)
	}

	// 清理测试日志文件
	os.Remove("./test.log")
}

func TestLogBuffer(t *testing.T) {
	// 创建一个 LogBuffer 实例, 默认不激活以直接输出日志
	buf := NewLogBuffer()

	// 添加一些日志条目到缓冲区
	buf.AddEntry(LogEntry{
		Level:   zerolog.InfoLevel,
		Message: "This is an info message.",
		Fields:  nil})

	buf.AddEntry(LogEntry{
		Level:   zerolog.ErrorLevel,
		Message: "This is an error message in buffer.",
		Fields:  nil})
	buf.AddEntry(LogEntry{
		Level:   zerolog.DebugLevel,
		Message: "This is a debug message in buffer.",
		Fields:  nil})
	// 刷新缓冲区，输出日志
	buf.Flush(zerolog.InfoLevel)

	// 激活缓冲区
	buf.SetActive(true)
	buf.AddEntry(LogEntry{
		Level:   zerolog.InfoLevel,
		Message: "This is another info message in buffer.",
		Fields:  nil})

	// 再次刷新缓冲区, 输出所有 >= InfoLevel 的日志
	buf.Flush(zerolog.InfoLevel)
}
