package resp

import (
	"bufio"
	"fmt"
	"strings"
)

func ReadCommand(r *bufio.Reader) (*Format, error) {
	payload, err := readValidateInput(r)
	if err != nil {
		return nil, err
	}
	format := Format{}
	msgType := payload[0]
	switch msgType {
	case TypeSimple, TypeError, TypeInt:
		format.Type = msgType
		if err := format.SimpleUnmarshal(payload); err != nil {
			return nil, err
		}
	case TypeBulk:
		format.Type = msgType
		if err := format.BulkUnmarshal(r, payload); err != nil {
			return nil, err
		}
	case TypeArray:
		format.Type = msgType
		if err := format.ArrayUnmarshal(r, payload); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid syntax: message type %c", msgType)
	}
	return &format, nil
}

func (f *Format) ToByteMatrix() [][]byte {
	var output [][]byte
	switch f.Type {
	case TypeArray:
		for _, fmt := range f.ArrayPayload {
			temp := fmt.ToByteMatrix()
			output = append(output, temp...)
		}
	default:
		upper := strings.ToUpper(string(f.Payload))
		output = append(output, []byte(upper))
	}
	return output
}
