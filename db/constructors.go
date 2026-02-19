package db

func NewString(b []byte) *Value {
	return &Value{
		ValType: TypeString,
		Val:     b,
	}
}

func NewDB() DB {
	database := make(DB)
	return database
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
