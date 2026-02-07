package db

func NewString(b []byte) *Value {
	return &Value{
		ValType: TypeString,
		Val:     b,
	}
}

func NewHash() *Value {
	hash := make(DB)
	return &Value{
		ValType: TypeHash,
		Hash:    hash,
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
