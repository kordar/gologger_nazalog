# gologger_nazalog

包装`nazalog`对象，实现[日志门面](https://github.com/kordar/gologger)接口

## 安装

```go
go get github.com/kordar/gologger_nazalog v1.0.2
```

## 初始化

```go
import (
    logger "github.com/kordar/gologger"
    "github.com/kordar/gologger_nazalog"
    "github.com/q191201771/naza/pkg/nazalog"
)

_ = nazalog.Init(func(option *nazalog.Option) {
		option.Level = nazalog.LevelInfo
		option.Filename = "./logs/progress.log"
		option.IsRotateDaily = true
		option.LevelFlag = true
	})
defer nazalog.Sync()
logger.InitGlobal(gologger_nazalog.NewNazalogAdapt())
```
