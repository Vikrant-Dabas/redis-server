package main

import (
	"bufio"
	"fmt"
	"net"

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

type Input struct{
	Type byte
	Payload []byte
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
		reader := bufio.NewReader(conn)
		go s.ReadLoop(conn,reader)
	}
}

func (s *Server) ReadLoop(conn net.Conn,r *bufio.Reader) {
	defer conn.Close()

	for{
		payload,err := ReadInput(r)
		if err != nil{
			fmt.Println(err)
			break
		}
		s.msgch<-Message{
			from: conn.RemoteAddr().String(),
			payload: payload,
		}
	}
	
	fmt.Printf("connection closed by %s\n",conn.RemoteAddr().String())
}

func ReadInput(r *bufio.Reader)([]byte,error){
	input := Input{}
	for{
		b,err := r.ReadByte()
		if err != nil{
			return nil,err
		}

		if b == '\n'{
			if input.Type == 0{
				return []byte("empty message"),nil
			} else {
				return resp.Parse((*resp.Input)(&input))
			}	
		}

		if input.Type == 0{
			input.Type = b
		} else {
			input.Payload = append(input.Payload, b)
		}
	}
}
