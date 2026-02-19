package db

type Value struct {
	ValType DBValueType
	Val     []byte
	List    *List
	Set     Set
	// ExpiresAt uint64    deal with this later
}

type DBValueType uint8

type DB map[string]*Value
type Set map[string]struct{}

const (
	TypeString DBValueType = iota
	TypeList
	TypeSet
)

type List struct {
	Head *Node
	Tail *Node
	Size int
}

type Node struct {
	Left, Right *Node
	Val         []byte
}
