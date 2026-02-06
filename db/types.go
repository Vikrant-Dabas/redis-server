package db

type Value struct {
	ValType DBValueType
	Val     []byte
	Hash    DB
	// ExpiresAt uint64    deal with this later
}

type DBValueType uint8

type DB map[string]*Value

const (
	TypeString DBValueType = iota
	TypeHash
)
