package exceptions

import (
	"chatsystem/internal/models"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Logger represents a logging interface.
type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})
}

// FileLogger implements the Logger interface using a file.
type FileLogger struct {
	file *os.File
}

// NewFileLogger creates a new FileLogger instance.
func NewFileLogger(filename string) (*FileLogger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &FileLogger{file: file}, nil
}

// Debug logs a debug message.
func (l *FileLogger) Debug(v ...interface{}) {
	l.log("DEBUG", v...)
}

// Info logs an info message.
func (l *FileLogger) Info(v ...interface{}) {
	l.log("INFO", v...)
}

// Warn logs a warning message.
func (l *FileLogger) Warn(v ...interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	scriptName := filepath.Base(funcName)
	l.log("WARNING", fmt.Sprintf("%s - %s - %s", funcName, scriptName, v))
}

// Error logs an error message.
func (l *FileLogger) Error(v ...interface{}) {
	l.log("ERROR", v...)
}

// Fatal logs a fatal message and exits the program.
func (l *FileLogger) Fatal(v ...interface{}) {
	l.log("FATAL", v...)
	os.Exit(1)
}

// log logs a message with the specified level.
func (l *FileLogger) log(level string, v ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05.999999-0700 MST")
	message := fmt.Sprintf("%s - %s - %s\n", timestamp, level, fmt.Sprint(v...))
	if _, err := l.file.WriteString(message); err != nil {
		log.Fatal(err)
	}
}

// Close closes the underlying file.
func (l *FileLogger) Close() error {
	return l.file.Close()
}

type CustomLoggers struct {
	System  *FileLogger
	UserOps *FileLogger
	ToolOps *FileLogger
	AiOps   *FileLogger
	RagOps  *FileLogger
}

var Loggers CustomLoggers

// findRootDir traverses up the directory tree to find the project root
// It looks for a specific marker like go.mod or a specific directory name
func findRootDir() (string, error) {
	curr, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Keep going up until we find the project root marker (e.g., go.mod)
	for {
		// Check for go.mod or any other project root indicator
		if _, err := os.Stat(filepath.Join(curr, "go.mod")); err == nil {
			return curr, nil
		}

		// Get the parent directory
		parent := filepath.Dir(curr)
		if parent == curr {
			// We've reached the root of the filesystem without finding the project root
			return "", fmt.Errorf("project root not found")
		}
		curr = parent
	}
}

func init() {
	// Find the project root directory
	rootDir, err := findRootDir()
	if err != nil {
		log.Fatal(err)
	}

	// Create the logs folder in the root directory
	logsFolder := filepath.Join(rootDir, "logs")
	if err := os.MkdirAll(logsFolder, 0755); err != nil {
		log.Fatal(err)
	}

	// Create log files in the root directory
	logFiles := []string{"system.log", "aiops.log", "int_aiops.log", "toolops.log", "usersops.log", "ragops.log"}
	for _, filename := range logFiles {
		logFile := filepath.Join(logsFolder, filename)
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			if err := os.WriteFile(logFile, []byte(""), 0644); err != nil {
				log.Fatal(err)
			}
		}
	}

	// Setup loggers with absolute paths
	systemLogger, err := NewFileLogger(filepath.Join(logsFolder, "system.log"))
	if err != nil {
		log.Fatal(err)
	}

	aipipelineLogger, err := NewFileLogger(filepath.Join(logsFolder, "aiops.log"))
	if err != nil {
		log.Fatal(err)
	}

	toolLogger, err := NewFileLogger(filepath.Join(logsFolder, "toolops.log"))
	if err != nil {
		log.Fatal(err)
	}

	usersopsLogger, err := NewFileLogger(filepath.Join(logsFolder, "usersops.log"))
	if err != nil {
		log.Fatal(err)
	}

	ragopsLogger, err := NewFileLogger(filepath.Join(logsFolder, "ragops.log"))
	if err != nil {
		log.Fatal(err)
	}

	Loggers.System = systemLogger
	Loggers.ToolOps = toolLogger
	Loggers.UserOps = usersopsLogger
	Loggers.AiOps = aipipelineLogger
	Loggers.RagOps = ragopsLogger
}

func IOLogger(rc int, detail, ext_ref string) models.Error {
	var error models.Error
	error.ResponseCode = rc
	error.Message = "Failed"
	error.Detail = detail
	error.ExternalReference = ext_ref

	return error
}
