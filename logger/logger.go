package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
)

const (
	prefixInfo     = "INFO: "
	prefixError    = "ERROR: "
	prefixReminder = "REMINDER: "
)

var mu sync.Mutex

type Logger struct {
	infoLogger     *log.Logger
	errorLogger    *log.Logger
	reminderLogger *log.Logger
	file           *os.File
}

func LoggerNew(filename string) (*Logger, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл %s: %w", filename, err)
	}

	return &Logger{
		infoLogger:     log.New(file, prefixInfo, log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger:    log.New(file, prefixError, log.Ldate|log.Ltime|log.Lshortfile),
		reminderLogger: log.New(file, prefixReminder, log.Ldate|log.Ltime|log.Lshortfile),
		file:           file,
	}, nil

}
func (l *Logger) Info(msg string) {
	mu.Lock()
	defer mu.Unlock()
	l.infoLogger.Output(2, msg)
}
func (l *Logger) Error(msg string) {
	mu.Lock()
	defer mu.Unlock()
	l.errorLogger.Output(2, msg)
}
func (l *Logger) Reminder(msg string) {
	mu.Lock()
	defer mu.Unlock()
	l.reminderLogger.Output(2, msg)
}
func (l *Logger) Close() error {
	err := l.file.Close()
	if err != nil {
		return fmt.Errorf("не удалось закрыть файл %s: %w", l.file.Name(), err)
	}
	return nil
}
