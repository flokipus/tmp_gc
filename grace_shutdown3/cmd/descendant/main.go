package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"syscall"
)

func main() {
	predcessor, err := net.Dial("unix", "/tmp/predcessor.sock")
	if err != nil {
		//TODO: This is first worker
		log.Fatal("predcessor net.Dial error:", err)
	}

	var (
		msg = make([]byte, 1024)
		oob = make([]byte, 1024)
	)
	msgn, oobn, _, _, err := predcessor.(*net.UnixConn).ReadMsgUnix(msg, oob)

	cmsg, _ := syscall.ParseSocketControlMessage(oob[:oobn])
	fds, _ := syscall.ParseUnixRights(&cmsg[0])
	file := os.NewFile(uintptr(fds[0]), string(msg[:msgn]))

	ln, _ := net.FileListener(file)
	(&http.Server{}).Serve(ln)
}
