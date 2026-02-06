package resp

type Input struct {
	Type    byte
	Payload []byte
}

type Format struct {
	Type         byte
	Payload      []byte
	ArrayPayload []Format
}

const (
	TypeSimple byte = '+'
	TypeBulk   byte = '$'
	TypeArray  byte = '*'
	TypeError  byte = '-'
	TypeInt    byte = ':'
	TypeNil    byte = 0
)

const terminator string = "\r\n"

var AllTypes = []byte{TypeArray, TypeSimple, TypeBulk, TypeError, TypeInt}

var NilResp = []byte("$-1\r\n")
