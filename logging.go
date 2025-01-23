// Package logging
// @Author Clover
// @Data 2024/7/18 上午10:24:00
// @Desc 日志输出
package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	defaultProjectKey = "project"
)

var (
	logfile      *os.File
	once         sync.Once
	logPath      string              // 日志文件路径
	ProjectKey   = defaultProjectKey // 项目唯一标识
	projectName  string              // 项目名称
	maxLogSize   int64               // 最大日志文件大小
	monitorTimer *time.Ticker        // 日志大小监控计时器
)

// Config 用于配置日志记录器
type Config struct {
	LogPath             string        // 日志文件路径
	ProjectKey          string        // 项目唯一标识
	ProjectName         string        // 项目名称
	MaxLogSize          int64         // 最大日志文件大小 (字节)
	MonitorInterval     time.Duration // 监控日志大小的间隔时间
	EnableConsoleOutput bool          // 是否启用控制台输出
	EnableFileOutput    bool          // 是否启用文件输出
}

// InitLogger 初始化日志记录器
func InitLogger(config Config) {
	logPath = config.LogPath
	ProjectKey = config.ProjectKey
	projectName = config.ProjectName
	maxLogSize = config.MaxLogSize

	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"

	var writers []io.Writer

	if config.EnableConsoleOutput {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if config.EnableFileOutput {
		_, err := validLogPath(logPath, true)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to validate log path")
		}

		logfile, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal().Err(err).Msg("Error opening log file")
		}
		writers = append(writers, logfile)
	}

	multi := zerolog.MultiLevelWriter(writers...)
	// 直接使用 log.Logger 作为基础日志记录器，并设置输出、时间戳和项目名称字段
	log.Logger = log.Output(multi).With().Timestamp().Str(ProjectKey, projectName).Logger()

	if config.EnableFileOutput && config.MonitorInterval > 0 {
		monitorTimer = time.NewTicker(config.MonitorInterval)
		go monitorLogSize(monitorTimer.C)
	}
}

// SetField 设置字段信息k-v
func SetField(fields map[string]interface{}) {
	// 直接使用 log.Logger
	tmpLogger := log.With().Fields(fields).Logger()
	log.Logger = tmpLogger // 设置
}

// monitorLogSize 监控日志文件大小并在超过限制时清除日志文件
func monitorLogSize(ticker <-chan time.Time) {
	for range ticker {
		// Get the current log file size
		fi, err := logfile.Stat()
		if err != nil {
			log.Error().Err(err).Msg("Error getting file info")
			continue
		}

		if fi.Size() > maxLogSize {
			log.Info().Msg("Log file size exceeds limit. Clearing log file.")
			clearLogFile()
		}
	}
}

func clearLogFile() {
	var err error
	if err = logfile.Close(); err != nil {
		log.Error().Err(err).Msg("Error closing log file before truncation")
		return
	}

	// Truncate the log file to clear its content
	if err := os.Truncate(logPath, 0); err != nil {
		log.Error().Err(err).Msg("Error truncating log file")
		return
	}

	// Reopen the log file
	logfile, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("Error reopening log file after truncation")
		return
	}

	// Update the zerolog writer with the new file descriptor
	writers := []io.Writer{zerolog.ConsoleWriter{Out: os.Stderr}}
	if logfile != nil {
		writers = append(writers, logfile)
	}
	multi := zerolog.MultiLevelWriter(writers...)
	// 直接更新 log.Logger 的输出
	log.Logger = log.Output(multi).With().Timestamp().Str("sdk", projectName).Logger()

	log.Info().Msg("Log file cleared successfully.")
}

// Close 关闭日志文件和监控计时器
func Close() {
	once.Do(func() {
		if logfile != nil {
			err := logfile.Close()
			if err != nil {
				log.Error().Msgf("Error closing log file: %v", err)
			}
			logfile = nil
		}
		if monitorTimer != nil {
			monitorTimer.Stop()
		}
	})
}

// Info 定义简化的日志函数
func Info(msg string, fields ...map[string]interface{}) {
	event := log.Info()
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

func Error(msg string, fields ...map[string]interface{}) {
	event := log.Error()
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

func ErrorWithErr(err error, msg string, fields ...map[string]interface{}) {
	event := log.Error()
	event.Err(err)
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

func Debug(msg string, fields ...map[string]interface{}) {
	event := log.Debug()
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

func Warn(msg string, fields ...map[string]interface{}) {
	event := log.Warn()
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

func WarnWithErr(err error, msg string, fields ...map[string]interface{}) {
	event := log.Warn()
	event.Err(err)
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
}

func Fatal(msg string, exitCode int, fields ...map[string]interface{}) {
	event := log.Fatal()
	for _, field := range fields {
		for k, v := range field {
			event = event.Interface(k, v)
		}
	}
	event.Msg(msg)
	os.Exit(exitCode)
}

func validLogPath(path string, isCreate bool) (bool, error) {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if isCreate {
			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				return false, fmt.Errorf("error creating log directory: %w", err)
			}
		} else {
			return false, fmt.Errorf("log directory does not exist: %w", err)
		}
	}
	return true, nil
}

// Logger 定义一个全局的 LogBuffer
var Logger = NewLogBuffer()

// LogEntry 定义一个结构体来存储日志消息
type LogEntry struct {
	Level   zerolog.Level
	Message string
	Fields  map[string]interface{}
}

// LogBuffer 用于存储日志的缓冲区
type LogBuffer struct {
	entries []LogEntry
	mu      sync.Mutex
	active  bool // 是否激活缓冲模式
}

// NewLogBuffer 创建一个新的日志缓冲区
func NewLogBuffer() *LogBuffer {
	return &LogBuffer{
		entries: make([]LogEntry, 0),
		active:  true, // 初始激活缓冲模式
	}
}

// AddEntry 向缓冲区中添加一个日志条目
func (lb *LogBuffer) AddEntry(entry LogEntry) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	if lb.active {
		lb.entries = append(lb.entries, entry)
	} else {
		// 直接输出日志
		evt := log.WithLevel(entry.Level).Fields(entry.Fields)
		evt.Msg(entry.Message)
	}
}

// Flush 清空缓冲区，并根据日志等级输出日志
func (lb *LogBuffer) Flush(minLevel zerolog.Level) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	for _, entry := range lb.entries {
		if entry.Level >= minLevel {
			evt := log.WithLevel(entry.Level).Fields(entry.Fields)
			evt.Msg(entry.Message)
		}
	}
	// 清空缓冲区
	lb.entries = make([]LogEntry, 0)
}

// SetActive 设置缓冲区的激活状态
func (lb *LogBuffer) SetActive(active bool) {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	lb.active = active
}

func init() {
	// 初始化一个默认的 Logger
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stderr})
	multi := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stderr})

	log.Logger = log.Output(multi).With().Timestamp().Logger()
}
