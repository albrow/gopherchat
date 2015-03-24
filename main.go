package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

const (
	srvAddr         = "224.0.0.1:9999"
	maxDatagramSize = 8192
)

type message struct {
	src  *net.UDPAddr
	body []byte
	n    int
}

func main() {
	// go ping(srvAddr)
	addr, err := net.ResolveUDPAddr("udp", srvAddr)
	if err != nil {
		log.Fatal(err)
	}

	msgs := make(chan message, 100)
	go receiveMessages(addr, msgs)
	go printMessages(msgs)
	go readInput(addr)
	done := make(chan bool)
	<-done
}

func readInput(addr *net.UDPAddr) {
	reader := bufio.NewReader(os.Stdin)
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		msg, err := reader.ReadBytes('\n')
		if err != nil {
			log.Fatal(err)
		}
		if _, err := conn.Write(msg); err != nil {
			log.Fatal(err)
		}
	}
}

func receiveMessages(addr *net.UDPAddr, msgs chan message) {
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	if err := conn.SetReadBuffer(maxDatagramSize); err != nil {
		log.Fatal(err)
	}
	for {
		buf := make([]byte, maxDatagramSize)
		n, src, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}
		if n > 0 {
			msg := message{
				src:  src,
				n:    n,
				body: buf,
			}
			msgs <- msg
		}
	}
}

func printMessages(msgs chan message) {
	for msg := range msgs {
		log.Print(string(msg.body[:msg.n]))
	}
}
