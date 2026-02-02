package resp

type Input struct {
	Type    byte
	Payload []byte
}

type Format struct {
	Type         byte
	Size         int
	Payload      [][]byte
	ArrayPayload []Format
}

const (
	typeSimple byte = '+'
	typeBulk   byte = '$'
	typeArray  byte = '*'
	typeError  byte = '-'
	typeInt    byte = ':'
)

const terminator string = "\r\n"

var AllTypes = []byte{typeArray, typeSimple, typeBulk, typeError, typeInt}
