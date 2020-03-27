package tcpProxy

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
)

type FooReader struct {}
type FooWriter struct {}

func (f FooReader)Read(b []byte) (int, error){
	fmt.Print("In > ")
	return os.Stdin.Read(b)
}

func (f FooWriter)Write(b []byte) (int, error) {
	return os.Stdout.Write(b)
}
func CustomReadAndWrite() {
	var reader FooReader
	var writer FooWriter
	input := make([]byte, 4096)
	s, err := reader.Read(input)
	if err != nil {
		log.Fatalln("Unable to read data")

	}
	fmt.Printf("Read %d byte from stdin \n", s)
	s, err = writer.Write(input)
	if err != nil {
		log.Fatalln("Unable to write data")
	}
	fmt.Printf("Write %d byte from stdout\n", s)
	if _, err = io.Copy(writer, reader); err != nil {
		log.Fatalln("Unable to write/read data")
	}
}
func echo(conn net.Conn) {
	defer conn.Close()
	b := make([]byte, 512)
	for {
		size, err := conn.Read(b[0:])
		if err == io.EOF {
			log.Println("Client disconnect")
			break
		}
		if err != nil {
			log.Println("Unexpected error")
			break
		}
		log.Printf("Received %d bytes: %s\n",size, string(b))
		log.Println("Writing data")
		if _, err =conn.Write(b[0:size]); err != nil {
			log.Fatalln("Unable to write data")
		}
	}
}

func echoImprove(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	s, err := reader.ReadString('\n')
	if err != nil{
		log.Fatalln("Unable to read data")
	}
	log.Printf("Read %d bytes: %s", len(s), s)
	log.Println("Writing Data")
	writer := bufio.NewWriter(conn)
	_, err = writer.WriteString(s);
	if err != nil {
		log.Fatalln("Unable to write Data")
	}
	_ = writer.Flush()
}
func echoFinal(conn net.Conn) {
	defer conn.Close()
	if _ ,err := io.Copy(conn, conn); err != nil {
		log.Fatalln("Unable to write/read data")
	}
}
func EchoServer() {
	listener, err := net.Listen("tcp", ":2080")
	if err != nil {
		log.Fatalln("Unable to bind port")
	}
	log.Println("Listening on 0.0.0.0:2080")
	for {
		conn , err := listener.Accept()
		log.Println("Received connection")
		if err != nil {
			log.Fatalln("Unable tp accept connection")
		}
		go echoImprove(conn)
	}
}

func ProxyServer() {
	listener, err := net.Listen("tcp",":80"); if err != nil {
		log.Fatalln("Unable to bind port")
	}
	log.Println("Hello")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("Unable to accept connection")
		}
		go handleProxy (conn)
	}

}
func handleProxy(conn net.Conn) {
	dst, err := net.Dial("tcp", "facebook.com:80")
	if err != nil {
		log.Fatalln("Unable to connect to our unreachable host")
	}
	defer dst.Close()
	go func() {
		if _, err := io.Copy(dst, conn); err != nil {
			log.Fatalln(err, "w")
		}
	}()
	if _, err := io.Copy(conn, dst); err != nil {
		log.Fatalln(err, "r")
	}
}

func useCmd(conn net.Conn) {
	cmd := exec.Command("/bin/sh", "-i")
	if err := cmd.Run(); err != nil {
		log.Fatalln("Run command failure")
	}
	rp, wp := io.Pipe()
	cmd.Stdin = conn
	cmd.Stdout = wp
	go io.Copy(conn, rp)
	cmd.Run()
	conn.Close()
}