package commands

type CmdType uint8

const (
	CmdUniversal CmdType = iota
	CmdString
	CmdList
)

var CmdTypes = map[string]CmdType{
	// String
	"GET":    CmdString,
	"SET":    CmdString,
	"INCR":   CmdString,
	"INCRBY": CmdString,
	"DECR":   CmdString,
	"DECRBY": CmdString,
	"MSET":   CmdString,
	"MGET":   CmdString,

	// Universal
	"PING": CmdUniversal,

	// List
	"LPUSH":  CmdList,
	"RPUSH":  CmdList,
	"LPOP":   CmdList,
	"RPOP":   CmdList,
	"LLEN":   CmdList,
	"LRANGE": CmdList,
	"LTRIM":  CmdList,
}
