package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"syscall"

	"golang.org/x/sys/unix"
)

type myHandler struct {
	srvID int
}

func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("service %d got / request", h.srvID)
	io.WriteString(w, fmt.Sprintf("Hello from service %d!\n", h.srvID))
}

func main() {
	n := flag.Int("n", 1, "number of service")
	flag.Parse()
	fmt.Println("hello! This is service", *n)
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			var err error
			c.Control(func(fd uintptr) {
				err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
				if err != nil {
					return
				}
				err = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
				if err != nil {
					return
				}
			})
			return err
		},
	}
	l, _ := lc.Listen(context.Background(), "tcp", "localhost:5001")

	wg := sync.WaitGroup{}
	wg.Add(1)
	// start server1
	go func(n int) {
		defer wg.Done()
		serverID := n
		mux := http.NewServeMux()
		handler := &myHandler{srvID: n}
		mux.Handle("/", handler)
		srv := http.Server{
			Handler: mux,
		}
		if err := srv.Serve(l); err != nil {
			fmt.Printf("Service %d got error: %s\n", serverID, err)
		}
	}(*n)
	wg.Wait()
	fmt.Println("We are done")
}
