package logs_hooks

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

var (
	all_logs *os.File
)

func initFileAll(rootDir string) {
	logFile := rootDir + "/logs/all.log"
	f, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create all logfile" + logFile)
		panic(err)
	}
	//defer f.Close()
	all_logs = f
}

type ToFileAllHook struct{}

func NewToFileHook(rootDir string) *ToFileAllHook {
	initFileAll(rootDir)
	return &ToFileAllHook{}
}

func (hook *ToFileAllHook) Fire(entry *log.Entry) error {

	msg := fmt.Sprintf("[%s](%s) %s:%d %s\n",
		strings.ToUpper(entry.Level.String()),
		entry.Time.Format("2006-01-02 15:04:05"),
		entry.Caller.File,
		entry.Caller.Line,
		entry.Message)

	_, err := all_logs.WriteString(msg)
	if err != nil {
		fmt.Printf("Cannot write to all log file. %s", err.Error())
		return err
	}
	return nil
}

func (hook *ToFileAllHook) Levels() []log.Level {
	return []log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
		log.WarnLevel,
		log.InfoLevel,
		log.DebugLevel,
	}
}
