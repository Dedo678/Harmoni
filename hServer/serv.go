package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

//Client is a wrapper for connections, giving a name identifier
type Client struct {
	name string
	conn net.Conn
}

//Server is a wrapper for underlying listener and Client list
type Server struct {
	t  string
	ls net.Listener
	cl []Client

	stopconn chan bool
}

//Loops forever watching for connections and adds it to Server Client list if one is found
func (s *Server) listenNewConn() {
	select {
	case <-s.stopconn: //Disables listener
		return
	default:
		for {
			c, err := s.ls.Accept()
			if err != nil {
				log.Println(err.Error())
				continue
			}
			log.Println(">> Client connected")
			cl := Client{}
			cl.conn = c
			clrw := bufio.NewReadWriter(bufio.NewReader(cl.conn), bufio.NewWriter(cl.conn))
			n, err := clrw.ReadString('\n') //Receive name data
			n = strings.Trim(n, "\n")
			cl.name = n
			clrw.WriteString("SERVERDETAILS\n") //TODO:Send actual Server details
			clrw.Flush()
			s.cl = append(s.cl, cl)
		}
	}
}

func (s *Server) listenMsgs() {

}

//NewServer does initial Server startup
func NewServer(title string) Server {
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		panic(err)
	}
	serv := Server{
		t:        title,
		ls:       listener,
		cl:       make([]Client, 0),
		stopconn: make(chan bool),
	}
	return serv
}

//StartConn starts up listeners
func (s *Server) StartConn() {
	go s.listenNewConn()
}
