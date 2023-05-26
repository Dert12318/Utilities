package logrus

import (
	"io"
	"os"

	"github.com/Dert12318/Utilities/logs"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

const (
	JSONFormatter Formatter = "JSON"
	TextFormatter Formatter = "TEXT"
)

type (
	Formatter string

	Option struct {
		Level       log.Lvl
		LogFilePath string
		Formatter   Formatter
		Prefix      string
		Masking     logs.MaskedEncoder
	}

	logger struct {
		instance *logrus.Logger
		level    log.Lvl
		prefix   string
		masking  logs.MaskedEncoder
	}
)

func (l *logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

func (l *logger) Level() log.Lvl {
	return l.level
}

func (l *logger) SetLevel(v log.Lvl) {
	l.level = v
	l.instance.SetLevel(getLevel(v))
}

func (l *logger) SetHeader(header string) {

}

func (l *logger) Print(args ...interface{}) {
	if l.level != log.OFF {
		l.instance.Print(args...)
	}
}

func (l *logger) Println(args ...interface{}) {
	if l.level != log.OFF {
		l.instance.Println(args...)
	}
}

func (l *logger) Printf(format string, args ...interface{}) {
	if l.level != log.OFF {
		l.instance.Printf(format, args...)
	}
}

func (l *logger) Printj(j log.JSON) {
	if l.level != log.OFF {
		l.Printf("%+v\n", j)
	}
}

func (l *logger) Debug(args ...interface{}) {
	if l.level != log.OFF {
		l.instance.Debug(args...)
	}
}

func (l *logger) Debugf(format string, args ...interface{}) {
	if l.level != log.OFF {
		l.instance.Debugf(format, args...)
	}
}

func (l *logger) Debugj(j log.JSON) {
	if l.level != log.OFF {
		l.Debugf("%+v\n", j)
	}
}

func (l *logger) Info(args ...interface{}) {
	if l.level != log.OFF {
		l.instance.Info(args...)
	}
}

func (l *logger) Infof(format string, args ...interface{}) {
	if l.level != log.OFF {
		l.instance.Infof(format, args...)
	}
}

func (l *logger) Infoj(j log.JSON) {
	if l.level != log.OFF {
		l.Infof("%+v\n", j)
	}
}

func (l *logger) Warn(i ...interface{}) {
	if l.level != log.OFF {
		l.instance.Warn(i...)
	}
}

func (l *logger) Warnf(format string, i ...interface{}) {
	if l.level != log.OFF {
		l.instance.Warnf(format, i...)
	}
}

func (l *logger) Warnj(j log.JSON) {
	if l.level != log.OFF {
		l.Warnf("%+v\n", j)
	}
}

func (l *logger) Error(args ...interface{}) {
	if l.level != log.OFF {
		l.instance.Error(args...)
	}
}

func (l *logger) Errorf(format string, args ...interface{}) {
	if l.level != log.OFF {
		l.instance.Errorf(format, args...)
	}
}

func (l *logger) Errorj(j log.JSON) {
	if l.level != log.OFF {
		l.Errorf("%+v\n", j)
	}
}

func (l *logger) Fatal(args ...interface{}) {
	if l.level != log.OFF {
		l.instance.Fatal(args...)
	}
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	if l.level != log.OFF {
		l.instance.Fatalf(format, args...)
	}
}

func (l *logger) Fatalj(j log.JSON) {
	if l.level != log.OFF {
		l.Fatalf("%+v\n", j)
	}
}

func (l *logger) Panic(i ...interface{}) {
	if l.level != log.OFF {
		l.instance.Panic(i...)
	}
}

func (l *logger) Panicf(format string, i ...interface{}) {
	if l.level != log.OFF {
		l.instance.Panicf(format, i...)
	}
}

func (l *logger) Panicj(j log.JSON) {
	if l.level != log.OFF {
		l.Panicf("%+v\n", j)
	}
}

func (l *logger) Instance() interface{} {
	return l.instance
}

func (l logger) Log(msg string) {
	if l.level != log.OFF {
		l.instance.Info(msg)
	}
}

func (l *logger) Output() io.Writer {
	return l.instance.Out
}

func (l *logger) SetOutput(w io.Writer) {
	l.instance.Out = w
}

func (l *logger) Prefix() string {
	return l.prefix
}

func New(option *Option) (logs.Logger, error) {
	instance := logrus.New()

	switch option.Level {
	case log.INFO:
		instance.Level = logrus.InfoLevel
		break
	case log.DEBUG:
		instance.Level = logrus.DebugLevel
		break
	case log.WARN:
		instance.Level = logrus.WarnLevel
		break
	case log.ERROR:
		instance.Level = logrus.ErrorLevel
		break
	default:
		instance.Level = logrus.ErrorLevel
		break
	}

	var formatter logrus.Formatter

	if option.Formatter == JSONFormatter {
		formatter = &logrus.JSONFormatter{}
	} else {
		formatter = &logrus.TextFormatter{}
	}

	instance.Formatter = formatter

	// - check if log file path does exists
	if option.LogFilePath != "" {
		if _, err := os.Stat(option.LogFilePath); os.IsNotExist(err) {
			if _, err = os.Create(option.LogFilePath); err != nil {
				return nil, errors.Wrapf(err, "failed to create log file %s", option.LogFilePath)
			}
		}
		maps := lfshook.PathMap{
			logrus.InfoLevel:  option.LogFilePath,
			logrus.DebugLevel: option.LogFilePath,
			logrus.ErrorLevel: option.LogFilePath,
		}
		instance.Hooks.Add(lfshook.NewHook(maps, formatter))
	}

	return &logger{
		instance: instance,
		level:    option.Level,
		prefix:   option.Prefix,
		masking:  option.Masking,
	}, nil
}

func DefaultLog() logs.Logger {
	logger, _ := New(&Option{
		Level:     log.INFO,
		Formatter: TextFormatter,
	})
	return logger
}

func getLevel(lvl log.Lvl) logrus.Level {
	switch lvl {
	case log.INFO:
		return logrus.InfoLevel
	case log.DEBUG:
		return logrus.DebugLevel
	case log.WARN:
		return logrus.WarnLevel
	case log.ERROR:
		return logrus.ErrorLevel
	default:
		return logrus.ErrorLevel
	}
}
