package resp

import (
	"bufio"
	"fmt"
	"strconv"
)

func (f *Format) SimpleUnmarshal(payload []byte) error {
	n := len(payload)
	if f.Type == TypeInt {
		valString := payload[1 : n-2]
		if _, err := strconv.ParseInt(string(valString), 10, 64); err != nil {
			return err
		}
	}
	f.Payload = payload[1 : n-2]
	return nil
}

func (f *Format) BulkUnmarshal(r *bufio.Reader, payload []byte) error {
	n := len(payload)
	sizeStr := payload[1 : n-2]
	size, err := strconv.ParseInt(string(sizeStr), 10, 64)
	if err != nil {
		return fmt.Errorf("syntax error: incorrect size input: %s\n", sizeStr)
	}

	payload, err = readValidateInput(r)
	if err != nil {
		return err
	}

	n = len(payload)
	if int(size) != n-2 {
		return fmt.Errorf("syntax error: input size different from what was specified: %d - %s\n", size, payload[:n-2])
	}

	f.Payload = payload[:n-2]
	return nil
}

func (f *Format) ArrayUnmarshal(r *bufio.Reader, payload []byte) error {
	n := len(payload)
	sizeStr := payload[1 : n-2]
	size, err := strconv.ParseInt(string(sizeStr), 10, 64)
	if err != nil {
		return fmt.Errorf("syntax error: incorrect size input: %s\n", sizeStr)
	}

	for i := 0; i < int(size); i++ {
		newFormat, err := ReadCommand(r)
		if err != nil {
			return err
		}
		f.ArrayPayload = append(f.ArrayPayload, *newFormat)
	}
	return nil
}

func contains(b []byte) (byte, bool) {
	for _, i := range b {
		for _, ch := range AllTypes {
			if i == ch {
				return ch, true
			}
		}
	}
	return 0, false
}

func validTerminator(b []byte) bool {
	n := len(b)
	if n < 2 {
		return false
	}
	if b[n-2] != '\r' || b[n-1] != '\n' {
		return false
	}

	for _, ch := range b[:n-2] {
		if ch == '\r' || ch == '\n' {
			return false
		}
	}
	return true
}

func terminatorInBetween(b []byte) bool {
	for _, ch := range b {
		if ch == '\r' || ch == '\n' {
			return true
		}
	}
	return false
}

func readValidateInput(r *bufio.Reader) ([]byte, error) {
	payload, err := r.ReadBytes('\n')
	// fmt.Printf("\033[94m%s-%d\033[0m\n", payload, len(payload))
	if err != nil {
		return nil, err
	}
	if !validTerminator(payload) {
		return nil, fmt.Errorf("syntax error: terminator: %d-%s\n", len(payload), payload)
	}
	if _, ok := contains([]byte{payload[0]}); ok {
		if ch, ok := contains(payload[1:]); ok {
			return nil, fmt.Errorf("syntax error: unexpted message type - %c", ch)
		}
	} else if ch, ok := contains(payload); ok {
		return nil, fmt.Errorf("syntax error: unexpted message type - %c", ch)
	}
	return payload, nil
}
