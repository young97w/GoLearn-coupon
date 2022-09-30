package util

import (
	"account/log"
	"net"
)

func GenRandomPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "0")
	if err != nil {
		log.Logger.Error(err.Error())
		return 0, err
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Logger.Error(err.Error())
		return 0, err
	}
	port := listener.Addr().(*net.TCPAddr).Port
	return port, nil
}

func GenRandomPort2() int {
	addr, err := net.ResolveTCPAddr("tcp", ":0")
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	return port
}
