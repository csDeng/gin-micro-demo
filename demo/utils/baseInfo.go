package utils

import (
	"errors"
	"gin-micro-demo/config"
	"net"
	"strings"
	"sync"
)

func GetIp() (string, error) {
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	s := conn.LocalAddr().String()
	arr := strings.Split(s, ":")
	if len(arr) < 2 {
		return "", errors.New("获取 ip 失败")
	}
	return arr[0], nil
}

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

var (
	ip_port *config.IpPort
	ip_once sync.Once
)

func GetIpPort() *config.IpPort {
	ip_once.Do(func() {
		p, err := GetFreePort()
		if err != nil {
			panic(err)
		}
		ip, err := GetIp()
		if err != nil {
			panic(err)
		}
		ip_port = &config.IpPort{
			Port: p,
			Ip:   ip,
		}
	})
	return ip_port
}
