package stdlog

import (
	"applicationDesignTest/utils/log"
	"fmt"
	stdlog "log"
	"os"
)

var _ log.Logger = (*Logger)(nil)

type Logger struct {
	l *stdlog.Logger
}

func NewLogger(prefix string) *Logger {
	if prefix != "" {
		prefix += " > "
	}
	return &Logger{
		l: stdlog.New(os.Stdout, prefix, stdlog.Ldate|stdlog.Ltime),
	}
}

func (l *Logger) Errorf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	l.l.Printf("[Error]: %s\n", msg)
}

func (l *Logger) Warnf(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	l.l.Printf("[Warn]: %s\n", msg)
}

func (l *Logger) Infof(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	l.l.Printf("[Info]: %s\n", msg)
}
