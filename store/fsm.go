// fsm.go
// Copyright (C) 2018 YanMing <yming0221@gmail.com>
//
// Distributed under terms of the MIT license.
//

package store

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
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

// generate FSMSnapshot
// Snapshot will be called duiring make snapshot
func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	f.logger.Printf("Generate FSMSnapshot")
	return &fsmSnapshot{
		db:     f.db,
		logger: log.New(os.Stderr, "[fsmSnapshot] ", log.LstdFlags),
	}, nil
}

// restore from FSMSnapshot
// TODO
func (f *fsm) Restore(rc io.ReadCloser) error {
	f.logger.Printf("Restore snapshot from FSMSnapshot")
	defer rc.Close()

	var (
		readBuf  []byte
		protoBuf *proto.Buffer
		err      error
		keyCount int = 0
	)
	// decode message from protobuf
	f.logger.Printf("Read all data")
	if readBuf, err = ioutil.ReadAll(rc); err != nil {
		// read done completely
		f.logger.Printf("Snapshot restore failed")
		return err
	}

	protoBuf = proto.NewBuffer(readBuf)

	f.logger.Printf("new protoBuf length %d bytes", len(protoBuf.Bytes()))

	// decode messages from 1M block file
	// the last message could decode failed with io.ErrUnexpectedEOF
	for {
		item := &ProtoKVItem{}
		if err = protoBuf.DecodeMessage(item); err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			f.logger.Printf("DecodeMessage failed %v", err)
			return err
		}
		// apply item to store
		f.logger.Printf("Set key %v to %v count: %d", item.Key, item.Value, keyCount)
		err = f.db.Set(item.Key, item.Value)
		if err != nil {
			f.logger.Printf("Snapshot load failed %v", err)
			return err
		}
		keyCount = keyCount + 1
	}

	f.logger.Printf("Restore total %d keys", keyCount)

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
