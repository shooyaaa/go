package log

import "fmt"

func Debug(msg string) {
	fmt.Println(msg)
}

func DebugF(format string, values ...interface{}) {
	fmt.Println(format, values)
}
