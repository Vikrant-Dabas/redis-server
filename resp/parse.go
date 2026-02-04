package resp

import (
	"bufio"
	"fmt"
)

func Parse(r *bufio.Reader)([]byte,error){
	format,err := ReadCommand(r)
	if err != nil{
		return nil,err
	}

	return format.Marshal()
}

func ReadCommand(r *bufio.Reader)(*Format,error){
	payload,err := readValidateInput(r)
	if err != nil{
		return nil,err
	}
	format := Format{}
	msgType := payload[0]
	switch msgType{
	case typeSimple, typeError, typeInt:
		format.Type = msgType
		if err := format.SimpleUnmarshal(payload);err != nil{
			return nil,err
		}
	case typeBulk:
		format.Type = msgType
		if err := format.BulkUnmarshal(r,payload);err != nil{
			return nil,err
		}
	case typeArray:
		format.Type = msgType
		if err := format.ArrayUnmarshal(r,payload);err != nil{
			return nil,err
		}
	default:
		return nil,fmt.Errorf("invalid syntax: message type %c",msgType)
	}
	return &format,nil
}
