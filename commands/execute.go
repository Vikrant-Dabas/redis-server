package commands

import (
	"fmt"
	"strings"

	"github.com/Vikrant-Dabas/redis/db"
	"github.com/Vikrant-Dabas/redis/resp"
)

func Execute(db db.DB, input [][]byte) (*resp.Format, error) {
	cmd := strings.ToUpper(string(input[0]))
	cmdType, ok := CmdTypes[cmd]
	if !ok {
		return nil, fmt.Errorf("invalid command %s", cmd)
	}
	switch cmdType {
	case CmdString:
		return ExecuteString(db, cmd, input[1:])
	case CmdUniversal:
		return ExecuteUniversal(cmd, input[1:])
	case CmdList:
		return ExecuteList(db, cmd, input[1:])
	default:
		return nil, fmt.Errorf("command not supported: %s", cmd)
	}
}
