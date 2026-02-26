package commands

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Vikrant-Dabas/redis/db"
	"github.com/Vikrant-Dabas/redis/resp"
)

func ExecuteUniversal(store *db.Store,cmd string, input [][]byte) (*resp.Format, error) {
	switch cmd {
	case "PING":
		return ping(input)
	case "EXPIRE":
		return expire(store,input)
	case "TTL":
		return ttl(store,input)
	case "PERSIST":
		return persist(store,input)
	}
	return nil, nil
}

func persist(store *db.Store,input [][]byte)(*resp.Format,error){
	if len(input) != 1{
		return nil, fmt.Errorf("invalid no of commands %d", len(input))
	}
	key := input[0]
	_,ok := store.GetDB(key)
	if !ok{
		return &resp.Format{
			Type: resp.TypeInt,
			Payload: []byte("0"),
		},nil
	}
	_,ok = store.GetExpDB(key)
	if !ok{
		return &resp.Format{
			Type: resp.TypeInt,
			Payload: []byte("0"),
		},nil
	}
	store.DeleteExpDB(key)
	return &resp.Format{
		Type: resp.TypeInt,
		Payload: []byte("1"),
	},nil
}


func ttl(store *db.Store,input [][]byte)(*resp.Format,error){
	if len(input) != 1{
		return nil, fmt.Errorf("invalid no of commands %d", len(input))
	}
	key := input[0]
	_,ok := store.GetDB(key)
	if !ok{
		return &resp.Format{
			Type: resp.TypeInt,
			Payload: []byte("-2"),
		},nil
	}
	timeLeft,ok := store.GetExpDB(key)
	if !ok{
		return &resp.Format{
			Type: resp.TypeInt,
			Payload: []byte("-1"),
		},nil
	}
	ttl := int(time.Until(timeLeft).Seconds())
	if ttl < 0 {
		return &resp.Format{
			Type: resp.TypeInt,
			Payload: []byte("-2"),
		}, nil
	}
	return &resp.Format{
		Type: resp.TypeInt,
		Payload: []byte(strconv.Itoa(ttl)),
	},nil
}

func expire(store *db.Store,input [][]byte)(*resp.Format,error){
	if len(input) != 2{
		return nil, fmt.Errorf("invalid no of commands %d", len(input))
	}

	key,t := input[0],input[1]
	expiry,err := strconv.Atoi(string(t))
	if err != nil{
		return nil, fmt.Errorf("value is not an integer or out of range")
	}
	_,ok := store.GetDB(key)
	if !ok{
		return &resp.Format{
			Type: resp.TypeInt,
			Payload: []byte("0"),
		},nil
	}

	if expiry == 0 {
		store.DeleteDB(key)
		return &resp.Format{Type: resp.TypeInt, Payload: []byte("1")}, nil
	}

	if expiry < 0 {
		return nil, fmt.Errorf("value is not an integer or out of range")
	}
	store.SetExpDB(key,expiry)

	return &resp.Format{
		Type: resp.TypeInt,
		Payload: []byte("1"),
	},nil 
}
func ping(input [][]byte) (*resp.Format, error) {
	size := len(input)
	if size == 0 {
		return &resp.Format{
			Type:    resp.TypeSimple,
			Payload: []byte("PONG"),
		}, nil
	} else if size == 1 {
		return &resp.Format{
			Type:    resp.TypeBulk,
			Payload: input[0],
		}, nil
	}
	return nil, fmt.Errorf("invalid PING command: %d", size)
}

