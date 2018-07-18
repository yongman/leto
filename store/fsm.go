// client.go
// Copyright (C) 2018 YanMing <yming0221@gmail.com>
//
// Distributed under terms of the MIT license.
//

package store

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/raft"
)

type fsm struct {
	db DB

	logger *log.Logger
}

func NewFSM(path string) (*fsm, error) {
	db, err := NewBadgerDB(path, path)
	if err != nil {
		return nil, err
	}
	return &fsm{
		logger: log.New(os.Stderr, "[fsm] ", log.LstdFlags),
		db:     db,
	}, nil
}

func (f *fsm) Get(key string) (string, error) {
	v, err := f.db.Get([]byte(key))
	if err != nil {
		return "", err
	}
	return string(v), nil
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
	f.logger.Printf("apply %s to %s\n", key, value)
	err := f.db.Set([]byte(key), []byte(value))

	return err
}

func (f *fsm) applyDelete(key string) interface{} {
	_, err := f.db.Delete([]byte(key))
	return err
}

func (f *fsm) Close() error {
	f.db.Close()
	return nil
}
