package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
)

var name string

func open(address string) (net.Conn, error) {
	var client net.Conn
	var err error
	var attempts int
	for {
		client, err = net.Dial("tcp", address)
		if err != nil {
			if attempts != 5 {
				attempts++
				continue
			} else {
				return nil, err
			}
		}
		return client, nil
	}
}

func client(address string) {
	var client net.Conn
	var err error

	client, err = open(address)
	defer client.Close()
	if err != nil {
		log.Println(">> Failed to connect --- " + err.Error())
		return
	}
	log.Println(">> Connected to " + client.RemoteAddr().String())
	go func() {
		for {
			msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			chw := bufio.NewWriter(client)
			chw.WriteString(name + msg)
			chw.Flush()
		}
	}()
	func() {
		for {
			msg, err := bufio.NewReader(client).ReadString('\n')
			if err != nil {
				log.Println(">> Could not read from buffer --- " + err.Error())
				break
			}
			msg = strings.Trim(msg, "\n")
			log.Println(">> " + msg)
		}
	}()
}

func main() {
	name = "dedo678"
	client("localhost:8888")

}
