package main

import (
	"bufio"
	"fmt"
	"net"
)

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan Message
}

type Message struct {
	from    string
	payload []byte
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan Message, 10),
	}
}

func (s *Server) Start() error {
	fmt.Printf("\033[32mserver started on %s\033[0m\n", s.listenAddr)
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.AcceptLoop()
	<-s.quitch
	close(s.msgch)
	return nil
}

func (s *Server) AcceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Printf("error in accepting connection: %s\n", err.Error())
			continue
		}
		fmt.Printf("connection established by %s\n",conn.RemoteAddr().String())
		conn.Write([]byte("Welcome\n"))
		go s.ReadLoop(conn)
	}
}

func (s *Server) ReadLoop(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan(){
		payload := make([]byte,len(scanner.Bytes()))
		copy(payload,scanner.Bytes())

		s.msgch <- Message{
			from: conn.RemoteAddr().Network(),
			payload: payload,
		}
	}

	if err := scanner.Err();err != nil{
		fmt.Printf("connection error from %s : %s",conn.RemoteAddr().String(),err)
	}
	fmt.Printf("connection closed by %s",conn.RemoteAddr().String())
}
