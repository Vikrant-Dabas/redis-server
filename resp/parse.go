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
	var ans []byte
	if format.Type == typeArray{
		ans = append(ans, format.Type,'-','\n')
		for _,f := range format.ArrayPayload{
			ansString := fmt.Sprintf("%c - %s\n",f.Type,f.Payload[0])
			ans = append(ans, []byte(ansString)...)
		}
	} else {
		ans = append(ans, format.Type,'-')
		ans = append(ans, format.Payload[0]...)
	}
	return ans,nil
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
		if err := SimpleUnmarshal(&format,payload);err != nil{
			return nil,err
		}
	case typeBulk:
		format.Type = msgType
		if err := BulkUnmarshal(r,&format,payload);err != nil{
			return nil,err
		}
	case typeArray:
		format.Type = msgType
		if err := ArrayUnmarshal(r,&format,payload);err != nil{
			return nil,err
		}
	default:
		return nil,fmt.Errorf("invalid syntax: message type %c",msgType)
	}
	return &format,nil
}

func contains(b []byte) (byte,bool){
	for _,i := range b{
		for _,ch := range AllTypes{
			if i == ch{
				return ch,true
			}
		}
	}
	return 0,false;
}

func validTerminator(b []byte)bool{
	n := len(b)
	if n < 2{
		return false
	}
	if b[n-2] != '\r' || b[n-1] != '\n'{
		return false
	}

	for _,ch := range b[:n-2]{
		if ch == '\r' || ch == '\n'{
			return false
		}
	}
	return true
}

func readValidateInput(r *bufio.Reader)([]byte,error){
	payload,err := r.ReadBytes('\n')	
	fmt.Printf("\033[94m%s-%d\033[0m\n",payload,len(payload))
	if err != nil{
		return nil,err
	}
	if(!validTerminator(payload)){
		return nil,fmt.Errorf("syntax error: terminator: %d-%s\n",len(payload),payload)
	}
	if _,ok := contains([]byte{payload[0]});ok{
		if ch,ok := contains(payload[1:]);ok{
			return nil,fmt.Errorf("syntax error: unexpted message type - %c",ch)
		}
	} else if ch,ok := contains(payload);ok{
		return nil,fmt.Errorf("syntax error: unexpted message type - %c",ch)
	}
	return payload,nil
}