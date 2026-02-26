package db

func NewString(b []byte) *Value {
	return &Value{
		ValType: TypeString,
		Val:     b,
	}
}

func NewExpDB() ExpDB{
	database := make(ExpDB)
	return database
}

func NewDB() DB {
	database := make(DB)
	return database
}

func NewStore() *Store{
	db,expDb := make(DB),make(ExpDB)
	return &Store{
		DB: db,
		ExpDB: expDb,
	}
}

func NewList() *Value {
	list := &List{}
	return &Value{
		ValType: TypeList,
		List:    list,
	}
}

func NewNode(value []byte) *Node {
	return &Node{
		Val: value,
	}
}

func NewSet() *Value {
	set := make(map[string]struct{})
	return &Value{
		ValType: TypeSet,
		Set:     set,
	}
}
