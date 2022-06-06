package logger_test

import (
	"sync"
	"testing"

	"github.com/lwinmgmg/logger"
)

var wg = sync.WaitGroup{}

func TestHandleFunc(t *testing.T) {
	mesgCh := []chan logger.Message{make(chan logger.Message), make(chan logger.Message)}
	confirmCh := make(chan struct{})
	confirmValue := 0
	doneCh := make(chan struct{})
	for _, v := range mesgCh {
		go logger.HandleFunc(v, confirmCh, doneCh, &logger.ConsoleWriter{}, "abcd", "")
	}
	wg.Add(1)
	go func(ch <-chan struct{}) {
		<-ch
		confirmValue += 1
		wg.Done()
	}(confirmCh)
	if confirmValue != 0 {
		t.Errorf("Expected confirm value : %v, Getting value : %v", 0, confirmValue)
	}
	close(doneCh)
	wg.Wait()
	if confirmValue == 0 {
		t.Errorf("Expected confirm value greater than zero, Getting value : %v", confirmValue)
	}
}
