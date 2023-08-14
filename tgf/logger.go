package tgf

import (
	"encoding/json"
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
	LogUpdate(tgbotapi.Update)
}

type DefaultLogger struct {
	infoLogger    *log.Logger
	warningLogger *log.Logger
	debugLogger   *log.Logger
	errorLogger   *log.Logger
	update        tgbotapi.Update
	debug         bool
}

func NewDefaultLogger(debug bool) *DefaultLogger {
	flags := log.Ldate | log.Ltime | log.LstdFlags

	infoLogger := log.New(os.Stdout, "", flags)
	warningLogger := log.New(os.Stdout, "", flags)
	debugLogger := log.New(os.Stdout, "", flags)
	errorLogger := log.New(os.Stderr, "", flags)

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
	l.infoLogger.Printf(prefix+"[INFO]:\t\t"+format, v...)
}

func (l *DefaultLogger) Error(format string, v ...any) {
	prefix := l.getFileName()
	l.errorLogger.Printf(prefix+"[ERROR]:\t"+format, v...)
}

func (l *DefaultLogger) Warn(format string, v ...any) {
	prefix := l.getFileName()
	l.warningLogger.Printf(prefix+"[WARN]:\t"+format, v...)
}

func (l *DefaultLogger) Debug(format string, v ...any) {
	if !l.debug {
		return
	}
	prefix := l.getFileName()
	l.debugLogger.Printf(prefix+"[DEBUG]:\t"+format, v...)
}

func (l *DefaultLogger) LogUpdate(u tgbotapi.Update) {
	prefix := l.getFileName()

	message := u.Message
	if message == nil {
		message = u.CallbackQuery.Message
	}

	updateID := u.UpdateID
	chatID := message.Chat.ID
	userID := message.From.ID

	messageJSON, _ := json.Marshal(u.Message)
	callbackQueryJSON, _ := json.Marshal(u.CallbackQuery)

	l.infoLogger.Printf(prefix+"[UPDATE]:\tupdateID: %d, chatID: %d, userID: %d, messageJSON: %s, callbackQueryJSON: %s",
		updateID, chatID, userID, messageJSON, callbackQueryJSON)
}
