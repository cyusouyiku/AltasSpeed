package main

import (
	"fmt"
	"log"

	netio "github.com/shirou/gopsutil/v4/net"
)

func main() {
	stats, err := netio.IOCounters(true)
	if err != nil {
		log.Fatalf("failed to read network stats: %v", err)
	}

	if len(stats) == 0 {
		fmt.Println("No network interfaces found.")
		return
	}

	for _, iface := range stats {
		if iface.BytesRecv == 0 && iface.BytesSent == 0 {
			continue
		}

		fmt.Printf("%s -> RX: %d bytes, TX: %d bytes\n", iface.Name, iface.BytesRecv, iface.BytesSent)
	}
}

/*
以下是Linux系统中读取网络接口流量的实例代码，通过读取/proc/net/dev文件来获取网络接口的流量信息。请注意，这段代码仅适用于Linux系统，并且需要具有适当的权限来访问该文件。

package main

import (
    "fmt"
    "os"
    "strings"
)

func main() {
    // 读取文件
    data, _ := os.ReadFile("/proc/net/dev")
    
    // 按行拆分
    lines := strings.Split(string(data), "\n")
    
    // 遍历每一行
    for _, line := range lines {
        // 1. 检查这一行是否包含 "eth0"
        if strings.Contains(line, "eth0:") {
            
            // 2. 按空白拆分这一行
            fields := strings.Fields(line)
            
            // 3. 取出需要的字段
            rx := fields[1]  // 接收字节
            tx := fields[9]  // 发送字节
            
            fmt.Printf("eth0 接收: %s 字节\n", rx)
            fmt.Printf("eth0 发送: %s 字节\n", tx)
        }
    }
}*/