package logger

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

type LogLevel uint8

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	CRITICAL
)

var (
	RED          = color.New(color.Bold, color.FgHiRed)
	GREEN        = color.New(color.Bold, color.FgHiGreen)
	BLUE         = color.New(color.Bold, color.FgHiBlue)
	YELLOW       = color.New(color.Bold, color.FgHiYellow)
	LOG_CODE_MAP = map[LogLevel]string{
		DEBUG:    " DEBUG  ",
		INFO:     "  INFO  ",
		WARNING:  "WARNING ",
		ERROR:    " ERROR  ",
		CRITICAL: "CRITICAL",
	}
	ColorFuncMap = map[string]func(string) string{
		LOG_CODE_MAP[DEBUG]:    func(s string) string { return BLUE.Sprint(s) },
		LOG_CODE_MAP[INFO]:     func(s string) string { return GREEN.Sprint(s) },
		LOG_CODE_MAP[WARNING]:  func(s string) string { return YELLOW.Sprint(s) },
		LOG_CODE_MAP[ERROR]:    func(s string) string { return RED.Sprint(s) },
		LOG_CODE_MAP[CRITICAL]: func(s string) string { return RED.Sprint(s) },
	}
)

type LoggerWriter interface {
	Write([]byte) (int, error)
	Close() error
}

type ConsoleWriter struct {
}

func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{}
}

func (coslWr *ConsoleWriter) Write(mesg []byte) (int, error) {
	newMesg := string(mesg)
	length := len(newMesg)
	for str, fun := range ColorFuncMap {
		newMesg = strings.Replace(newMesg, str, fun(str), -1)
	}
	_, err := fmt.Print(newMesg)
	return length, err
}

func (coslWr *ConsoleWriter) Close() error {
	return nil
}

type FileWriter struct {
	//log filename, should be full path
	FileName string
	//log file rotation
	RotationCount int
	MaxSize       int64
	//you can provide writercloser but please do not assign filename
	WriterCloser io.WriteCloser

	//private variables
	currentSize int64
}

func NewFileWriter(filename string, rotateCount int, maxSize int64, wr io.WriteCloser) *FileWriter {
	fileWr := &FileWriter{
		FileName:      filename,
		RotationCount: rotateCount,
		MaxSize:       maxSize,
		WriterCloser:  wr,
	}
	if fileWr.FileName == "" && fileWr.WriterCloser == nil {
		panic(map[string]any{
			"status": 500,
			"mesg":   "FileName or WriterCloser does not exist!",
		})
	} else if fileWr.WriterCloser != nil && fileWr.FileName == "" {
	} else {
		file, err := os.OpenFile(fileWr.FileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			panic(map[string]any{
				"status": 500,
				"mesg":   fmt.Sprintf("Error on opening file : %v", err),
			})
		}
		fileWr.WriterCloser = file
	}
	if file, ok := fileWr.WriterCloser.(*os.File); ok {
		fInfo, err := file.Stat()
		if err != nil {
			panic(map[string]any{
				"status": 500,
				"mesg":   fmt.Sprintf("Error on checking file : %v", err),
			})
		}
		fileWr.currentSize = fInfo.Size()
		fileWr.FileName = file.Name()
		if fileWr.MaxSize == 0 {
			fileWr.MaxSize = 100 * 1024 * 1024 //100 MB default
		}
	} else {
		fileWr.MaxSize = 0
	}
	return fileWr
}

func (fileWr *FileWriter) rotateFile() error {
	for i := fileWr.RotationCount; i >= 0; i-- {
		if i == fileWr.RotationCount {
			filename := fileWr.FileName + strconv.Itoa(i)
			if CheckFileExist(filename) {
				if err := os.Remove(filename); err != nil {
					return err
				}
			}
		} else if i == 0 {
			if err := os.Rename(fileWr.FileName, fileWr.FileName+strconv.Itoa(i+1)); err != nil {
				return err
			}
		} else {
			filename := fileWr.FileName + strconv.Itoa(i)
			if CheckFileExist(filename) {
				if err := os.Rename(filename, fileWr.FileName+strconv.Itoa(i+1)); err != nil {
					return err
				}
			}
		}
	}
	file, err := os.OpenFile(fileWr.FileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	fileWr.WriterCloser = file
	fileWr.currentSize = 0
	return nil
}

func (fileWr *FileWriter) Write(mesg []byte) (int, error) {
	if fileWr.MaxSize > 0 && fileWr.currentSize >= fileWr.MaxSize {
		if err := fileWr.rotateFile(); err != nil {
			return 0, err
		}
	}
	wrote, err := fileWr.WriterCloser.Write(mesg)
	fileWr.currentSize += int64(wrote)
	return wrote, err
}

func (fileWr *FileWriter) Close() error {
	return fileWr.WriterCloser.Close()
}

func CheckFileExist(fname string) bool {
	if _, err := os.Stat(fname); err != nil {
		return false
	}
	return true
}
