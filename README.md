# logging

[![Go Reference](https://pkg.go.dev/badge/github.com/Clov614/logging.svg)](https://pkg.go.dev/github.com/Clov614/logging)

一个基于 zerolog 的 Go 日志库，提供灵活的配置选项，包括文件输出、日志轮转和自定义字段。

## 特性

* 支持控制台和文件输出
* 可配置的日志文件大小限制和轮转
* 支持自定义日志字段
* 简洁的 API

## 安装

```text
go get github.com/Clov614/logging
```

## 使用方法

1. **初始化**:

   首先，使用 `logging.Config` 结构体配置日志参数，然后调用 `logging.InitLogger()` 初始化日志库。

```golang
   package main

import (
   "time"
   "github.com/Clov614/logging"
)

func main() {
   logConfig := logging.Config{
      LogPath:             "./log/app.log",
      ProjectName:         "my-app",
      MaxLogSize:          10 * 1024 * 1024, // 10MB
      MonitorInterval:     4 * time.Hour,
      EnableConsoleOutput: true,
      EnableFileOutput:    true,
   }
   logging.InitLogger(logConfig)
   defer logging.Close()

   // ... other code ...
}

```

2. **记录日志**:

   使用 `logging` 包提供的函数记录不同级别的日志信息，例如 `Info`、`Error`、`Warn`、`Debug` 和 `Fatal`。可以添加自定义字段以提供更多上下文信息。

   ```golang
   logging.Info("启动程序", map[string]interface{}{"version": "1.0.0"})
   logging.Error("发生错误", map[string]interface{}{"error": "some error"})
   logging.Warn("警告信息", map[string]interface{}{"code": 123})
   logging.Debug("调试信息", map[string]interface{}{"data": "some data"})
   logging.Fatal("致命错误", 1, map[string]interface{}{"reason": "critical error"})
   ```

3. **创建自定义 Logger**:

   你可以基于 `baseLogger` 创建一个新的 `Logger` 实例，并添加自定义字段。

   ```golang
   myLogger := logging.NewLogger(map[string]interface{}{
   	"component": "main",
   })
   myLogger.Info().Msg("This is a log message from myLogger.")
   ```

4. **关闭日志**:

   在程序结束时，调用 `logging.Close()` 关闭日志文件和监控计时器，以确保所有日志信息都已写入磁盘。

   ```golang
   defer logging.Close()
   ```

## 配置选项

*   **`LogPath`**: 日志文件的路径。
*   **`ProjectName`**: 项目名称，用于在日志中标识项目。
*   **`MaxLogSize`**: 日志文件的最大大小（单位：字节）。当日志文件大小超过此限制时，将自动清空日志文件。
*   **`MonitorInterval`**: 监控日志文件大小的间隔时间。
*   **`EnableConsoleOutput`**: 是否启用控制台输出。
*   **`EnableFileOutput`**: 是否启用文件输出。

## 示例

以下是一个完整的示例，演示如何使用 `logging` 包记录不同级别的日志信息：

```golang
package main

import (
   "errors"
   "time"
   "github.com/Clov614/logging"
)

func main() {
   logConfig := logging.Config{
      LogPath:             "./log/application.log",
      ProjectName:         "my-application",
      MaxLogSize:          50 * 1024 * 1024, // 50MB
      MonitorInterval:     4 * time.Hour,
      EnableConsoleOutput: true,
      EnableFileOutput:    true,
   }
   logging.InitLogger(logConfig)
   defer logging.Close()

   logging.Info("应用程序启动", map[string]interface{}{"version": "1.2.3"})

   err := errors.New("示例错误")
   logging.ErrorWithErr(err, "处理请求时出错", map[string]interface{}{"request_id": "12345"})

   logging.Warn("资源不足", map[string]interface{}{"resource": "memory"})

   logging.Debug("进入调试模式", map[string]interface{}{"debug_level": 2})

   // 使用自定义字段创建新的 Logger
   myLogger := logging.NewLogger(map[string]interface{}{
      "component": "api",
   })
   myLogger.Info().Msg("API 请求成功")

   // ... other code ...
}

```

## 多项目依赖示例

假设我们有两个项目：`projectA` 和 `projectB`，其中 `projectB` 依赖于 `projectA`。

**项目结构:**

```
multi-project-example/
├── projectA/
│   └── main.go
├── projectB/
│   └── main.go
└── go.mod
```

**go.mod:**

```text
module multi-project-example

go 1.20

require github.com/Clov614/logging v0.0.0-00010101000000-000000000000

replace github.com/Clov614/logging => ./logging // 假设 logging 包位于同一仓库中
```

**projectA/main.go:**

```golang
package main

import (
	"time"

	"github.com/Clov614/logging"
)

func init() {
	logConfig := logging.Config{
		LogPath:             "./log/projectA.log",
		ProjectName:         "projectA",
		MaxLogSize:          10 * 1024 * 1024, // 10MB
		MonitorInterval:     4 * time.Hour,
		EnableConsoleOutput: true,
		EnableFileOutput:    true,
	}
	logging.InitLogger(logConfig)
}

func main() {
	defer logging.Close()

	logging.Info("This is project A.")

	// 使用自定义字段创建新的 Logger
	projectALogger := logging.NewLogger(map[string]interface{}{
		"module": "projectA",
	})
	projectALogger.Info().Msg("This is a log message from projectA's logger.")
}
```

**projectB/main.go:**

```golang
package main

import (
	"time"

	"github.com/Clov614/logging"
	"multi-project-example/projectA"
)

func init() {
	logConfig := logging.Config{
		LogPath:             "./log/projectB.log",
		ProjectName:         "projectB",
		MaxLogSize:          10 * 1024 * 1024, // 10MB
		MonitorInterval:     4 * time.Hour,
		EnableConsoleOutput: true,
		EnableFileOutput:    true,
	}
	logging.InitLogger(logConfig)
}

func main() {
	defer logging.Close()

	logging.Info("This is project B.")

	// 调用 projectA 的函数
	projectA.main()

	// 使用自定义字段创建新的 Logger
	projectBLogger := logging.NewLogger(map[string]interface{}{
		"module": "projectB",
	})
	projectBLogger.Info().Msg("This is a log message from projectB's logger.")
}
```

**运行结果:**

运行 `projectB` (因为它依赖于 `projectA`) 将生成两个日志文件：`projectA.log` 和 `projectB.log`。

**projectA.log:**

```text
{"level":"info","sdk":"projectA","time":"2024-07-19 15:30:00","message":"This is project A."}
{"level":"info","sdk":"projectA","module":"projectA","time":"2024-07-19 15:30:00","message":"This is a log message from projectA's logger."}
```

**projectB.log:**

```text
{"level":"info","sdk":"projectB","time":"2024-07-19 15:30:00","message":"This is project B."}
{"level":"info","sdk":"projectB","module":"projectB","time":"2024-07-19 15:30:00","message":"This is a log message from projectB's logger."}
```

**说明:**

*   每个项目都有自己的 `init()` 函数，用于初始化 `logging` 包，并配置不同的日志文件路径和项目名称。
*   `projectB` 导入 `projectA` 并调用其 `main()` 函数。
*   每个项目都可以使用 `logging.Info()` 记录日志，也可以使用 `logging.NewLogger()` 创建带有自定义字段的 `Logger` 实例。
*   运行 `projectB` 会同时触发 `projectA` 的日志记录。

这个示例展示了如何在多个相互依赖的项目中使用 `logging` 包，并保持各自的日志配置和输出。每个项目都可以独立配置日志记录器，并记录到不同的日志文件中。