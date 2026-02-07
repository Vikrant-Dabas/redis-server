package commands

import (
	"fmt"
	"strconv"

	"github.com/Vikrant-Dabas/redis/db"
	"github.com/Vikrant-Dabas/redis/resp"
)

func ExecuteString(db db.DB, cmd string, input [][]byte) (*resp.Format, error) {
	switch cmd {
	case "GET":
		return get(db, input)
	case "SET":
		return set(db, input)
	case "INCR":
		return incdec(db, input, true)
	case "DECR":
		return incdec(db, input, false)
	case "INCRBY":
		return incdecby(db, input, true)
	case "DECRBY":
		return incdecby(db, input, false)
	case "MGET":
		return mget(db, input)
	case "MSET":
		return mset(db, input)
	}
	return nil, fmt.Errorf("invalid command: %s", cmd)
}

func incdec(db db.DB, input [][]byte, increase bool) (*resp.Format, error) {
	if len(input) != 1 {
		return nil, fmt.Errorf("invalid no of commands %d", len(input))
	}
	key := input[0]
	var change int
	if increase {
		change = 1
	} else {
		change = -1
	}
	val, err := db.ChangeIntValue(key, change)
	if err != nil {
		return nil, err
	}
	return &resp.Format{
		Type:    resp.TypeInt,
		Payload: []byte(strconv.Itoa(val)),
	}, nil
}

func incdecby(db db.DB, input [][]byte, increase bool) (*resp.Format, error) {
	if len(input) != 2 {
		return nil, fmt.Errorf("invalid no of commands %d", len(input))
	}
	key := input[0]
	change, err := strconv.Atoi(string(input[1]))
	if err != nil {
		return nil, err
	}
	if increase {
		change *= 1
	} else {
		change *= -1
	}
	val, err := db.ChangeIntValue(key, change)
	if err != nil {
		return nil, err
	}
	return &resp.Format{
		Type:    resp.TypeInt,
		Payload: []byte(strconv.Itoa(val)),
	}, nil
}

func get(database db.DB, input [][]byte) (*resp.Format, error) {
	if len(input) != 1 {
		return nil, fmt.Errorf("invalid no of commands %d", len(input))
	}
	key := input[0]
	output, ok := database.Get(key)
	if !ok {
		return &resp.Format{
			Type: resp.TypeNil,
		}, nil
	}
	if output.ValType != db.TypeString {
		return nil, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	return &resp.Format{
		Type:    resp.TypeBulk,
		Payload: output.Val,
	}, nil
}

func set(database db.DB, input [][]byte) (*resp.Format, error) {
	if len(input) != 2 {
		return nil, fmt.Errorf("invalid no of commands %d", len(input))
	}
	key, value := input[0], input[1]
	database.Set(key, db.NewString(value))
	return &resp.Format{
		Type:    resp.TypeSimple,
		Payload: []byte("OK"),
	}, nil
}

func mset(database db.DB, input [][]byte) (*resp.Format, error) {
	n := len(input)
	if n == 0 || n%2 == 1 {
		return nil, fmt.Errorf("wrong number of commands for mset command")
	}
	for i := 0; i < n; i += 2 {
		key, value := input[i], input[i+1]
		val := db.NewString(value)
		database.Set(key, val)
	}

	return &resp.Format{
		Type:    resp.TypeSimple,
		Payload: []byte("OK"),
	}, nil
}

func mget(database db.DB, input [][]byte) (*resp.Format, error) {
	output := &resp.Format{
		Type: resp.TypeArray,
	}
	for _, i := range input {
		val, ok := database.Get(i)
		if ok && val.ValType != db.TypeString {
			return nil, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		if ok {
			output.ArrayPayload = append(output.ArrayPayload, resp.Format{
				Type:    resp.TypeBulk,
				Payload: val.Val,
			})
		} else {
			output.ArrayPayload = append(output.ArrayPayload, resp.Format{
				Type: resp.TypeNil,
			})
		}
	}
	return output, nil
}
