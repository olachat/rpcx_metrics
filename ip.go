package prom

import (
	"errors"
	"fmt"
	"net"
)

var (
	// ErrorEmptyInterfaceAddrs 定义错误，无法找到网卡信息时
	ErrorEmptyInterfaceAddrs = fmt.Errorf("empty found in InterfaceAddrs")
)

// IP 单例ip，并且导出
var IP = &ip{}

type ip struct {
}

// LocalIPv4s 获取本机局域网ipv4地址列表
func (*ip) LocalIPv4s() (ips []string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}
	if len(ips) == 0 {
		return ips, ErrorEmptyInterfaceAddrs
	}
	return ips, nil
}

// LocalIPv4s 获取本机局域网ipv4地址
func (*ip) LocalIPv4() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String(), nil
		}
	}

	return "", ErrorEmptyInterfaceAddrs
}

// IsIPV4 判断字符串是否是ipv4
func (*ip) IsIPV4(ipv4 string) bool {
	address := net.ParseIP(ipv4)
	return address != nil
}

// IsLanIP 判断字符串是否是内网IP
func (*ip) IsLanIP(ipv4 string) bool {
	ip := net.ParseIP(ipv4)
	if ip == nil {
		return false
	}
	if ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		return true
	}
	if ip4 := ip.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return true
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return true
		case ip4[0] == 192 && ip4[1] == 168:
			return true
		default:
			return false
		}
	}
	return false
}

func (*ip) Long2ip(ipInt int64) (string, error) {
	if ipInt < 0 || ipInt > 4294967295 {
		return "", errors.New("invalid ip address")
	}

	ip := make([]byte, 4)
	for i := 0; i < 4; i++ {
		ip[i] = byte(ipInt >> uint(i*8) & 0xFF)
	}

	return fmt.Sprintf("%d.%d.%d.%d", ip[3], ip[2], ip[1], ip[0]), nil
}

// GetAvailablePort 获取一个可用端口
func GetAvailablePort() int {
	// Listen on a random TCP port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = listener.Close()
	}()

	// Extract the actual port number assigned
	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port
}
