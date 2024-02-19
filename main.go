package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	ipAddr      string
	maxWorkers  int
	timeoutSecs int
)

var wg sync.WaitGroup

func checkPortWorker(jobs <-chan int) {
	defer wg.Done()
	for port := range jobs {
		conn, err := net.DialTimeout("tcp", ipAddr+":"+strconv.Itoa(port), time.Duration(timeoutSecs)*time.Second)
		if err != nil {
			continue
		}
		defer conn.Close()
		fmt.Println("[+] Port", port, "is open")
	}
}

func main() {
	flag.StringVar(&ipAddr, "target", "", "ip address to scan ports.")
	flag.IntVar(&maxWorkers, "workers", 1000, "num of workers.")
	flag.IntVar(&timeoutSecs, "timeout", 2, "seconds for port connection.")
	flag.Parse()
	if ipAddr == "" {
		fmt.Println("[!] usage", os.Args[0], "-target <ip>")
		flag.PrintDefaults()
		return
	}
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
