package logger

import (
	"fmt"
	"strings"
	"time"
)

type Message struct {
	Level LogLevel
	Mesg  string
	Args  []any
	Time  time.Time
	Other string
}

func NewMessage(level LogLevel, other, mesg string, args ...any) Message {
	return Message{
		Level: level,
		Mesg:  mesg,
		Args:  args,
		Time:  time.Now(),
		Other: other,
	}
}

func (msg *Message) Format(pattern string, timeFormat string) string {
	pattern = strings.Replace(pattern, "$(other)", "", -1)
	pattern = strings.Replace(pattern, "$(time)", msg.Time.Format(timeFormat), -1)
	pattern = strings.Replace(pattern, "$(message)", fmt.Sprintf(msg.Mesg, msg.Args...), -1)
	pattern = strings.Replace(pattern, "$(level)", LOG_CODE_MAP[msg.Level], -1)
	return pattern + "\n"
}
