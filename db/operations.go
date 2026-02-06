package db

import (
	"fmt"
	"strconv"
)

func (db DB) Set(key, value []byte) {
	db[string(key)] = NewString(value)
}

func (db DB) Get(key []byte) ([]byte, bool) {
	val, ok := db[string(key)]
	if !ok {
		return nil, ok
	}
	return val.Val, ok
}

func (db DB) ChangeIntValue(key []byte, amount int) (int, error) {
	valString, ok := db.Get(key)
	var val int
	var err error
	if ok {
		val, err = strconv.Atoi(string(valString))
		if err != nil {
			return 0, fmt.Errorf("value is not an integer or out of range")
		}
	} else {
		val = 0
	}

	val += amount
	db.Set(key, []byte(strconv.Itoa(val)))
	return val, nil
}
