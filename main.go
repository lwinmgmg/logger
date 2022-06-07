package logger

import (
	"os"
	"time"
)

const (
	DEFAULT_TFORMAT string = time.RFC3339
	DEFAULT_PATTERN string = "$(other)$(time) | $(level) | $(message)"
)

type Logger interface {
	Debug(string, ...any)
	Info(string, ...any)
	Warning(string, ...any)
	Error(string, ...any)
	Critical(string, ...any)
	Fatal(string, ...any)
}

type Logging struct {
	Handlers        []LoggerWriter
	HandlerChannels []chan Message
	Done            chan struct{}
	Confirm         chan struct{}
	Level           LogLevel
	Pattern         string
	TFormat         string
	Other           func() string
	PendingCount    int
}

func NewLogging(level LogLevel, other func() string, pendingCount int, pattern, tFormat string, handlers ...LoggerWriter) *Logging {
	if pendingCount < 0 {
		pendingCount = 1000
	}
	lgr := &Logging{
		Handlers:        handlers,
		HandlerChannels: make([]chan Message, 0, len(handlers)),
		Done:            make(chan struct{}),
		Confirm:         make(chan struct{}, len(handlers)),
		Pattern:         pattern,
		TFormat:         tFormat,
		Other:           other,
		PendingCount:    pendingCount,
	}
	for i, wr := range lgr.Handlers {
		lgr.HandlerChannels = append(lgr.HandlerChannels, make(chan Message, lgr.PendingCount))
		go HandleFunc(lgr.HandlerChannels[i], lgr.Confirm, lgr.Done, wr, lgr.Pattern, lgr.TFormat)
	}
	return lgr
}

func DefaultLogging(level LogLevel, handlers ...LoggerWriter) *Logging {
	lgr := &Logging{
		Handlers: []LoggerWriter{NewConsoleWriter()},
	}
	if len(handlers) > 0 {
		lgr.Handlers = append(lgr.Handlers, handlers...)
	}
	lgr.HandlerChannels = make([]chan Message, 0, len(lgr.Handlers))
	lgr.Done = make(chan struct{})
	lgr.Confirm = make(chan struct{}, len(lgr.Handlers))
	lgr.Pattern = DEFAULT_PATTERN
	lgr.TFormat = DEFAULT_TFORMAT
	lgr.Level = level
	lgr.Other = func() string { return "" }
	lgr.PendingCount = 1000
	for i, wr := range lgr.Handlers {
		lgr.HandlerChannels = append(lgr.HandlerChannels, make(chan Message, lgr.PendingCount))
		go HandleFunc(lgr.HandlerChannels[i], lgr.Confirm, lgr.Done, wr, lgr.Pattern, lgr.TFormat)
	}
	return lgr
}

func (lgr *Logging) write(level LogLevel, mesg string, args ...any) {
	if lgr.Level <= level {
		mesgObj := NewMessage(level, lgr.Other(), mesg, args...)
		for _, ch := range lgr.HandlerChannels {
			ch <- mesgObj
		}
	}
}

func (lgr *Logging) Close() {
	close(lgr.Done)
	for i := 0; i < len(lgr.HandlerChannels); i++ {
		<-lgr.Confirm
	}
}

func (lgr *Logging) Debug(mesg string, args ...any) {
	lgr.write(DEBUG, mesg, args...)
}

func (lgr *Logging) Info(mesg string, args ...any) {
	lgr.write(INFO, mesg, args...)
}

func (lgr *Logging) Warning(mesg string, args ...any) {
	lgr.write(WARNING, mesg, args...)
}

func (lgr *Logging) Error(mesg string, args ...any) {
	lgr.write(ERROR, mesg, args...)
}

func (lgr *Logging) Critical(mesg string, args ...any) {
	lgr.write(CRITICAL, mesg, args...)
}

func (lgr *Logging) Fatal(mesg string, args ...any) {
	lgr.write(ERROR, mesg, args...)
	lgr.Close()
	os.Exit(1)
}
