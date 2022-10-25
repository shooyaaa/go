package log

import (
	"fmt"
	"time"
)

func formatMsg(msg string) string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05.000 ") + msg + "\n"
}

func Debug(msg string) {
	fmt.Print(formatMsg(msg))
}
func DebugF(format string, values ...interface{}) {
	fmt.Printf(formatMsg(format), values...)
}

func Warn(format string, values ...interface{}) {
	fmt.Printf(formatMsg(format), values...)
}
func WarnF(format string, values ...interface{}) {
	fmt.Printf(formatMsg(format), values...)
}

func Fatal(format string, values ...interface{}) {
	fmt.Printf(formatMsg(format), values)
}
func FatalF(format string, values ...interface{}) {
	fmt.Printf(formatMsg(format), values...)
}

func Info(format string, values ...interface{}) {
	fmt.Printf(formatMsg(format), values...)
}
func InfoF(format string, values ...interface{}) {
	fmt.Printf(formatMsg(format), values...)
}

func Error(format string) {
	fmt.Println(format)
}
func ErrorF(format string, values ...interface{}) {
	fmt.Printf(format, values...)
}
