package logger_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/lwinmgmg/logger"
)

func TestConsoleWriter(t *testing.T) {
	mesg := "[   INFO   ]My Name is Lwin Maung Maung"
	consoleWriter := logger.ConsoleWriter{}
	length, err := consoleWriter.Write([]byte(mesg))
	if err != nil {
		t.Errorf("Getting error on writing console : %v", err)
	}
	if length != len(mesg) {
		t.Errorf("Expect lenght of the message : %v, Getting : %v", len(mesg), length)
	}
}

func TestFileWriterInit(t *testing.T) {
	filename := path.Join(os.TempDir(), "testing_logging.log")
	fileWriter := logger.NewFileWriter(filename, 0, 0, nil)

	if _, err := os.Open(filename); err != nil {
		t.Errorf("Getting error on checking log file : %v", err)
	}
	if fileWriter.MaxSize == 0 {
		t.Errorf("Setting default max size does not work")
	}
	if err := os.Remove(filename); err != nil {
		t.Errorf("Getting error on remove tmp log file : %v", err)
	}
}

func TestFileWriterWrite(t *testing.T) {
	filename := path.Join(os.TempDir(), "testing_logging.log")
	fileWriter := logger.NewFileWriter(filename, 1, 5, nil)
	mesg := "My Name is Lwin Maung Maung\n"
	wrote, err := fileWriter.Write([]byte(mesg))
	if err != nil {
		t.Errorf("Getting error on write function : %v", err)
	}
	if wrote != len(mesg) {
		t.Errorf("Expect lenght of the message : %v, Getting : %v", len(mesg), wrote)
	}
	wrote, err = fileWriter.Write([]byte(mesg))
	if err != nil {
		t.Errorf("Getting error on second write function : %v", err)
	}
	if err := os.Remove(filename); err != nil {
		t.Errorf("Getting error on remove file [%v] : %v", filename, err)
	}
	if err := os.Remove(fmt.Sprintf("%v%v", filename, 1)); err != nil {
		t.Errorf("Getting error on remove file [%v] : %v", filename+"1", err)
	}
}

func TestFileWriterClose(t *testing.T) {
	filename := path.Join(os.TempDir(), "testing_logging.log")
	fileWriter := logger.NewFileWriter(filename, 0, 0, nil)
	if err := fileWriter.Close(); err != nil {
		t.Errorf("Getting error on file writer close : %v", err)
	}
	if _, err := fileWriter.Write([]byte("LMM")); err == nil {
		t.Error("Expecting error on write after close")
	}
}
