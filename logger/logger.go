
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// LogLevel represents the severity of log messages
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

// String returns the string representation of log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger represents a structured logger
type Logger struct {
	level      LogLevel
	logger     *log.Logger
	component  string
	outputFile *os.File
}

var defaultLogger *Logger

// Config holds logger configuration
type Config struct {
	Level      string `mapstructure:"level"`
	Output     string `mapstructure:"output"`     // "console", "file", or "both"
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int64  `mapstructure:"max_size"`   // Maximum log file size in MB
	Component  string `mapstructure:"component"`
}

// Init initializes the default logger
func Init(config Config) error {
	level := parseLogLevel(config.Level)
	
	var writer io.Writer
	var outputFile *os.File
	
	switch config.Output {
	case "file":
		if err := os.MkdirAll(filepath.Dir(config.FilePath), 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %v", err)
		}
		
		file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open log file: %v", err)
		}
		writer = file
		outputFile = file
		
	case "both":
		if err := os.MkdirAll(filepath.Dir(config.FilePath), 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %v", err)
		}
		
		file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open log file: %v", err)
		}
		writer = io.MultiWriter(os.Stdout, file)
		outputFile = file
		
	default: // console
		writer = os.Stdout
	}

	defaultLogger = &Logger{
		level:      level,
		logger:     log.New(writer, "", 0),
		component:  config.Component,
		outputFile: outputFile,
	}

	return nil
}

// NewLogger creates a new logger for a specific component
func NewLogger(component string) *Logger {
	if defaultLogger == nil {
		// Fallback to console logging if not initialized
		defaultLogger = &Logger{
			level:     INFO,
			logger:    log.New(os.Stdout, "", 0),
			component: "default",
		}
	}

	return &Logger{
		level:      defaultLogger.level,
		logger:     defaultLogger.logger,
		component:  component,
		outputFile: defaultLogger.outputFile,
	}
}

// Close closes the logger and any open files
func (l *Logger) Close() error {
	if l.outputFile != nil {
		return l.outputFile.Close()
	}
	return nil
}

// SetLevel sets the logging level
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// Debug logs a debug message
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info logs an info message
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warning logs a warning message
func (l *Logger) Warning(format string, args ...interface{}) {
	l.log(WARNING, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
	os.Exit(1)
}

// log performs the actual logging
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	_, file, line, _ := runtime.Caller(2)
	file = filepath.Base(file)
	
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] [%s] [%s] %s:%d - %s", 
		timestamp, level.String(), l.component, file, line, message)
	
	l.logger.Println(logLine)
}

// parseLogLevel parses string log level to LogLevel
func parseLogLevel(level string) LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARNING", "WARN":
		return WARNING
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	default:
		return INFO
	}
}

// Global convenience functions
func Debug(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debug(format, args...)
	}
}

func Info(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Info(format, args...)
	}
}

func Warning(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warning(format, args...)
	}
}

func Error(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Error(format, args...)
	}
}

func Fatal(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Fatal(format, args...)
	}
}

func Close() error {
	if defaultLogger != nil {
		return defaultLogger.Close()
	}
	return nil
}
