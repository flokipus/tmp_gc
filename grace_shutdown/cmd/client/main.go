package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	for i := 0; i < 100; i++ {
		resp, err := http.Get("http://localhost:5001/")
		if err != nil {
			fmt.Println("Got error:", err)
		} else {
			s, _ := io.ReadAll(resp.Body)
			fmt.Println("Got response:", string(s))
		}
	}
}
