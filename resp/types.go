package resp

type Input struct {
	Type    byte
	Payload []byte
}

type Format struct {
	Type    byte
	Size    int
	Payload [][]byte
}

const (
	typeSimple = '+'
	typeBulk   = '$'
	typeArray  = '*'
	typeError  = '-'
	typeInt    = ':'
)
