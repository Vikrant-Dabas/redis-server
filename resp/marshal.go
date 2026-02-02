package resp

import (
	"fmt"
	"strconv"
)

func (f *Format)Marshal() ([]byte,error){
	switch f.Type{
	case typeSimple,typeError,typeInt:
		return f.SimpleMarshal()
	case typeBulk:
		return f.BulkMarshal()
	case typeArray:
		return f.ArrayMarshal()
	default:
		return nil,fmt.Errorf("marshal error: type not supported\n")
	}
} 

func (f *Format)SimpleMarshal()([]byte,error){
	if !(f.Type == typeSimple || f.Type == typeError || f.Type == typeInt){
		return nil,fmt.Errorf("marhal error: wrong msg type%c\n",f.Type)
	}
	if f.Type == typeInt{
		if _,err := strconv.ParseInt(string(f.Payload),10,64);err != nil{
			return nil,err
		}
	}
	if terminatorInBetween(f.Payload){
		return nil,fmt.Errorf("marshal error: unexpexted terminator\n")
	}
	output := fmt.Sprintf("%c%s%s",f.Type,f.Payload,terminator)
	return []byte(output),nil
}

func (f *Format)BulkMarshal()([]byte,error){
	if f.Type != typeBulk{
		return nil,fmt.Errorf("marshal error: wrong msg type %c\n",f.Type)
	}
	if terminatorInBetween(f.Payload){
		return nil,fmt.Errorf("marshal error: unexpexted terminator\n")
	}
	size := len(f.Payload)
	sizeStr := strconv.Itoa(size)
	output := fmt.Sprintf("$%s%s%s%s",sizeStr,terminator,f.Payload,terminator)
	return []byte(output),nil
}

func (f *Format)ArrayMarshal()([]byte,error){
	size := len(f.ArrayPayload)
	if f.Type != typeArray{
		return nil,fmt.Errorf("marshal error: wrong error type %c\n",f.Type)
	}
	output := fmt.Sprintf("*%d%s",size,terminator)
	for _,format := range f.ArrayPayload{
		str,err := format.Marshal()
		if err != nil{
			return nil,err
		}
		output += string(str)
	}
	return []byte(output),nil
}