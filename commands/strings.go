package commands

import (
	"fmt"

	"github.com/Vikrant-Dabas/redis/db"
	"github.com/Vikrant-Dabas/redis/resp"
)

func ExecuteString(db db.DB,cmd string,input [][]byte)(*resp.Format,error){
	switch cmd{
	case "GET":
		return get(db,input)
	case "SET":
		return set(db,input)
	}
	return nil,fmt.Errorf("invalid command: %s",cmd)
}

func get(db db.DB,input[][]byte)(*resp.Format,error){
	if len(input) != 1{
		return nil,fmt.Errorf("invalid no of commands %d",len(input))
	}
	key := input[0]
	output,ok := db.Get(key)
	if !ok{
		return nil,nil
	}
	return &resp.Format{
		Type: resp.TypeBulk,
		Payload: output,
	},nil
}

func set(db db.DB,input [][]byte)(*resp.Format,error){
	if len(input) != 2{
		return nil,fmt.Errorf("invalid no of commands %d",len(input))
	}
	key,value := input[0],input[1]
	db.Set(key,value)
	return &resp.Format{
		Type: resp.TypeSimple,
		Payload: []byte("OK"),
	},nil
}