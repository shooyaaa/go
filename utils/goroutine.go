package utils

import (
	"bytes"
	"runtime"
	"strconv"
	"strings"
)

// GetGoroutineID 获取当前 goroutine 的 ID
// 通过解析 runtime.Stack() 的输出来获取 goroutine ID
func GetGoroutineID() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	// runtime.Stack 的输出格式为: "goroutine 123 [running]:\n..."
	// 我们需要提取 "123" 这个数字
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.ParseInt(idField, 10, 64)
	if err != nil {
		return -1
	}
	return id
}

// GetGoroutineIDString 获取当前 goroutine 的 ID（字符串形式）
func GetGoroutineIDString() string {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	// runtime.Stack 的输出格式为: "goroutine 123 [running]:\n..."
	// 我们需要提取 "123" 这个数字
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	return idField
}

// GetGoroutineIDFast 使用更高效的方式获取 goroutine ID
// 这个方法使用 bytes 操作，性能更好
func GetGoroutineIDFast() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)

	// 查找 "goroutine " 的位置
	goroutinePrefix := []byte("goroutine ")
	prefixLen := len(goroutinePrefix)

	// 在 buf 中查找 "goroutine " 前缀
	idx := bytes.Index(buf[:n], goroutinePrefix)
	if idx == -1 {
		return -1
	}

	// 跳过 "goroutine " 前缀，找到 ID 的开始位置
	start := idx + prefixLen

	// 找到 ID 的结束位置（空格或换行符）
	end := start
	for end < n && buf[end] != ' ' && buf[end] != '\n' && buf[end] != '[' {
		end++
	}

	// 解析 ID
	id, err := strconv.ParseInt(string(buf[start:end]), 10, 64)
	if err != nil {
		return -1
	}

	return id
}
