// Filename: test2/internal/jsonlog/jsonlog.go
package jsonlog

import (
	"encoding/json"
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

// logging levels
type Level int8

// Levels start at zero
const (
	LevelInfo  Level = iota //0
	LevelError              //1
	LevelFatal              //2
	LevelOff                //3
)

// The levels presented for humans to read
func (l Level) String() string {
	switch l {
	case LevelInfo:
		return "INFO"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

// Custom Logger
type Logger struct {
	out      io.Writer
	minLevel Level
	mu       sync.Mutex
}

// creats a new instance of logger
func New(out io.Writer, minLevel Level) *Logger {
	return &Logger{
		out:      out,
		minLevel: minLevel,
	}
}

// Helper Methods
func (l *Logger) PrintInfo(message string, properties map[string]string) {
	l.Print(LevelInfo, message, properties)
}

func (l *Logger) PrintError(err error, properties map[string]string) {
	l.Print(LevelError, err.Error(), properties)
}

func (l *Logger) PrintFatal(err error, properties map[string]string) {
	l.Print(LevelFatal, err.Error(), properties)
	os.Exit(1)
}

func (l *Logger) Print(level Level, message string, properties map[string]string) (int, error) {
	//Ensure serverity level is at least minimal
	if level < l.minLevel {
		return 0, nil
	}

	//struct for holding the log entry data
	data := struct {
		Level         string            `json:"level"`
		Time          string            `json:"time"`
		Message       string            `json:"message"`
		Proproperties map[string]string `json:"properties,omitempty"`
		Trace         string            `json:"trace,omitempty"`
	}{
		Level:         level.String(),
		Time:          time.Now().UTC().Format(time.RFC3339),
		Message:       message,
		Proproperties: properties,
	}

	//stack trace?
	if level >= LevelError {
		data.Trace = string(debug.Stack())
	}

	//Encoding the log entry to JSON
	var line []byte
	line, err := json.Marshal(data)
	if err != nil {
		line = []byte(LevelError.String() + ":unable to marshal log message" + err.Error())
	}

	//prepare to write the log entry
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.out.Write(append(line, '\n'))
}

// Implementing the io.write interface
func (l *Logger) Write(message []byte) (n int, err error) {
	return l.Print(LevelError, string(message), nil)
}
