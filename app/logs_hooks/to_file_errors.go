package logs_hooks

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

var (
	errors_logs *os.File
)

func initFile(rootDir string) {
	logFile := rootDir + "/logs/errors.log"
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create error logfile" + logFile)
		panic(err)
	}
	//defer f.Close()
	errors_logs = f
}

type ToFileErrorHook struct{}

func NewToFileErrorHook(rootDir string) *ToFileErrorHook {
	initFile(rootDir)
	return &ToFileErrorHook{}
}

func (hook *ToFileErrorHook) Fire(entry *log.Entry) error {

	msg := fmt.Sprintf("[%s](%s) %s:%d %s\n",
		strings.ToUpper(entry.Level.String()),
		entry.Time.Format("2006-01-02 15:04:05"),
		entry.Caller.File,
		entry.Caller.Line,
		entry.Message)

	_, err := errors_logs.WriteString(msg)
	if err != nil {
		fmt.Printf("Cannot write error to log file. %s", err.Error())
		return err
	}
	return nil
}

func (hook *ToFileErrorHook) Levels() []log.Level {
	return []log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
	}
}
