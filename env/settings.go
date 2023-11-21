package env

import (
	"gkui/pkg/bootstrap"
	"gkui/pkg/logstream"
)

var (
	BootLoader bootstrap.Config[Settings]
	Config     Settings
	Logger     logstream.LogStream
)

func init() {
	Config = Settings{}
}

type Settings struct {
	Hello string `yaml:"hello"`
}

func (s Settings) InitializeSettings(ls logstream.LogStream) (Settings, chan bootstrap.Config[Settings], chan error) {
	Logger = ls
	BootLoader = bootstrap.Config[Settings]{BootStrap: Settings{}}
	BootLoader.Init(&Logger).Load()
	Config = BootLoader.BootStrap
	return Config, BootLoader.BootCh, BootLoader.ErrCh
}

func (s Settings) Save() {
	BootLoader.BootStrap = Config
	Logger.DebugLog("Saved Config")
	BootLoader.Save()
}
