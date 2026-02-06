package db

func (db DB)Set(key,value []byte){
	db[string(key)] = NewString(value)
}

func (db DB)Get(key []byte)([]byte,bool){
	val,ok := db[string(key)]
	if !ok {
		return nil,ok
	}
	return val.Val,ok
}