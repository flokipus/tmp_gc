package main

import (
	"fmt"
	"log"
	"net"
	"net/netip"
	"syscall"
)

func main() {
	// 1. Новый экземпляр получает от старого дескрипторы соединений
	// 2. Посылает ему сигнал о штатном завершении
	// 3. Получает стейт соединений (их буффер)

	// Если использовать SO_REUSEPORT, можно не передавать дескриптор, только его стейт

	tcpLn, err := net.ListenTCP(
		"tcp",
		net.TCPAddrFromAddrPort(netip.AddrPortFrom(netip.IPv4Unspecified(), 6000)),
	)
	if err != nil {
		log.Fatal("listener error: ", err)
	}
	file, err := tcpLn.File()
	if err != nil {
		log.Fatal("converting listener to file descr error: ", err)
	}
	fmt.Println(int(file.Fd()), file.Name())

	var unixAddr *net.UnixAddr = nil
	unixLn, err := net.ListenUnix("unix", unixAddr)
	for {
		worker, err := unixLn.Accept()
		if err != nil {
			log.Fatal("unix accept error: ", err)
		}
		worker.(*net.UnixConn).WriteMsgUnix(
			[]byte("listener"),
			syscall.UnixRights(int(file.Fd())),
			unixAddr,
		)
	}
}
