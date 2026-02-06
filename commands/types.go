package commands

type CmdType uint8

const (
	CmdUniversal CmdType = iota
	CmdString
	CmdList
)

var CmdTypes = map[string]CmdType{
	// String
	"GET": CmdString,
	"SET": CmdString,

	// Universal
	"PING": CmdUniversal,
}
