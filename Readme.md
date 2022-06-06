# Logger (Golang) Async

# Installation
> go get github.com/lwinmgmg/logger

# Example
```
func main(){
    logger := logger.DefaultLogging(logger.INFO)
    logger.Info("This is info")
    logger.Close()
}
```
Notice : don't forget to close the logger.<br>As it is async logger, you can miss the last messages if you don't close the logger

# Options
You can modify Handlers, Level, Pattern, Pending Count, Time Format using logger.NewLogging(...) function.<br>
Pending Count equal zero means sync logger.<br>
Use logger.NewFileWriter(...) function for rotation file log.
```
// This format and pattern is for default logger
const (
	DEFAULT_TFORMAT string = time.RFC3339
	DEFAULT_PATTERN string = "$(other)$(time) | $(level) | $(message)"
)

// file handler LoggerWriter
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

type Logging struct {
	Handlers        []LoggerWriter
	HandlerChannels []chan Message
	Done            chan struct{}
	Confirm         chan struct{}
	LEVEL           LogLevel
	Pattern         string
	TFormat         string
}
```

