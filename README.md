# gologger-nazalog

为 `log/slog` 提供一个输出到 [nazalog](https://github.com/q191201771/naza/tree/master/pkg/nazalog) 的 `slog.Handler` 实现。

## 安装

```bash
go get github.com/kordar/gologger_nazalog@v0.1.0
```

## 使用

### 基础用法

```go
package main

import (
	"log/slog"

	"github.com/kordar/gologger_nazalog"
	"github.com/q191201771/naza/pkg/nazalog"
)

func main() {
	_ = nazalog.Init(func(o *nazalog.Option) {
		o.IsToStdout = true
		o.Filename = ""
	})

	log := gologger_nazalog.NewSlogLogger(nil, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	log.Info("hello", "k", "v")
	log.With("a", 1).WithGroup("g").Warn("warn", "x", true)
}
```

### 自定义 nazalog.Logger

```go
package main

import (
	"log/slog"

	"github.com/kordar/gologger_nazalog"
	"github.com/q191201771/naza/pkg/nazalog"
)

func main() {
	l, _ := nazalog.New(func(o *nazalog.Option) {
		o.IsToStdout = false
		o.Filename = "./app.log"
	})

	log := gologger_nazalog.NewSlogLogger(l, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		AddSource: true,
	})

	log.Info("m", "k", "v")
}
```

## API

- `NewSlogHandler(l nazalog.Logger, opts *slog.HandlerOptions) *SlogHandler`
- `NewSlogLogger(l nazalog.Logger, opts *slog.HandlerOptions) *slog.Logger`

