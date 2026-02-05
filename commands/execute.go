package commands

import (
	"fmt"

	"github.com/Vikrant-Dabas/redis/resp"
)

func Execute(input [][]byte) (*resp.Format, error) {
	noOfCommands := len(input)
	switch noOfCommands{
	case 1:
		cmd := input[0]
		return ExecuteSingle(cmd)
	default:
		return nil,fmt.Errorf("invalid command")
	}
}

func PingPong() *resp.Format {
	return &resp.Format{
		Type: resp.TypeSimple,
		Payload: []byte("PONG"),
	}
}

func ExecuteSingle(input []byte)(*resp.Format,error){
	switch string(input){
	case "PING":
		return PingPong(),nil
	default:
		return nil,fmt.Errorf("invalid command: %s",input)
	}
}