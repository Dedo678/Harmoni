package main

func main() {
	myServer := NewServer("My Server")
	myServer.listenNewConn()
}
