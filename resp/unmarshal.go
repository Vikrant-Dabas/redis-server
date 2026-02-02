package resp

import (
	"bufio"
	"fmt"
	"strconv"
)

func (f *Format)SimpleUnmarshal(payload []byte) error {
	n := len(payload)
	if f.Type == typeInt{
		valString := payload[1:n-2]
		if _,err := strconv.ParseInt(string(valString),10,64);err != nil{
			return err
		}
	}
	f.Payload = payload[1:n-2]
	return nil
}

func (f *Format)BulkUnmarshal(r *bufio.Reader,payload []byte) error {
	n := len(payload)
	sizeStr := payload[1:n-2]
	size,err := strconv.ParseInt(string(sizeStr),10,64)
	if err != nil{
		return fmt.Errorf("syntax error: incorrect size input: %s\n",sizeStr)
	}

	payload,err = readValidateInput(r)
	if err != nil{
		return err
	}

	n = len(payload)
	if int(size) != n-2{
		return fmt.Errorf("syntax error: input size different from what was specified: %d - %s\n",size,payload[:n-2])
	}

	f.Payload = payload[:n-2]
	return nil
}

func (f *Format)ArrayUnmarshal(r *bufio.Reader,payload []byte) error{
	n := len(payload)
	sizeStr := payload[1:n-2]
	size,err := strconv.ParseInt(string(sizeStr),10,64)
	if err != nil{
		return fmt.Errorf("syntax error: incorrect size input: %s\n",sizeStr)
	}

	for i:=0;i<int(size);i++{
	newFormat,err := ReadCommand(r)
		if err != nil{
			return err
		}
		f.ArrayPayload = append(f.ArrayPayload, *newFormat)
	}
	return nil
}