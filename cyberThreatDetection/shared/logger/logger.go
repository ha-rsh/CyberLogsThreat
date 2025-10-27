package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	logger      *Logger
	levelColors = map[Level]string{
		DEBUG: "\033[36m", // Cyan
		INFO:  "\033[32m", // Green
		WARN:  "\033[33m", // Yellow
		ERROR: "\033[31m", // Red
		FATAL: "\033[35m", // Magenta
	}
	resetColor = "\033[0m"
)

type Logger struct {
	serviceName string
	minLevel    Level
}

func Init(serviceName string, minLevel Level) {
	logger = &Logger{
		serviceName: serviceName,
		minLevel:    minLevel,
	}
}

func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.minLevel {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelStr := l.getLevelString(level)
	message := fmt.Sprintf(format, args...)
	
	// Get caller info
	_, file, line, ok := runtime.Caller(2)
	fileInfo := ""
	if ok {
		parts := strings.Split(file, "/")
		fileName := parts[len(parts)-1]
		fileInfo = fmt.Sprintf("%s:%d", fileName, line)
	}

	// Format: [TIMESTAMP] [SERVICE] [LEVEL] [FILE:LINE] MESSAGE
	logMessage := fmt.Sprintf("%s[%s]%s %s[%s]%s %s[%-5s]%s %s[%s]%s %s",
		"\033[90m", timestamp, resetColor,
		"\033[34m", l.serviceName, resetColor,
		levelColors[level], levelStr, resetColor,
		"\033[90m", fileInfo, resetColor,
		message,
	)

	if level == FATAL {
		log.Fatal(logMessage)
	} else {
		fmt.Println(logMessage)
	}
}

func (l *Logger) getLevelString(level Level) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

func Debug(format string, args ...interface{}) {
	if logger == nil {
		Init("default", DEBUG)
	}
	logger.log(DEBUG, format, args...)
}

func Info(format string, args ...interface{}) {
	if logger == nil {
		Init("default", INFO)
	}
	logger.log(INFO, format, args...)
}

func Warn(format string, args ...interface{}) {
	if logger == nil {
		Init("default", WARN)
	}
	logger.log(WARN, format, args...)
}

func Error(format string, args ...interface{}) {
	if logger == nil {
		Init("default", ERROR)
	}
	logger.log(ERROR, format, args...)
}

func Fatal(format string, args ...interface{}) {
	if logger == nil {
		Init("default", FATAL)
	}
	logger.log(FATAL, format, args...)
	os.Exit(1)
}

// Convenience functions for HTTP logging
func HTTP(method, path string, statusCode int, duration time.Duration) {
	statusColor := "\033[32m" // Green
	if statusCode >= 400 && statusCode < 500 {
		statusColor = "\033[33m" // Yellow
	} else if statusCode >= 500 {
		statusColor = "\033[31m" // Red
	}

	Info("%s %-7s %s%s %s%d%s %s[%s]%s",
		"\033[1m", method, resetColor,
		path,
		statusColor, statusCode, resetColor,
		"\033[90m", duration.String(), resetColor,
	)
}

func Success(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("\033[1;32m✓\033[0m %s\n", message)
}

func Failure(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	fmt.Printf("\033[1;31m✗\033[0m %s\n", message)
}