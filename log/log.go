package log

import (
	"io"
	"log"
	"sync/atomic"
)

const (
	LevelDebug = 0
	LevelInfo  = 1
	LevelWarn  = 2
	LevelError = 3
)

type Logger struct {
	level atomic.Int32
	debug *log.Logger
	info  *log.Logger
	warn  *log.Logger
	err   *log.Logger
}

func New(out io.Writer, level int) *Logger {
	l := &Logger{
		debug: log.New(out, "[DEBUG] ", 0),
		info:  log.New(out, "[INFO] ", 0),
		warn:  log.New(out, "[WARN] ", 0),
		err:   log.New(out, "[ERROR] ", 0),
	}
	l.level.Store(int32(level))
	return l
}

func (l *Logger) SetLevel(level int) {
	l.level.Store(int32(level))
}

func (l *Logger) Debug(msg string) {
	if l.level.Load() <= LevelDebug {
		l.debug.Println(msg)
	}
}

func (l *Logger) Debugf(format string, args ...any) {
	if l.level.Load() <= LevelDebug {
		l.debug.Printf(format, args...)
	}
}

func (l *Logger) Info(msg string) {
	if l.level.Load() <= LevelInfo {
		l.info.Println(msg)
	}
}

func (l *Logger) Infof(format string, args ...any) {
	if l.level.Load() <= LevelInfo {
		l.info.Printf(format, args...)
	}
}

func (l *Logger) Warn(msg string) {
	if l.level.Load() <= LevelWarn {
		l.warn.Println(msg)
	}
}

func (l *Logger) Warnf(format string, args ...any) {
	if l.level.Load() <= LevelWarn {
		l.warn.Printf(format, args...)
	}
}

func (l *Logger) Error(msg string) {
	if l.level.Load() <= LevelError {
		l.err.Println(msg)
	}
}

func (l *Logger) Errorf(format string, args ...any) {
	if l.level.Load() <= LevelError {
		l.err.Printf(format, args...)
	}
}
