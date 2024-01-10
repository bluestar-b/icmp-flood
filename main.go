package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func ping(host string, packetSize int, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		start := time.Now()

		conn, err := net.Dial("ip4:icmp", host)
		if err != nil {
			fmt.Println("Error connecting to the host:", err)
			continue
		}

		header := make([]byte, 8+packetSize)

		header[0] = 8
		header[1] = 0
		header[2] = 0
		header[3] = 0

		header[4] = 0
		header[5] = 0
		header[6] = 0
		header[7] = byte(time.Now().Nanosecond())
		for j := 8; j < len(header); j++ {
			header[j] = byte(j % 256)
		}

		conn.Write(header)

		reply := make([]byte, 1024)
		conn.SetReadDeadline(time.Now().Add(time.Second * 2))
		_, err = conn.Read(reply)
		if err != nil {
			fmt.Println("No response from the host:", err)
		} else {
			elapsed := time.Since(start)
			fmt.Printf("Ping to %s, RTT: %v\n", host, elapsed)
		}

	}
}

func main() {
	host := "nigga56.net"
	packetSize := 64
	numThreads := 4

	var wg sync.WaitGroup

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go ping(host, packetSize, &wg)
	}

	wg.Wait()
}
