package utils

import (
	"errors"
	"net"
	"strings"
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
