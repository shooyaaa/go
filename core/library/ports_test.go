package library

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCurrentProcessListeningPorts(t *testing.T) {
	ports, err := GetCurrentProcessListeningPorts()

	// 即使没有监听端口，也不应该返回错误
	assert.NoError(t, err)
	// 注意：ports 可能是空列表（如果没有监听端口），这是正常的
	// 在 Go 中，空切片不是 nil，所以这里只检查 err
	_ = ports // 使用 ports 避免未使用变量警告
}

func TestParseAddress(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantPort  int
		wantAddr  string
		wantError bool
	}{
		{
			name:      "wildcard IPv4",
			input:     "*:8080",
			wantPort:  8080,
			wantAddr:  "0.0.0.0",
			wantError: false,
		},
		{
			name:      "localhost IPv4",
			input:     "127.0.0.1:8080",
			wantPort:  8080,
			wantAddr:  "127.0.0.1",
			wantError: false,
		},
		{
			name:      "IPv6",
			input:     "[::1]:8080",
			wantPort:  8080,
			wantAddr:  "[::1]",
			wantError: false,
		},
		{
			name:      "with LISTEN suffix",
			input:     "*:8080 (LISTEN)",
			wantPort:  8080,
			wantAddr:  "0.0.0.0",
			wantError: false,
		},
		{
			name:      "invalid format",
			input:     "invalid",
			wantPort:  0,
			wantAddr:  "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			port, addr, err := parseAddress(tt.input)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantPort, port)
				assert.Equal(t, tt.wantAddr, addr)
			}
		})
	}
}

func TestGetListeningPortsByProtocol(t *testing.T) {
	ports := []ListeningPort{
		{Port: 8080, Protocol: "tcp", Address: "0.0.0.0"},
		{Port: 8081, Protocol: "udp", Address: "0.0.0.0"},
		{Port: 8082, Protocol: "tcp6", Address: "[::1]"},
		{Port: 8083, Protocol: "tcp", Address: "127.0.0.1"},
	}

	tcpPorts := GetListeningPortsByProtocol(ports, "tcp")
	assert.Len(t, tcpPorts, 3) // tcp, tcp6, tcp

	udpPorts := GetListeningPortsByProtocol(ports, "udp")
	assert.Len(t, udpPorts, 1) // udp

	tcp6Ports := GetListeningPortsByProtocol(ports, "tcp6")
	assert.Len(t, tcp6Ports, 1) // tcp6
}

func TestFormatListeningPorts(t *testing.T) {
	ports := []ListeningPort{
		{Port: 8080, Protocol: "tcp", Address: "0.0.0.0"},
		{Port: 8081, Protocol: "udp", Address: "127.0.0.1"},
	}

	formatted := FormatListeningPorts(ports)
	assert.Contains(t, formatted, "Found 2 listening port(s)")
	assert.Contains(t, formatted, "tcp://0.0.0.0:8080")
	assert.Contains(t, formatted, "udp://127.0.0.1:8081")

	empty := FormatListeningPorts([]ListeningPort{})
	assert.Contains(t, empty, "No listening ports found")
}
