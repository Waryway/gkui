package logstream

import "github.com/charmbracelet/log"

type Message struct {
	Level  log.Level
	Msg    interface{}
	Fields []interface{}
}
