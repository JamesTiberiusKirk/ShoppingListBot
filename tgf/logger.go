package tgf

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Logger interface {
	Info(string, ...any)
	Error(string, ...any)
	Warn(string, ...any)
	Debug(string, ...any)
	SetUpdate(tgbotapi.Update)
	LogUpdate()
}

type DefaultLogger struct {
	infoLogger    *log.Logger
	warningLogger *log.Logger
	debugLogger   *log.Logger
	errorLogger   *log.Logger
	update        tgbotapi.Update
	debug         bool
}

func (l *DefaultLogger) SetUpdate(update tgbotapi.Update) {
	l.update = update
}

func NewDefaultLogger(debug bool) *DefaultLogger {
	flags := log.Ldate | log.Ltime | log.LstdFlags

	infoLogger := log.New(os.Stdout, "INFO: ", flags)
	warningLogger := log.New(os.Stdout, "WARNING: ", flags)
	debugLogger := log.New(os.Stdout, "WARNING: ", flags)
	errorLogger := log.New(os.Stderr, "ERROR: ", flags)

	return &DefaultLogger{
		infoLogger:    infoLogger,
		warningLogger: warningLogger,
		debugLogger:   debugLogger,
		errorLogger:   errorLogger,
		debug:         debug,
	}
}

func (l *DefaultLogger) getFileName() string {
	_, file, line, _ := runtime.Caller(2)
	fileName := filepath.Base(file)
	return fmt.Sprintf("[%s:%d] ", fileName, line)
}

func (l *DefaultLogger) Info(format string, v ...any) {
	prefix := l.getFileName()
	l.infoLogger.Printf(prefix+format, v...)
}

func (l *DefaultLogger) Error(format string, v ...any) {
	prefix := l.getFileName()
	l.errorLogger.Printf(prefix+format, v...)
}

func (l *DefaultLogger) Warn(format string, v ...any) {
	prefix := l.getFileName()
	l.warningLogger.Printf(prefix+format, v...)
}

func (l *DefaultLogger) Debug(format string, v ...any) {
	if !l.debug {
		return
	}
	prefix := l.getFileName()
	l.debugLogger.Printf(prefix+format, v...)
}

// TODO: Implement
func (l *DefaultLogger) LogUpdate() {
	prefix := l.getFileName()
	l.infoLogger.Printf(prefix + "UPDATE: NEED TO IMPLEMENT")
}
