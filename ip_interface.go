package main

import (
	"fmt"
	"net"
)

func main() {
	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error fetching interfaces:", err)
		return
	}

	for _, iface := range interfaces {
		// 过滤掉禁用的接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// 获取接口的地址
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println("Error fetching addresses:", err)
			continue
		}

		for _, addr := range addrs {
			// 只处理 IPv4 地址
			if ipNet, ok := addr.(*net.IPNet); ok && ipNet.IP.To4() != nil {
				fmt.Printf("Interface: %s, IP: %s\n", iface.Name, ipNet.IP.String())
			}
		}
	}
}
