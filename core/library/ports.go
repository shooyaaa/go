package library

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// ListeningPort 表示一个监听的端口信息
type ListeningPort struct {
	Port     int    // 端口号
	Protocol string // 协议类型 (tcp, udp, tcp6, udp6)
	Address  string // 监听地址
}

// GetCurrentProcessListeningPorts 获取当前进程监听的端口和协议类型
func GetCurrentProcessListeningPorts() ([]ListeningPort, error) {
	pid := os.Getpid()
	return GetProcessListeningPorts(pid)
}

// GetProcessListeningPorts 获取指定进程ID监听的端口和协议类型
func GetProcessListeningPorts(pid int) ([]ListeningPort, error) {
	var ports []ListeningPort
	var err error

	switch runtime.GOOS {
	case "darwin", "linux":
		ports, err = getListeningPortsUnix(pid)
	case "windows":
		ports, err = getListeningPortsWindows(pid)
	default:
		return nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	return ports, err
}

// getListeningPortsUnix 在 Unix 系统（macOS/Linux）上获取监听端口
func getListeningPortsUnix(pid int) ([]ListeningPort, error) {
	// 使用 lsof 命令获取监听端口
	cmd := exec.Command("lsof", "-i", "-P", "-n", "-p", strconv.Itoa(pid))
	output, err := cmd.Output()
	if err != nil {
		// lsof 可能没有找到监听端口，返回空列表而不是错误
		return []ListeningPort{}, nil
	}

	return parseLsofOutput(string(output))
}

// parseLsofOutput 解析 lsof 命令的输出
func parseLsofOutput(output string) ([]ListeningPort, error) {
	var ports []ListeningPort
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "COMMAND") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		// lsof 输出格式：
		// COMMAND PID USER FD TYPE DEVICE SIZE/OFF NODE NAME
		// 例如: go 12345 user 3u IPv4 0x... TCP *:8080 (LISTEN)

		// 检查是否是 LISTEN 状态
		if !strings.Contains(line, "LISTEN") {
			continue
		}

		// 获取协议类型（第8个字段，索引7）
		protocol := fields[7]
		if !strings.HasPrefix(protocol, "TCP") && !strings.HasPrefix(protocol, "UDP") {
			continue
		}

		// 获取 NAME 字段（最后一个字段），格式如: *:8080 或 127.0.0.1:8080
		nameField := fields[len(fields)-1]

		// 解析地址和端口
		port, address, err := parseAddress(nameField)
		if err != nil {
			continue
		}

		// 标准化协议名称
		protocolName := strings.ToLower(protocol)
		if strings.Contains(protocolName, "tcp") {
			if strings.Contains(protocolName, "6") || strings.Contains(nameField, "[") {
				protocolName = "tcp6"
			} else {
				protocolName = "tcp"
			}
		} else if strings.Contains(protocolName, "udp") {
			if strings.Contains(protocolName, "6") || strings.Contains(nameField, "[") {
				protocolName = "udp6"
			} else {
				protocolName = "udp"
			}
		}

		ports = append(ports, ListeningPort{
			Port:     port,
			Protocol: protocolName,
			Address:  address,
		})
	}

	return ports, nil
}

// parseAddress 解析地址字符串，提取端口和地址
// 输入格式: *:8080, 127.0.0.1:8080, [::1]:8080
func parseAddress(addr string) (int, string, error) {
	// 移除可能的 (LISTEN) 后缀
	addr = strings.TrimSpace(addr)
	if idx := strings.Index(addr, " ("); idx != -1 {
		addr = addr[:idx]
	}

	// 处理 IPv6 地址 [::1]:8080
	if strings.HasPrefix(addr, "[") {
		idx := strings.LastIndex(addr, "]:")
		if idx == -1 {
			return 0, "", fmt.Errorf("invalid IPv6 address format: %s", addr)
		}
		address := addr[1:idx]  // 移除 [ 和 ]
		portStr := addr[idx+2:] // 移除 ]:
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return 0, "", err
		}
		return port, "[" + address + "]", nil
	}

	// 处理 IPv4 地址 *:8080 或 127.0.0.1:8080
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("invalid address format: %s", addr)
	}

	address := parts[0]
	if address == "*" {
		address = "0.0.0.0"
	}

	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, "", err
	}

	return port, address, nil
}

// getListeningPortsWindows 在 Windows 系统上获取监听端口
func getListeningPortsWindows(pid int) ([]ListeningPort, error) {
	// 使用 netstat 命令
	cmd := exec.Command("netstat", "-ano")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to execute netstat: %w", err)
	}

	return parseNetstatOutput(string(output), pid)
}

// parseNetstatOutput 解析 netstat 命令的输出（Windows）
func parseNetstatOutput(output string, pid int) ([]ListeningPort, error) {
	var ports []ListeningPort
	lines := strings.Split(output, "\n")
	pidStr := strconv.Itoa(pid)

	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "Proto") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		// netstat 输出格式（Windows）:
		// Proto Local Address Foreign Address State PID
		// TCP 0.0.0.0:8080 0.0.0.0:0 LISTENING 12345

		// 检查 PID 是否匹配
		if fields[len(fields)-1] != pidStr {
			continue
		}

		// 检查是否是 LISTENING 状态
		if !strings.Contains(line, "LISTENING") {
			continue
		}

		protocol := strings.ToLower(fields[0])
		localAddr := fields[1]

		// 解析本地地址
		port, address, err := parseAddress(localAddr)
		if err != nil {
			continue
		}

		ports = append(ports, ListeningPort{
			Port:     port,
			Protocol: protocol,
			Address:  address,
		})
	}

	return ports, nil
}

// GetListeningPortsByProtocol 按协议类型过滤监听端口
func GetListeningPortsByProtocol(ports []ListeningPort, protocol string) []ListeningPort {
	var filtered []ListeningPort
	protocol = strings.ToLower(protocol)

	for _, port := range ports {
		if strings.Contains(strings.ToLower(port.Protocol), protocol) {
			filtered = append(filtered, port)
		}
	}

	return filtered
}

// FormatListeningPorts 格式化监听端口信息为字符串
func FormatListeningPorts(ports []ListeningPort) string {
	if len(ports) == 0 {
		return "No listening ports found"
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Found %d listening port(s):\n", len(ports)))
	for i, port := range ports {
		builder.WriteString(fmt.Sprintf("  %d. %s://%s:%d\n", i+1, port.Protocol, port.Address, port.Port))
	}

	return builder.String()
}
