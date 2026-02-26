package commands

import (
	"fmt"
	"strings"

	"github.com/Vikrant-Dabas/redis/db"
	"github.com/Vikrant-Dabas/redis/resp"
)

func Execute(store *db.Store,input [][]byte) (*resp.Format, error) {
	cmd := strings.ToUpper(string(input[0]))
	cmdType, ok := CmdTypes[cmd]
	if !ok {
		return nil, fmt.Errorf("invalid command %s", cmd)
	}
	switch cmdType {
	case CmdString:
		return ExecuteString(store, cmd, input[1:])
	case CmdUniversal:
		return ExecuteUniversal(store,cmd, input[1:])
	case CmdList:
		return ExecuteList(store, cmd, input[1:])
	case CmdSet:
		return ExecuteSet(store, cmd, input[1:])
	default:
		return nil, fmt.Errorf("command not supported: %s", cmd)
	}
}
