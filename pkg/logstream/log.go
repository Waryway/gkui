package logstream

import (
	"context"
	"github.com/charmbracelet/log"
	"os"
)

func InitLogStream(ctx context.Context, cancelFunc context.CancelFunc) LogStream {
	var Logger LogStream
	Logger.New(ctx, cancelFunc)
	return Logger
}

type LogStream struct {
	StdErr  *log.Logger
	StdLog  *log.Logger
	ErrChan *chan Message
	StdChan *chan Message
	Ctx     context.Context
}

func (ls *LogStream) New(ctx context.Context, cancelFunc context.CancelFunc) *LogStream {
	ls.StdErr = log.New(os.Stderr)
	ls.StdLog = log.New(os.Stdout)
	ls.StdErr.SetLevel(log.ErrorLevel)
	ls.StdLog.SetLevel(log.DebugLevel)

	errChan := make(chan Message)
	stdChan := make(chan Message)
	ls.ErrChan = &errChan
	ls.StdChan = &stdChan
	ls.Ctx = ctx

	go func() {
		for {
			select {
			case errMsg := <-*ls.ErrChan:
				if errMsg.Level == log.FatalLevel {
					cancelFunc()
				}
				ls.writeStdLog(ls.StdErr, errMsg)
			case stdMsg := <-*ls.StdChan:
				if stdMsg.Level == log.FatalLevel {
					cancelFunc()
				}
				ls.writeStdLog(ls.StdLog, stdMsg)

			case <-ls.Ctx.Done():
				return
			}
		}
	}()

	return ls
}

func (ls *LogStream) DebugLog(msg interface{}, fields ...interface{}) {
	ls.Log(log.DebugLevel, msg, fields...)
}
func (ls *LogStream) InfoLog(msg interface{}, fields ...interface{}) {
	ls.Log(log.InfoLevel, msg, fields...)
}
func (ls *LogStream) WarnLog(msg interface{}, fields ...interface{}) {
	ls.Log(log.WarnLevel, msg, fields...)
}
func (ls *LogStream) ErrorLog(msg interface{}, fields ...interface{}) {
	ls.Log(log.ErrorLevel, msg, fields...)
}
func (ls *LogStream) FatalLog(msg interface{}, fields ...interface{}) {
	ls.Log(log.FatalLevel, msg, fields...)
}

func (ls *LogStream) Log(level log.Level, msg interface{}, fields ...interface{}) {
	streamMsg := Message{
		Level:  level,
		Msg:    msg,
		Fields: fields,
	}
	*ls.StdChan <- streamMsg
}

func (ls *LogStream) Err(level log.Level, msg interface{}, fields ...interface{}) {
	streamMsg := Message{
		Level:  level,
		Msg:    msg,
		Fields: fields,
	}
	*ls.ErrChan <- streamMsg
}

func (ls *LogStream) writeStdLog(logger *log.Logger, msg Message) {
	switch msg.Level {
	case log.DebugLevel:
		logger.Debug(msg.Msg, msg.Fields...)
	case log.InfoLevel:
		logger.Info(msg.Msg, msg.Fields...)
	case log.WarnLevel:
		logger.Warn(msg.Msg, msg.Fields...)
	case log.ErrorLevel:
		logger.Error(msg.Msg, msg.Fields...)
	case log.FatalLevel:
		logger.Fatal(msg.Msg, msg.Fields...)
	}
}
