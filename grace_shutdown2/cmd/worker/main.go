package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
)

type myHandler struct {
	srvID int
}

func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("service %d got / request", h.srvID)
	io.WriteString(w, fmt.Sprintf("Hello from service %d!\n", h.srvID))
}

func main() {
	// получаем сокет
	file := os.NewFile(3, "listener")
	ln, err := net.FileListener(file)
	if err != nil {
		log.Fatal("error to get listener from file")
	}
	// стартуем сервис, как и раньше
	n := flag.Int("n", 1, "number of service")
	flag.Parse()
	fmt.Println("hello! This is service", *n)
	wg := sync.WaitGroup{}
	wg.Add(1)
	// start server
	go func(n int) {
		defer wg.Done()
		serverID := n
		mux := http.NewServeMux()
		handler := &myHandler{srvID: n}
		mux.Handle("/", handler)
		srv := http.Server{
			Handler: mux,
		}
		if err := srv.Serve(ln); err != nil {
			fmt.Printf("Service %d got error: %s\n", serverID, err)
		}
	}(*n)
	wg.Wait()
	fmt.Println("We are done")
}
