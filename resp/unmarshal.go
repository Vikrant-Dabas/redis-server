package resp

import (
	"bufio"
	"fmt"
	"strconv"
)

func SimpleUnmarshal(format *Format,payload []byte) error {
	n := len(payload)
	if format.Type == typeInt{
		valString := payload[1:n-2]
		if _,err := strconv.ParseInt(string(valString),10,64);err != nil{
			return err
		}
	}
	format.Size = 1
	format.Payload = append(format.Payload, payload[1:n-2])
	return nil
}

func BulkUnmarshal(r *bufio.Reader, format *Format,payload []byte) error {
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

	format.Size = 1
	format.Payload = append(format.Payload, payload[:n-2])
	return nil
}

func ArrayUnmarshal(r *bufio.Reader,format *Format,payload []byte) error{
	n := len(payload)
	sizeStr := payload[1:n-2]
	size,err := strconv.ParseInt(string(sizeStr),10,64)
	if err != nil{
		return fmt.Errorf("syntax error: incorrect size input: %s\n",sizeStr)
	}
	format.Size = int(size)

	for i:=0;i<format.Size;i++{
		newFormat,err := ReadCommand(r)
		if err != nil{
			return err
		}
		format.ArrayPayload = append(format.ArrayPayload, *newFormat)
	}
	return nil
}