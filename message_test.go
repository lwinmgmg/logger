package logger_test

import (
	"fmt"
	"testing"

	"github.com/lwinmgmg/logger"
)

func TestMessageFormat(t *testing.T) {
	mesg := logger.NewMessage(logger.DEBUG, "", "This is debug log")
	output := mesg.Format(logger.DEFAULT_PATTERN, logger.DEFAULT_TFORMAT)
	fmt.Println(output)
	if output == "" {
		t.Errorf("Expecting output message, Getting : %v", output)
	}
}
