package system

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type Info struct {
	Hostname string
	CPU      int64
	Memory   int64
	Arch     string
}

type Collector struct{}

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) Collect() (Info, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return Info{}, fmt.Errorf("failed to read hostname: %w", err)
	}

	memory, err := totalMemoryBytes()
	if err != nil {
		return Info{}, fmt.Errorf("failed to read memory: %w", err)
	}

	return Info{
		Hostname: hostname,
		CPU:      int64(runtime.NumCPU()),
		Memory:   memory,
		Arch:     runtime.GOARCH,
	}, nil
}

func totalMemoryBytes() (int64, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "MemTotal:") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			return 0, fmt.Errorf("invalid MemTotal format")
		}

		kb, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, err
		}
		return kb * 1024, nil
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return 0, fmt.Errorf("MemTotal not found")
}

func FirstPrivateIPv4() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range interfaces {
		if (iface.Flags&net.FlagUp) == 0 || (iface.Flags&net.FlagLoopback) != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			netAddr, ok := addr.(*net.IPNet)
			if !ok || netAddr.IP == nil {
				continue
			}
			ip := netAddr.IP.To4()
			if ip != nil && ip.IsPrivate() {
				return ip.String()
			}
		}
	}

	return ""
}
