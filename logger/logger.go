
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// LogLevel represents the logging level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

var (
	levelNames = map[LogLevel]string{
		DEBUG:   "DEBUG",
		INFO:    "INFO",
		WARNING: "WARN",
		ERROR:   "ERROR",
		FATAL:   "FATAL",
	}
)

// Config holds logger configuration
type Config struct {
	Level     string // debug, info, warning, error, fatal
	Output    string // console, file, both
	FilePath  string // path to log file
	MaxSize   int64  // max file size in MB
	Component string // component name for logging
}

// Logger represents a component-specific logger
type Logger struct {
	component string
	level     LogLevel
	output    io.Writer
	mu        sync.Mutex
}

var (
	globalLogger *Logger
	globalConfig Config
	logFile      *os.File
	mu           sync.Mutex
)

// Init initializes the global logger with configuration
func Init(config Config) error {
	mu.Lock()
	defer mu.Unlock()

	globalConfig = config

	// Parse log level
	level := parseLogLevel(config.Level)

	// Setup output writer
	var output io.Writer
	switch strings.ToLower(config.Output) {
	case "console":
		output = os.Stdout
	case "file":
		if config.FilePath == "" {
			config.FilePath = "./logs/blockchain.log"
		}
		
		// Ensure log directory exists
		if err := os.MkdirAll(filepath.Dir(config.FilePath), 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %v", err)
		}
		
		file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %v", err)
		}
		logFile = file
		output = file
	case "both":
		if config.FilePath == "" {
			config.FilePath = "./logs/blockchain.log"
		}
		
		// Ensure log directory exists
		if err := os.MkdirAll(filepath.Dir(config.FilePath), 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %v", err)
		}
		
		file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %v", err)
		}
		logFile = file
		output = io.MultiWriter(os.Stdout, file)
	default:
		output = os.Stdout
	}

	globalLogger = &Logger{
		component: config.Component,
		level:     level,
		output:    output,
	}

	// Log initialization
	globalLogger.logWithLevel(INFO, "Logger initialized", "level", levelNames[level], "output", config.Output)

	return nil
}

// NewLogger creates a new logger for a specific component
func NewLogger(component string) *Logger {
	if globalLogger == nil {
		// Fallback to console logging if not initialized
		return &Logger{
			component: component,
			level:     INFO,
			output:    os.Stdout,
		}
	}

	return &Logger{
		component: component,
		level:     globalLogger.level,
		output:    globalLogger.output,
	}
}

// parseLogLevel parses log level string to LogLevel
func parseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warning", "warn":
		return WARNING
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return INFO
	}
}

// logWithLevel logs a message with the specified level
func (l *Logger) logWithLevel(level LogLevel, message string, keyValues ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	levelStr := levelNames[level]
	
	// Build key-value pairs string
	var kvPairs strings.Builder
	for i := 0; i < len(keyValues); i += 2 {
		if i+1 < len(keyValues) {
			if kvPairs.Len() > 0 {
				kvPairs.WriteString(", ")
			}
			kvPairs.WriteString(fmt.Sprintf("%v=%v", keyValues[i], keyValues[i+1]))
		}
	}

	// Format: [TIMESTAMP] [LEVEL] [COMPONENT] MESSAGE key=value, key=value
	logLine := fmt.Sprintf("[%s] [%s] [%s] %s", timestamp, levelStr, l.component, message)
	if kvPairs.Len() > 0 {
		logLine += " " + kvPairs.String()
	}
	logLine += "\n"

	// Write to output
	if l.output != nil {
		l.output.Write([]byte(logLine))
	}

	// For FATAL level, exit the program
	if level == FATAL {
		os.Exit(1)
	}
}

// Debug logs a debug message
func (l *Logger) Debug(message string, keyValues ...interface{}) {
	l.logWithLevel(DEBUG, message, keyValues...)
}

// Info logs an info message
func (l *Logger) Info(message string, keyValues ...interface{}) {
	l.logWithLevel(INFO, message, keyValues...)
}

// Warning logs a warning message
func (l *Logger) Warning(message string, keyValues ...interface{}) {
	l.logWithLevel(WARNING, message, keyValues...)
}

// Error logs an error message
func (l *Logger) Error(message string, keyValues ...interface{}) {
	l.logWithLevel(ERROR, message, keyValues...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(message string, keyValues ...interface{}) {
	l.logWithLevel(FATAL, message, keyValues...)
}

// Global logging functions for convenience
func Debug(message string, keyValues ...interface{}) {
	if globalLogger != nil {
		globalLogger.Debug(message, keyValues...)
	}
}

func Info(message string, keyValues ...interface{}) {
	if globalLogger != nil {
		globalLogger.Info(message, keyValues...)
	}
}

func Warning(message string, keyValues ...interface{}) {
	if globalLogger != nil {
		globalLogger.Warning(message, keyValues...)
	}
}

func Error(message string, keyValues ...interface{}) {
	if globalLogger != nil {
		globalLogger.Error(message, keyValues...)
	}
}

func Fatal(message string, keyValues ...interface{}) {
	if globalLogger != nil {
		globalLogger.Fatal(message, keyValues...)
	}
	os.Exit(1)
}

// Close closes the logger and associated resources
func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if logFile != nil {
		err := logFile.Close()
		logFile = nil
		return err
	}

	return nil
}

// SetLevel dynamically changes the log level
func SetLevel(level string) {
	if globalLogger != nil {
		globalLogger.level = parseLogLevel(level)
		globalLogger.Info("Log level changed", "new_level", level)
	}
}

// GetLevel returns the current log level
func GetLevel() string {
	if globalLogger != nil {
		return levelNames[globalLogger.level]
	}
	return "INFO"
}

// IsDebugEnabled returns true if debug logging is enabled
func IsDebugEnabled() bool {
	return globalLogger != nil && globalLogger.level <= DEBUG
}

// WithComponent creates a new logger with a specific component name
func WithComponent(component string) *Logger {
	return NewLogger(component)
}
