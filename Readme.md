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

# Options
You can modify Handlers, Level, Pattern, Time Format using logger.NewLogging(...) function.
```
// This format and pattern is for default logger
const (
	DEFAULT_TFORMAT string = time.RFC3339
	DEFAULT_PATTERN string = "$(other)$(time) | $(level) | $(message)"
)

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

