package commands

import (
	"fmt"

	"github.com/Vikrant-Dabas/redis/resp"
)

func ExecuteUniversal(cmd string, input [][]byte) (*resp.Format, error) {
	switch cmd {
	case "PING":
		return ping(input)
	}
	return nil, nil
}

func ping(input [][]byte) (*resp.Format, error) {
	size := len(input)
	if size == 0 {
		return &resp.Format{
			Type:    resp.TypeSimple,
			Payload: []byte("PONG"),
		}, nil
	} else if size == 1 {
		return &resp.Format{
			Type:    resp.TypeBulk,
			Payload: input[0],
		}, nil
	}
	return nil, fmt.Errorf("invalid PING command: %d", size)
}
