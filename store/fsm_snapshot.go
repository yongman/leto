// fsm_snapshot.go
// Copyright (C) 2018 YanMing <yming0221@gmail.com>
//
// Distributed under terms of the MIT license.
//

package store

import (
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/raft"
)

type fsmSnapshot struct {
	db     DB
	logger *log.Logger
}

// Persist data in specific type
// kv item serialize in google protubuf
func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	f.logger.Printf("Persist action in fsmSnapshot")
	defer sink.Close()

	ch := f.db.SnapshotItems()

	keyCount := 0

	// read kv item from channel
	for {
		buff := proto.NewBuffer([]byte{})

		dataItem := <-ch
		item := dataItem.(*KVItem)

		if item.IsFinished() {
			break
		}

		// create new protobuf item
		protoKVItem := &ProtoKVItem{
			Key:   item.key,
			Value: item.value,
		}

		keyCount = keyCount + 1

		// encode message
		buff.EncodeMessage(protoKVItem)

		if _, err := sink.Write(buff.Bytes()); err != nil {
			return err
		}
	}
	f.logger.Printf("Persist total %d keys", keyCount)

	return nil
}

func (f *fsmSnapshot) Release() {
	f.logger.Printf("Release action in fsmSnapshot")
}
