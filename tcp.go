package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/Vikrant-Dabas/redis/commands"
	"github.com/Vikrant-Dabas/redis/resp"
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
		reader := bufio.NewReader(conn)
		go s.ReadLoop(conn,reader)
	}
}

func (s *Server) ReadLoop(conn net.Conn,r *bufio.Reader) {
	defer conn.Close()
	w := bufio.NewWriter(conn)
	for{
		format,err := resp.ReadCommand(r)
		if err != nil{
			w.WriteString("-ERR " + err.Error() + "\r\n")
			w.Flush()
			continue
		}
		cmd := format.ToByteMatrix()
		format,err = commands.Execute(cmd)
		if err != nil{
			w.WriteString("-ERR " + err.Error() + "\r\n")
			w.Flush()
			continue
		}
		output,err := format.Marshal()
		if err != nil{
			w.WriteString("-ERR " + err.Error() + "\r\n")
			w.Flush()
			continue
		}
		_, err = w.Write(output)
		if err != nil {
    		return
		}
		w.Flush()

		s.msgch<-Message{
			from: conn.RemoteAddr().String(),
			payload: output, 
		}
	}
	
}

