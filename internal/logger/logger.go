package log

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type Level int32

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	default:
		return "unknown"
	}
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	SetLevel(level Level)
	GetLevel() Level
}

type BaseLogger struct {
	level Level
	mu    sync.Mutex
	impl  BaseImplementation
}

type BaseImplementation interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Fatal(msg string)
}

func New(writer io.Writer, level Level) Logger {
	return &BaseLogger{
		level: level,
		impl:  newStdImplementation(writer),
		mu:    sync.Mutex{},
	}
}

func (l *BaseLogger) SetLevel(level Level) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if level < DebugLevel || level > FatalLevel {
		panic("invalid log level")
	}
	l.level = level
}

func (l *BaseLogger) GetLevel() Level {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

func (l *BaseLogger) canLogAt(v Level) bool {
	return v >= l.GetLevel()
}

func (l *BaseLogger) Debug(args ...interface{}) {
	if !l.canLogAt(DebugLevel) {
		return
	}
	l.impl.Debug(fmt.Sprint(args...))
}

func (l *BaseLogger) Info(args ...interface{}) {
	if !l.canLogAt(InfoLevel) {
		return
	}
	l.impl.Info(fmt.Sprint(args...))
}

func (l *BaseLogger) Warn(args ...interface{}) {
	if !l.canLogAt(WarnLevel) {
		return
	}
	l.impl.Warn(fmt.Sprint(args...))
}

func (l *BaseLogger) Error(args ...interface{}) {
	if !l.canLogAt(ErrorLevel) {
		return
	}
	l.impl.Error(fmt.Sprint(args...))
}

func (l *BaseLogger) Fatal(args ...interface{}) {
	l.impl.Fatal(fmt.Sprint(args...))
	os.Exit(1)
}

func (l *BaseLogger) Debugf(format string, args ...interface{}) {
	if !l.canLogAt(DebugLevel) {
		return
	}
	l.impl.Debug(fmt.Sprintf(format, args...))
}

func (l *BaseLogger) Infof(format string, args ...interface{}) {
	if !l.canLogAt(InfoLevel) {
		return
	}
	l.impl.Info(fmt.Sprintf(format, args...))
}

func (l *BaseLogger) Warnf(format string, args ...interface{}) {
	if !l.canLogAt(WarnLevel) {
		return
	}
	l.impl.Warn(fmt.Sprintf(format, args...))
}

func (l *BaseLogger) Errorf(format string, args ...interface{}) {
	if !l.canLogAt(ErrorLevel) {
		return
	}
	l.impl.Error(fmt.Sprintf(format, args...))
}

func (l *BaseLogger) Fatalf(format string, args ...interface{}) {
	l.impl.Fatal(fmt.Sprintf(format, args...))
	os.Exit(1)
}
