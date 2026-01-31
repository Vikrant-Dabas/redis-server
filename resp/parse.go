package resp

import "fmt"

func Parse(input *Input) ([]byte, error) {
	switch input.Type {
	case typeSimple:
		input.Payload = append(input.Payload, []byte(" Simple String")...)
	case typeBulk:
		input.Payload = append(input.Payload, []byte(" Bulk String")...)
	case typeInt:
		input.Payload = append(input.Payload, []byte(" Integer")...)
	case typeError:
		input.Payload = append(input.Payload, []byte(" Simple Error")...)
	case typeArray:
		input.Payload = append(input.Payload, []byte(" Array")...)
	}
	return makeByteSlice(input),nil
}

func makeByteSlice(input *Input) []byte {
	return []byte(fmt.Sprintf("type : %s :: %s", string(input.Type), input.Payload))
}
