package env

import (
	"fmt"
	"gkui/pkg/bootstrap"
)

var (
	BootLoader bootstrap.Config[Settings]
	Config     Settings
)

func init() {
	Config = Settings{}
}

type Settings struct {
	Hello string `yaml:"hello"`
}

func (s Settings) InitializeSettings() (Settings, chan bootstrap.Config[Settings], chan error) {
	BootLoader = bootstrap.Config[Settings]{BootStrap: Settings{}}
	BootLoader.Init().Load()
	Config = BootLoader.BootStrap
	return Config, BootLoader.BootCh, BootLoader.ErrCh
}

func (s Settings) Save() {
	BootLoader.BootStrap = Config
	fmt.Print(Config.Hello)
	BootLoader.Save()
}
