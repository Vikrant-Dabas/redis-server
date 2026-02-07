package commands

import (
	"fmt"
	"strconv"

	"github.com/Vikrant-Dabas/redis/db"
	"github.com/Vikrant-Dabas/redis/resp"
)

func ExecuteList(db db.DB, cmd string, input [][]byte) (*resp.Format, error) {
	switch cmd {
	case "LPUSH":
		return listPush(db, input, true)
	case "RPUSH":
		return listPush(db, input, false)
	case "LLEN":
		return llen(db, input)
	case "LPOP":
		return listPop(db, input, true)
	case "RPOP":
		return listPop(db, input, false)
	case "LRANGE":
		return lrange(db, input)
	}
	return nil, fmt.Errorf("invalid command: %s", cmd)
}

func lrange(database db.DB, input [][]byte) (*resp.Format, error) {
	if len(input) != 3 {
		return nil, fmt.Errorf("invalid no of commands %d", len(input))
	}

	key := input[0]
	start, err1 := strconv.Atoi(string(input[1]))
	end, err2 := strconv.Atoi(string(input[2]))
	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("value is not an integer or out of range")
	}
	output := &resp.Format{
		Type: resp.TypeArray,
	}

	value, ok := database.Get(key)
	if ok && value.ValType != db.TypeList {
		return nil, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	if !ok || value.List.Size == 0 {
		return output, nil
	}
	if end < 0 {
		end += value.List.Size
	}
	if start < 0 {
		start += value.List.Size
	}
	if end >= value.List.Size || start < 0 {
		return nil, fmt.Errorf("value is not an integer or out of range")
	}
	node := value.List.Head
	for i := 0; i < value.List.Size; i++ {
		if i >= start && i <= end {
			newBulk := &resp.Format{
				Type:    resp.TypeBulk,
				Payload: node.Val,
			}
			output.ArrayPayload = append(output.ArrayPayload, *newBulk)
		}
		node = node.Right
	}
	return output, nil
}

func llen(database db.DB, input [][]byte) (*resp.Format, error) {
	if len(input) != 1 {
		return nil, fmt.Errorf("invalid no of commands %d", len(input))
	}

	key := input[0]
	value, ok := database.Get(key)
	if ok && value.ValType != db.TypeList {
		return nil, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	if !ok {
		return &resp.Format{
			Type:    resp.TypeInt,
			Payload: []byte("0"),
		}, nil
	}
	return &resp.Format{
		Type:    resp.TypeInt,
		Payload: []byte(strconv.Itoa(value.List.Size)),
	}, nil
}

func listPush(database db.DB, input [][]byte, left bool) (*resp.Format, error) {
	if len(input) < 2 {
		return nil, fmt.Errorf("invalid no of commands %d", len(input))
	}
	list := &db.Value{}
	key, values := input[0], input[1:]
	getVal, ok := database.Get(key)
	if ok && getVal.ValType != db.TypeList {
		return nil, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	if !ok {
		list = db.NewList()
	} else {
		list = getVal
	}

	for _, node := range values {
		newNode := db.NewNode(node)
		if list.List.Head == nil {
			list.List.Head = newNode
			list.List.Tail = newNode
		} else {
			if left {
				newNode.Right = list.List.Head
				list.List.Head.Left = newNode
				list.List.Head = newNode
			} else {
				newNode.Left = list.List.Tail
				list.List.Tail.Right = newNode
				list.List.Tail = newNode
			}
		}
		list.List.Size++
	}

	database.Set(key, list)

	return &resp.Format{
		Type:    resp.TypeInt,
		Payload: []byte(strconv.Itoa(len(values))),
	}, nil
}

func listPop(database db.DB, input [][]byte, left bool) (*resp.Format, error) {
	if len(input) != 1 && len(input) != 2 {
		return nil, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	noOfPops := 1
	var err error
	key := input[0]
	if len(input) == 2 {
		noOfPops, err = strconv.Atoi(string(input[1]))
		if err != nil {
			return nil, fmt.Errorf("value is not an integer or out of range")
		}
	}

	val, ok := database.Get(key)
	if ok && val.ValType != db.TypeList {
		return nil, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	if !ok || val.List.Size == 0 {
		return &resp.Format{
			Type: resp.TypeNil,
		}, nil
	}
	if noOfPops == 1 {
		popped := listPopHelper(val, left)
		return &resp.Format{
			Type:    resp.TypeBulk,
			Payload: popped,
		}, nil
	}
	noOfPops = min(noOfPops, val.List.Size)
	output := &resp.Format{
		Type: resp.TypeArray,
	}
	for i := 0; i < noOfPops; i++ {
		arrayElement := &resp.Format{
			Type:    resp.TypeBulk,
			Payload: listPopHelper(val, left),
		}
		output.ArrayPayload = append(output.ArrayPayload, *arrayElement)
	}
	return output, nil
}

func listPopHelper(val *db.Value, left bool) []byte {
	var popped []byte
	if left {
		popped = val.List.Head.Val
		val.List.Head = val.List.Head.Right
		if val.List.Head != nil {
			val.List.Head.Left = nil
		} else {
			val.List.Tail = nil
		}
	} else {
		popped = val.List.Tail.Val
		val.List.Tail = val.List.Tail.Left
		if val.List.Tail != nil {
			val.List.Tail.Right = nil
		} else {
			val.List.Head = nil
		}
	}
	val.List.Size--
	return popped
}
