package main

import (
	"log"
	"net"
	"strings"
	"sync"
)

var connectioncount int
var infolock sync.Mutex
var msgwritelock sync.Mutex

//ServerDetails holds information about the server, which it will pass to clients when it updates
type ServerDetails struct {
	Name     string
	Channels []string
	Clients  []string
}

type channel struct {
	name string
}

//Server is a wrapper for underlying listener and Client list
type server struct {
	name     string
	ls       net.Listener
	clients  map[int]client
	channels []channel
	details  ServerDetails
}

//Sets up the server data type with the basics, gives it a default channel
func newServer(title string) server {
	log.Println(">> Creating server..")
	var err error
	s := server{}
	s.name = title
	s.ls, err = net.Listen("tcp", "localhost:8888")
	if err != nil {
		panic(err)
	}
	s.clients = make(map[int]client)
	s.channels = make([]channel, 0)
	s.details = ServerDetails{
		Name:     s.name,
		Channels: make([]string, 0),
		Clients:  make([]string, 0),
	}
	s.addChannel("Default Channel")
	log.Println(">> Server finalized")
	return s
}

func (s *server) listen() {
	defer s.ls.Close()
	var conn net.Conn
	var err error
	for {
		log.Println(">> Waiting for client...")
		conn, err = s.ls.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		newClient := createClient(conn)
		log.Println(">> Client connect! : " + newClient.name)
		s.addClient(newClient)
		go s.runClient(newClient)

	}
}

//Client listens for message and bounces it back to user. Takes a channel first then the message
func (s *server) runClient(c client) {
	defer c.conn.Close()
	for {
		ch := c.receiveMessage()
		switch ch {
		case "NAME":
			newname := c.receiveMessage()
			c.name = newname
			s.updateinfo()
		case "ERRORINVALID":
			log.Println(">> Invalid command from " + c.name)
			log.Println(">> Disconnecting client " + c.name)
			s.removeClient(c)
			return
		default:
			switch s.checkChannel(ch) {
			case "NULLCHANNEL":
				log.Println(">> Invalid channel")
				msg := c.receiveMessage()
				log.Println(">> Unsent msg: " + strings.Trim(msg, "\n"))
			default: //If not errors have occured, 'ch' will contain a '\n' at the end
				msg := c.receiveMessage()
				err := c.sendMessage(ch, msg)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func (s *server) sendInfo() {
	detailscopy := s.details
	msgwritelock.Lock()
	for _, v := range s.clients {
		v.rw.WriteString("DETAILS\n")
		v.rw.Flush()
		v.gob.Encode(detailscopy)
	}
	msgwritelock.Unlock()
}

func (s *server) updateinfo() {
	infolock.Lock()
	s.details.Name = s.name
	s.details.Clients = nil
	s.details.Channels = nil
	for _, v := range s.clients {
		s.details.Clients = append(s.details.Clients, v.name)
	}
	for _, v := range s.channels {
		s.details.Channels = append(s.details.Channels, v.name)
	}
	s.sendInfo()
	log.Println(">> Updated and sent out server info")
	infolock.Unlock()
}

func (s *server) addChannel(n string) {
	log.Println(">> Added channel: " + n)
	s.channels = append(s.channels, channel{name: n})
	s.updateinfo()
}

func (s *server) removeChannel(c channel) {
	newchanlist := make([]channel, len(s.channels)-1, len(s.channels)-1)
	for _, v := range s.channels {
		if v != c {
			newchanlist = append(newchanlist, v)
		}
	}
	s.channels = nil
	s.channels = newchanlist
	s.updateinfo()
}

func (s *server) checkChannel(chstr string) string {
	for _, v := range s.channels {
		if strings.Trim(chstr, "\n") == v.name {
			return chstr
		}
	}
	return "NULLCHANNEL"
}

func (s *server) addClient(c client) {
	s.clients[c.uid] = c
	s.updateinfo()
}

func (s *server) removeClient(c client) {
	delete(s.clients, c.uid)
	s.updateinfo()
}
