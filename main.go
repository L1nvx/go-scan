package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

const (
	ipAddr      = "10.10.11.252"
	maxWorkers  = 1000
	timeoutSecs = 2
)

var wg sync.WaitGroup

func checkPortWorker(jobs <-chan int) {
	defer wg.Done()
	for port := range jobs {
		conn, err := net.DialTimeout("tcp", ipAddr+":"+strconv.Itoa(port), timeoutSecs*time.Second)
		if err != nil {
			continue
		}
		defer conn.Close()
		fmt.Println("[+] Port", port, "is open")
	}
}

func main() {
	start := time.Now()
	jobs := make(chan int, maxWorkers)
	wg.Add(maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		go checkPortWorker(jobs)
	}
	for port := 1; port <= 65535; port++ {
		jobs <- port
	}
	close(jobs)
	wg.Wait()
	fmt.Println("[*] Tiempo total:", time.Since(start))
}
