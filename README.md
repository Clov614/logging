# logging

[![Go Reference](https://pkg.go.dev/badge/github.com/Clov614/logging.svg)](https://pkg.go.dev/github.com/Clov614/logging)

一个基于 zerolog 的 Go 日志库，提供灵活的配置选项，包括文件输出、日志轮转和自定义字段。

## 特性

*   支持控制台和文件输出
*   可配置的日志文件大小限制和轮转
*   支持自定义日志字段
*   简洁的 API

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
            ProjectKey:          "project", // 项目唯一标识，默认为 "project"
            ProjectName:         "my-app",
            MaxLogSize:          10 * 1024 * 1024, // 10MB, 日志文件最大大小
            MonitorInterval:     4 * time.Hour,    // 日志文件大小监控间隔
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

3. **设置全局日志字段**:

    使用 `logging.SetField()` 函数可以设置全局日志的字段。之后所有的日志记录都会包含这些字段。

    ```golang
    logging.SetField(map[string]interface{}{
        "component": "main",
    })
    logging.Info().Msg("This log message will contain the 'component' field.")
    ```

4. **关闭日志**:

    在程序结束时，调用 `logging.Close()` 关闭日志文件和监控计时器，以确保所有日志信息都已写入磁盘。

    ```golang
    defer logging.Close()
    ```

## 配置选项

*   **`LogPath`**: 日志文件的路径。
*   **`ProjectKey`**: 项目唯一标识，用于区分不同项目的日志，默认为 `"project"`。
*   **`ProjectName`**: 项目名称，用于在日志中标识项目。
*   **`MaxLogSize`**: 日志文件的最大大小（单位：字节）。当日志文件大小超过此限制时，将自动**清空并重新创建**日志文件。
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
        ProjectKey:          "project",
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

    // 设置全局字段
    logging.SetField(map[string]interface{}{
        "component": "api",
    })
    logging.Info().Msg("API 请求成功")

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

replace github.com/Clov614/logging => ./logging // 假设 logging 包位于同一仓库的 ./logging 目录中
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
        ProjectKey:          "projectA",
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

	// 设置全局字段
	logging.SetField(map[string]interface{}{
		"module": "projectA",
	})
	logging.Info("This is a log message from projectA.")
}
```

**projectB/main.go:**

```golang
package main

import (
	"time"

	"github.com/Clov614/logging"
	// "multi-project-example/projectA" // projectB 仅记录自己的日志，不再调用 projectA
)

func init() {
	logConfig := logging.Config{
		LogPath:             "./log/projectB.log",
        ProjectKey:          "projectB",
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

	// 设置全局字段
	logging.SetField(map[string]interface{}{
		"module": "projectB",
	})
	logging.Info("This is a log message from projectB.")
}
```

**运行结果:**

分别运行 `projectA` 和 `projectB` 将生成两个日志文件：`projectA.log` 和 `projectB.log`。

**projectA.log:**

```text
{"level":"info","sdk":"projectA","projectA":"projectA","time":"2024-07-19 15:30:00","message":"This is project A."}
{"level":"info","sdk":"projectA","projectA":"projectA","module":"projectA","time":"2024-07-19 15:30:00","message":"This is a log message from projectA."}
```

**projectB.log:**

```text
{"level":"info","sdk":"projectB","projectB":"projectB","time":"2024-07-19 15:30:00","message":"This is project B."}
{"level":"info","sdk":"projectB","projectB":"projectB","module":"projectB","time":"2024-07-19 15:30:00","message":"This is a log message from projectB."}
```

**说明:**

*   每个项目都有自己的 `init()` 函数，用于初始化 `logging` 包，并配置不同的日志文件路径、项目名称和项目唯一标识。
*   `projectB` 不再调用 `projectA` 的 `main()` 函数，只负责记录自己的日志。
*   每个项目都可以使用 `logging.Info()` 记录日志，也可以使用 `logging.SetField()` 设置全局日志字段。
*   运行 `projectA` 和 `projectB` 会分别生成各自的日志文件。

这个示例展示了如何在多个相互依赖的项目中使用 `logging` 包，并保持各自的日志配置和输出。每个项目都可以独立配置日志记录器，并记录到不同的日志文件中。