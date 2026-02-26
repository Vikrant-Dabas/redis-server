package db

import "time"

type Value struct {
	ValType   DBValueType
	Val       []byte
	List      *List
	Set       Set
}

type DBValueType uint8

type ExpDB map[string]time.Time
type DB map[string]*Value
type Store struct{
	DB DB
	ExpDB ExpDB
}
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
