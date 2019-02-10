package main

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
	"strings"
	"time"
)

var name string

type serverDetails struct {
	Name     string
	Channels []string
	Clients  []string
}

type connection struct {
	name    string
	conn    net.Conn
	rw      *bufio.ReadWriter
	details serverDetails
}

func (c *connection) listen() {
	cwr := bufio.NewReadWriter(bufio.NewReader(c.conn), bufio.NewWriter(c.conn))
	cgob := gob.NewDecoder(c.conn)
	for {
		log.Println(">> Checking for incoming messages")
		msg, _ := cwr.ReadString('\n')
		msg = strings.Trim(msg, "\n")
		switch msg {
		case "SERVERDETAILS":
			log.Println(">> Incoming GOB")
			err := cgob.Decode(&c.details)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Println(">> Received server details")
			log.Println(">> ", c.details)
		case "STRINGMESSAGE":
			log.Println(">> Incoming string")
			strmsg, err := cwr.ReadString('\n')
			if err != nil {
				log.Println(">> " + err.Error())
				continue
			}
			strmsg = strings.Trim(strmsg, "\n")
			log.Println(">> Receive string message")
			log.Println(">> " + strmsg)
		}
	}
}

func (c *connection) sendStringMsg(msg string) {
	c.rw.WriteString("STRINGMESSAGE\n")
	c.rw.WriteString(msg + "\n")
	c.rw.Flush()
}

func createConnection(name string, address string) connection {
	conn, _ := net.Dial("tcp", "localhost:8888")
	rwconn := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	myclient := connection{
		name:    name,
		conn:    conn,
		rw:      rwconn,
		details: serverDetails{},
	}
	go myclient.listen()
	return myclient
}

func main() {
	myclient := createConnection("Dedo", "localhost:8888")
	time.Sleep(time.Second)
	myclient.rw.WriteString("Dedo\n")
	myclient.rw.Flush()
	for {
	}
}
