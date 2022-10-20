package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
)

func main() {
	ln, err := net.Listen("tcp", "localhost:12345")
	if err != nil {
		log.Fatal("listener error: ", err)
	}
	file, err := ln.(*net.TCPListener).File()
	if err != nil {
		log.Fatal("converting listener to file descr error: ", err)
	}

	cmd := exec.Command("go", "run", "./cmd/worker/main.go")
	var outStd bytes.Buffer
	cmd.Stdout = &outStd
	var outErr bytes.Buffer
	cmd.Stderr = &outErr
	cmd.ExtraFiles = []*os.File{file}
	err = cmd.Start()
	fmt.Println("StartErr:", err)
	fmt.Println("Stdout:", outStd.String())
	fmt.Println("Stderr:", outErr.String())
}
