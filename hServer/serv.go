package main

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
	"strings"
)

type serverDetails struct {
	Name     string
	Channels []string
	Clients  []string
}

//Client is a wrapper for connections, giving a name identifier
type Client struct {
	name string
	conn net.Conn
	rw   *bufio.ReadWriter
}

//Server is a wrapper for underlying listener and Client list
type Server struct {
	t       string
	ls      net.Listener
	cl      []Client
	details serverDetails

	stopconn chan bool
}

func (s *Server) updateDetails() {
	clist := make([]string, len(s.cl))
	for i, v := range s.cl {
		clist[i] = v.name
	}
	s.details.Clients = clist
}

func (s *Server) sendDetails(c Client) {
	s.updateDetails()
	c.rw.WriteString("SERVERDETAILS\n")
	c.rw.Flush()
	g := gob.NewEncoder(c.conn)
	err := g.Encode(s.details)
	if err != nil {
		log.Println(">> Failed to send Server details over GOB - ", err.Error())
	}
}

//Loops forever watching for connections and adds it to Server Client list. Receives client name and sends out server details
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
			log.Println(">> Client connecting..")
			newclient := Client{
				conn: c,
				rw:   bufio.NewReadWriter(bufio.NewReader(c), bufio.NewWriter(c)),
			}
			newclient.name, _ = newclient.rw.ReadString('\n')
			newclient.name = strings.Trim(newclient.name, "\n")
			s.cl = append(s.cl, newclient)
			s.sendDetails(newclient)
			log.Println(">> Client " + newclient.name + " has fully connected")
		}
	}
}

func (s *Server) listenClients() { //TODO:Main message relay
	for {

	}
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
		details:  serverDetails{Name: title},
		stopconn: make(chan bool),
	}
	return serv
}

//StartConn starts up listeners
func (s *Server) StartConn() {
	go s.listenNewConn()
	go s.listenClients()
}
