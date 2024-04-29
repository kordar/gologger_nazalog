package nazalog

import (
	logger "github.com/kordar/gologger"
	"github.com/q191201771/naza/pkg/nazalog"
	"testing"
)

func InitLogger() {
	_ = nazalog.Init(func(option *nazalog.Option) {
		option.Level = nazalog.LevelInfo
		option.Filename = "./logs/progress.log"
		option.IsRotateDaily = true
		option.LevelFlag = true
	})
	defer nazalog.Sync()
	logger.InitGlobal(NewNazalogAdapt())
}

func TestT22(t *testing.T) {
	InitLogger()
	logger.Infof("this is info!")
}
