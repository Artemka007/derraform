package logging

import (
    "fmt"
    "io"
    "os"
    "strings"
    "time"
    
    "github.com/fatih/color"
)

type LogLevel int

const (
    DEBUG LogLevel = iota
    INFO
    WARN
    ERROR
    FATAL
)

var (
    debugColor = color.New(color.FgHiBlue)
    infoColor  = color.New(color.FgHiGreen)
    warnColor  = color.New(color.FgHiYellow)
    errorColor = color.New(color.FgHiRed)
    fatalColor = color.New(color.FgHiRed, color.Bold)
)

type Logger struct {
    level  LogLevel
    writer io.Writer
}

func NewLogger(level LogLevel) *Logger {
    return &Logger{
        level:  level,
        writer: os.Stdout,
    }
}

func (l *Logger) log(level LogLevel, levelStr, msg string, args ...interface{}) {
    if level < l.level {
        return
    }

    timestamp := time.Now().Format("2006-01-02 15:04:05")
    message := fmt.Sprintf(msg, args...)
    
    var output string
    switch level {
    case DEBUG:
        output = debugColor.Sprintf("[%s] DEBUG: %s", timestamp, message)
    case INFO:
        output = infoColor.Sprintf("[%s] INFO: %s", timestamp, message)
    case WARN:
        output = warnColor.Sprintf("[%s] WARN: %s", timestamp, message)
    case ERROR:
        output = errorColor.Sprintf("[%s] ERROR: %s", timestamp, message)
    case FATAL:
        output = fatalColor.Sprintf("[%s] FATAL: %s", timestamp, message)
    }

    fmt.Fprintln(l.writer, output)
}

func (l *Logger) Debug(msg string, args ...interface{}) {
    l.log(DEBUG, "DEBUG", msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
    l.log(INFO, "INFO", msg, args...)
}

func (l *Logger) Warn(msg string, args ...interface{}) {
    l.log(WARN, "WARN", msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
    l.log(ERROR, "ERROR", msg, args...)
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
    l.log(FATAL, "FATAL", msg, args...)
    os.Exit(1)
}