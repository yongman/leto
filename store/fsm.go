// client.go
// Copyright (C) 2018 YanMing <yming0221@gmail.com>
//
// Distributed under terms of the MIT license.
//

package store

import (
	"encoding/json"
	"io"
	"strings"
	"sync"

	"github.com/hashicorp/raft"
)

type fsm struct {
	// just a sample in-memory storage fsm
	// this will replace with a storage such as bolt/badger/goleveldb etc.
	mu sync.Mutex
	m  map[string]string
}

func NewFSM() *fsm {
	return &fsm{
		m: make(map[string]string),
	}
}

func (f *fsm) Get(key string) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	v, ok := f.m[key]
	if !ok {
		return "", nil
	}
	return v, nil
}

func (f *fsm) Apply(l *raft.Log) interface{} {
	var c command
	if err := json.Unmarshal(l.Data, &c); err != nil {
		panic("failed to unmarshal raft log")
	}

	switch strings.ToLower(c.Op) {
	case "set":
		return f.applySet(c.Key, c.Value)
	case "delete":
		return f.applyDelete(c.Key)
	default:
		panic("command type not support")
	}
}

// TODO
func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	return &fsmSnapshot{}, nil
}

// TODO
func (f *fsm) Restore(rc io.ReadCloser) error {
	return nil
}

func (f *fsm) applySet(key, value string) interface{} {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.m[key] = value

	return nil
}

func (f *fsm) applyDelete(key string) interface{} {
	f.mu.Lock()
	defer f.mu.Unlock()

	delete(f.m, key)

	return nil
}
