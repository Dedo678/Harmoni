package main

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
	"strings"
)

//Client is a wrapper for connections, giving a name identifier
type client struct {
	name string
	uid  int
	conn net.Conn
	rw   *bufio.ReadWriter
	gob  *gob.Encoder
}

func createClient(c net.Conn) client {
	connectioncount++
	newc := client{}
	newc.conn = c
	newc.rw = bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c))
	newc.gob = gob.NewEncoder(c)
	n, err := newc.rw.ReadString('\n')
	if err != nil {
		log.Println(err)
		newc.name = "User"
	} else {
		n = strings.Trim(n, "\n")
		newc.name = n
	}
	newc.uid = connectioncount
	return newc
}

func (c *client) sendMessage(ch string, message string) error {
	msgwritelock.Lock()
	var err error
	_, err = c.rw.WriteString(message)
	if err != nil {
		return err
	}
	_, err = c.rw.WriteString(ch)
	if err != nil {
		return err
	}
	c.rw.Flush()
	msgwritelock.Unlock()
	return nil
}

func (c *client) receiveMessage() string {
	msg, err := c.rw.ReadString('\n')
	if err != nil {
		log.Println(">>", err)
		c.rw.ReadString('\n')
		return "ERRORINVALID"
	}
	if msg == "NAME\n" {
		return "NAME"
	}
	return msg
}
