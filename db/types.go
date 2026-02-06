package db

type Value struct {
	ValType DBValueType
	Val     []byte
	Hash    *DB
	// ExpiresAt uint64    deal with this later
}

type DBValueType uint8

type DB map[string]*Value

const (
	TypeString DBValueType = iota
	TypeHash
)

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
		Hash:    &hash,
	}
}

