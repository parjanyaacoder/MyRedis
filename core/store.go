package core

import (
	"time"
)

type Obj struct {
	Value interface{}
	ExpiresAt int64
}

var store map[string]*Obj

func init() {
	store = make(map[string]*Obj)
}

func NewObj (value interface{}, durationMs int64) *Obj {
	var expiresAt int64 = -1
	if durationMs > 0 {
		expiresAt = time.Now().UnixMilli() + durationMs 
	}

	return &Obj{
		Value:  value,
		ExpiresAt: expiresAt,
	}
}

func Put(k string, obj *Obj) {
	store[k] = obj
}

func GET(k string) *Obj {
	return store[k]
}

