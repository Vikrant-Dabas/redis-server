package commands

type CmdType uint8

const (
	CmdUniversal CmdType = iota
	CmdString
	CmdList
	CmdSet
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
	"PING":   CmdUniversal,
	"EXPIRE": CmdUniversal,
	"TTL":CmdUniversal,
	"PERSIST":CmdUniversal,

	// List
	"LPUSH":  CmdList,
	"RPUSH":  CmdList,
	"LPOP":   CmdList,
	"RPOP":   CmdList,
	"LLEN":   CmdList,
	"LRANGE": CmdList,
	"LTRIM":  CmdList,

	// Set
	"SADD":        CmdSet,
	"SREM":        CmdSet,
	"SPOP":        CmdSet,
	"SRANDMEMBER": CmdSet,
	"SISMEMBER":   CmdSet,
	"SMISMEMBER":  CmdSet,
	"SINTER":      CmdSet,
	"SUNION":      CmdSet,
	"SDIFF":       CmdSet,
	"SCARD":       CmdSet,
	"SMEMBERS":    CmdSet,
}
