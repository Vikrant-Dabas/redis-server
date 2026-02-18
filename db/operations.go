package db

import (
	"fmt"
	"strconv"
)

func (db DB) Set(key []byte, value *Value) {
	db[string(key)] = value
}

func (db DB) Get(key []byte) (*Value, bool) {
	val, ok := db[string(key)]
	if !ok {
		return nil, ok
	}
	return val, ok
}

func (db DB) Delete(key []byte) {
	delete(db, string(key))
}

func (db DB) DeleteSetMember(key, member []byte) (bool, error) {
	val, ok := db[string(key)]
	if !ok {
		return false, nil
	}

	if val.ValType != TypeSet {
		return false, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	_, existed := val.Set[string(member)]
	delete(val.Set, string(member))

	if len(val.Set) == 0 {
		delete(db, string(key))
	}

	return existed, nil
}

func (db DB) ChangeIntValue(key []byte, amount int) (int, error) {
	valString, ok := db.Get(key)
	var val int
	var err error
	if ok {
		val, err = strconv.Atoi(string(valString.Val))
		if err != nil {
			return 0, fmt.Errorf("value is not an integer or out of range")
		}
	} else {
		val = 0
	}

	val += amount
	db.Set(key, NewString([]byte(strconv.Itoa(val))))
	return val, nil
}
