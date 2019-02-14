package main

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
	"strings"
)

//ServerDetails holds information about the server, which is passed from the server when a change in serverdetails occurs
type ServerDetails struct {
	Name     string
	Channels []string
	Clients  []string
}

type connection struct {
	name     string
	conn     net.Conn
	gob      *gob.Decoder
	rw       *bufio.ReadWriter
	sdetails *ServerDetails
}

func (c connection) receiveMessages() {
	defer c.conn.Close()
	for {
		chmsg, _ := c.rw.ReadString('\n')
		switch chmsg {
		case "DETAILS\n":
			c.receiveInfo()
		default:
			msg, _ := c.rw.ReadString('\n')
			chmsg = strings.Trim(chmsg, "\n")
			msg = strings.Trim(msg, "\n")
			log.Println(">> " + chmsg + ": " + msg)
		}
	}
}

func (c connection) sendMessage(ch string, msg string) {

}

func (c connection) receiveInfo() {
	c.sdetails = nil
	c.gob.Decode(&c.sdetails)
	log.Println(">> Received server info update", c.sdetails)
}

func (c connection) changeName(n string) {
	c.name = n
	c.rw.WriteString("NAME\n")
	c.rw.WriteString(c.name + "\n")
}

func createConnection(name string, address string) connection {
	var conn net.Conn
	var err error
	for {
		if conn, err = net.Dial("tcp", address); err != nil {
			log.Println(err)
			continue
		}
		break
	}
	log.Println(">> Connected!")
	newclient := connection{
		name:     name,
		conn:     conn,
		gob:      gob.NewDecoder(conn),
		rw:       bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)),
		sdetails: &ServerDetails{},
	}
	newclient.rw.WriteString(newclient.name + "\n")
	newclient.rw.Flush()
	return newclient
}

func main() {
	myclient := createConnection("Dedo", "127.0.0.1:8888")
	myclient.receiveMessages()

	for {
	}
}
