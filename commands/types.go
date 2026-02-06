package commands

type CmdType uint8

const(
	CmdUversal CmdType = iota
	CmdString
	CmdList
)

var CmdTypes = map[string]CmdType{
	// String
	"GET":CmdString,
	"SET":CmdString,
}