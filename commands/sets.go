package commands

import (
	"fmt"
	"strconv"

	"github.com/Vikrant-Dabas/redis/db"
	"github.com/Vikrant-Dabas/redis/resp"
)

func ExecuteSet(db db.DB, cmd string, input [][]byte) (*resp.Format, error) {
	switch cmd {
	case "SADD":
		return sadd(db, input)
	case "SREM":
		return srem(db, input)
	case "SPOP":
		return spopRandMem(db, input, false)
	case "SRANDMEMBER":
		return spopRandMem(db, input, true)
	case "SISMEMBER":
		return sismember(db, input)
	case "SMISMEMBER":
		return smismember(db, input)
	case "SMEMBERS":
		return smembers(db, input)
	}
	return nil, fmt.Errorf("invalid command: %s", cmd)
}

func smembers(database db.DB, input [][]byte) (*resp.Format, error) {
	if len(input) != 1 {
		return nil, fmt.Errorf("ERR wrong number of arguments for 'smembers' command")
	}

	key := input[0]

	value, ok := database.Get(key)
	if !ok {
		return &resp.Format{
			Type: resp.TypeArray,
		}, nil
	}

	if value.ValType != db.TypeSet {
		return nil, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	output := &resp.Format{
		Type: resp.TypeArray,
	}

	for member := range value.Set {
		output.ArrayPayload = append(output.ArrayPayload, resp.Format{
			Type:    resp.TypeBulk,
			Payload: []byte(member),
		})
	}

	return output, nil
}

func sismember(database db.DB, input [][]byte) (*resp.Format, error) {
	if len(input) != 2 {
		return nil, fmt.Errorf("invalid no of commands %d", len(input))
	}

	key, val := input[0], input[1]
	value, ok := database.Get(key)
	if !ok {
		return resp.FalseFormat, nil
	}
	if value.ValType != db.TypeSet {
		return nil, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	if _, ok := value.Set[string(val)]; !ok {
		return resp.FalseFormat, nil
	} else {
		return resp.TrueFormat, nil
	}
}

func smismember(database db.DB, input [][]byte) (*resp.Format, error) {
	if len(input) < 2 {
		return nil, fmt.Errorf("ERR wrong number of arguments for 'smismember' command")
	}

	key, members := input[0], input[1:]

	value, ok := database.Get(key)
	if ok && value.ValType != db.TypeSet {
		return nil, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	output := &resp.Format{
		Type: resp.TypeArray,
	}

	for _, m := range members {
		exists := false
		if ok {
			_, exists = value.Set[string(m)]
		}

		if exists {
			output.ArrayPayload = append(output.ArrayPayload, *resp.TrueFormat)
		} else {
			output.ArrayPayload = append(output.ArrayPayload, *resp.FalseFormat)
		}
	}
	return output, nil
}

func sadd(database db.DB, input [][]byte) (*resp.Format, error) {
	if len(input) < 2 {
		return nil, fmt.Errorf("invalid no of commands %d", len(input))
	}
	key, values := input[0], input[1:]
	set := &db.Value{}
	value, ok := database.Get(key)
	if !ok {
		set = db.NewSet()
	} else {
		if value.ValType != db.TypeSet {
			return nil, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
		set = value
	}
	count := 0
	for _, items := range values {
		_, ok := set.Set[string(items)]
		if !ok {
			set.Set[string(items)] = struct{}{}
			count++
		}
	}
	database.Set(key, set)
	return &resp.Format{
		Type:    resp.TypeInt,
		Payload: []byte(strconv.Itoa(count)),
	}, nil
}

func srem(database db.DB, input [][]byte) (*resp.Format, error) {
	if len(input) < 2 {
		return nil, fmt.Errorf("invalid no of commands %d", len(input))
	}
	key, values := input[0], input[1:]
	value, ok := database.Get(key)
	if !ok {
		return &resp.Format{
			Type:    resp.TypeInt,
			Payload: []byte("0"),
		}, nil
	} else {
		if value.ValType != db.TypeSet {
			return nil, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
	}
	count := 0
	for _, item := range values {
		ok, err := database.DeleteSetMember(key, item)
		if err != nil {
			return nil, err
		}
		if ok {
			count++
		}
	}
	return &resp.Format{
		Type:    resp.TypeInt,
		Payload: []byte(strconv.Itoa(count)),
	}, nil
}

func spopRandMem(database db.DB, input [][]byte, randMember bool) (*resp.Format, error) {
	if len(input) > 2 || len(input) < 1 {
		return nil, fmt.Errorf("invalid no of commands %d", len(input))
	}
	key, num := input[0], input[1:]
	value, ok := database.Get(key)
	if !ok {
		return &resp.Format{Type: resp.TypeNil}, nil
	} else {
		if value.ValType != db.TypeSet {
			return nil, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
	}
	pops := 1
	if len(num) == 1 {
		n, err := strconv.Atoi(string(num[0]))
		if err != nil || n < 0 {
			return nil, fmt.Errorf("value is not an integer or 	out of range")
		}
		pops = n
	}

	pops = min(pops, len(value.Set))

	output := &resp.Format{}
	if pops == 1 {
		output.Type = resp.TypeBulk
	} else {
		output.Type = resp.TypeArray
	}
	for member := range value.Set {
		if pops <= 0 {
			break
		}

		if output.Type == resp.TypeBulk {
			output.Payload = []byte(member)
			if !randMember {
				database.DeleteSetMember(key, []byte(member))
			}
			return output, nil
		}

		output.ArrayPayload = append(output.ArrayPayload, resp.Format{
			Type:    resp.TypeBulk,
			Payload: []byte(member),
		})

		pops--
		if !randMember {
			database.DeleteSetMember(key, []byte(member))
		}
	}

	return output, nil
}
