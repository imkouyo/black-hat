package portScanner

import (
	"fmt"
	"net"
	"sort"
)

func PortScanner() {
	ports := make(chan int, 10)
	results := make(chan int)
	var openPort []int
	for i := 0; i < cap(ports); i++ {
		go scanner(ports, results)
	}

	go func() {
		for i := 80; i < 83; i++ {
			ports <- i
		}
	}()
	for i := 0; i < 3; i++ {
		port := <-results
		if port != 0 {
			openPort = append(openPort, port)
		}
	}
	close(ports)
	close(results)
	sort.Ints(openPort)
	for _, port := range openPort {
		fmt.Printf("%d open\n", port)
	}

}

func scanner(ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("tw.yahoo.com:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		_ = conn.Close()
		results <- p
	}
}