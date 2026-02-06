package main

import (
	"fmt"
	"log"
)

func main() {
	port := ":6380"
	s := NewServer(port)
	go func() {
		for msg := range s.msgch {
			fmt.Printf("message from %s : %s\n", msg.from, msg.payload)
		}
	}()
	err := s.Start()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("thanks for the visit")
	}
}
