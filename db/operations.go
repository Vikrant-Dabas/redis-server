package db

import (
	"fmt"
	"strconv"
	"time"
)

func (store *Store) SetDB(key []byte, value *Value) {
	store.DB[string(key)] = value
}

func (store *Store) GetDB(key []byte) (*Value, bool) {
	expiry, ok := store.GetExpDB(key)
	if ok{
		if time.Now().After(expiry){
			store.DeleteExpDB(key)
			store.DeleteDB(key)
			return nil,false
		}
	}
	val, ok := store.DB[string(key)]
	if !ok {
		return nil, ok
	}
	return val, ok
}

func (store *Store) DeleteDB(key []byte) {
	delete(store.DB, string(key))
	delete(store.ExpDB,string(key))
}

func (store *Store) DeleteSetMember(key, member []byte) (bool, error) {
	val, ok := store.DB[string(key)]
	if !ok {
		return false, nil
	}

	if val.ValType != TypeSet {
		return false, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	_, existed := val.Set[string(member)]
	delete(val.Set, string(member))

	if len(val.Set) == 0 {
		delete(store.DB, string(key))
		delete(store.ExpDB, string(key))
	}

	return existed, nil
}

func (store *Store) ChangeIntValue(key []byte, amount int) (int, error) {
	valString, ok := store.GetDB(key)
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
	store.SetDB(key, NewString([]byte(strconv.Itoa(val))))
	return val, nil
}

func (store *Store) SetExpDB(key []byte,value int){
	store.ExpDB[string(key)] = time.Now().Add(time.Duration(value)*time.Second)
}

func (store *Store) GetExpDB(key []byte)(time.Time,bool){
	val, ok := store.ExpDB[string(key)]
	if !ok {
		return time.Time{}, ok
	}
	return val, ok
}

func (store *Store) DeleteExpDB(key []byte) {
	delete(store.ExpDB, string(key))
}