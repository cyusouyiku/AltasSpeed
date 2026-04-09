package main

import (
	"context"//上下文包
	"fmt"
	"log"//一个简单的日志包
	"os"//提供了与操作系统交互的接口，比如读取环境变量，执行系统命令等等
	"os/signal"//专门用来处理系统发送程序的信号，在这个文件中监听os.Interrupt信号（通常是Ctrl+C）来优雅地停止监控
	"sort"
	"time"

	netio "github.com/shirou/gopsutil/v4/net"
)

type ifaceSpeed struct {
	Name      string
	Download  uint64
	Upload    uint64
	TotalRecv uint64
	TotalSent uint64
}

func main() {
	monitor()
}

func monitor() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := runMonitor(ctx, time.Second); err != nil {
		log.Fatalf("network monitor stopped with error: %v", err)
	}
}

func runMonitor(ctx context.Context, interval time.Duration) error {
	if interval <= 0 {
		interval = time.Second
	}

	previous, err := netio.IOCounters(true)
	if err != nil {
		return fmt.Errorf("read initial network stats: %w", err)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	fmt.Println("Network monitor started. Press Ctrl+C to stop.")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("\nNetwork monitor stopped.")
			return nil
		case <-ticker.C:
			current, err := netio.IOCounters(true)
			if err != nil {
				return fmt.Errorf("read current network stats: %w", err)
			}

			renderStats(previous, current, interval)
			previous = current
		}
	}
}

func renderStats(previous, current []netio.IOCountersStat, interval time.Duration) {
	previousByName := make(map[string]netio.IOCountersStat, len(previous))
	for _, stat := range previous {
		previousByName[stat.Name] = stat
	}

	seconds := interval.Seconds()
	if seconds <= 0 {
		seconds = 1
	}

	speeds := make([]ifaceSpeed, 0, len(current))
	for _, stat := range current {
		prev, ok := previousByName[stat.Name]
		if !ok {
			prev = stat
		}

		download := uint64(float64(diffUint64(prev.BytesRecv, stat.BytesRecv)) / seconds)
		upload := uint64(float64(diffUint64(prev.BytesSent, stat.BytesSent)) / seconds)

		if download == 0 && upload == 0 && stat.BytesRecv == 0 && stat.BytesSent == 0 {
			continue
		}

		speeds = append(speeds, ifaceSpeed{
			Name:      stat.Name,
			Download:  download,
			Upload:    upload,
			TotalRecv: stat.BytesRecv,
			TotalSent: stat.BytesSent,
		})
	}

	sort.Slice(speeds, func(i, j int) bool {
		return speeds[i].Download+speeds[i].Upload > speeds[j].Download+speeds[j].Upload
	})

	fmt.Print("\033[H\033[2J")
	fmt.Printf("Time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("Interface            Download          Upload            Total RX          Total TX")
	fmt.Println("--------------------------------------------------------------------------------")

	if len(speeds) == 0 {
		fmt.Println("No active network interfaces detected.")
		return
	}

	for _, stat := range speeds {
		fmt.Printf("%-20s %-17s %-17s %-17s %-17s\n",
			stat.Name,
			formatBytes(stat.Download)+"/s",
			formatBytes(stat.Upload)+"/s",
			formatBytes(stat.TotalRecv),
			formatBytes(stat.TotalSent),
		)
	}
}

func diffUint64(previous, current uint64) uint64 {
	if current < previous {
		return 0
	}
	return current - previous
}

func formatBytes(value uint64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	size := float64(value)
	unit := 0

	for size >= 1024 && unit < len(units)-1 {
		size /= 1024
		unit++
	}

	if unit == 0 {
		return fmt.Sprintf("%d %s", value, units[unit])
	}

	return fmt.Sprintf("%.2f %s", size, units[unit])
}
