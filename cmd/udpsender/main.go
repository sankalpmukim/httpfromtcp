package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	port := "42069"
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:"+port)
	if err != nil {
		log.Fatalf("cannot resolve udp address %v", err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatalf("cannot dial udp address %v", err)
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		inp, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("can't read string")
		}
		_, err = conn.Write([]byte(inp))
		if err != nil {
			log.Fatal("can't write to the connection")
		}
	}
}
