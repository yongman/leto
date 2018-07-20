// fsm_snapshot.go
// Copyright (C) 2018 YanMing <yming0221@gmail.com>
//
// Distributed under terms of the MIT license.
//

package store

import (
	"github.com/golang/protobuf/proto"
	"github.com/hashicorp/raft"
)

type fsmSnapshot struct {
	db DB
}

// Persist data in specific type
// kv item serialize in google protubuf
func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	defer sink.Close()

	ch := f.db.SnapshotItems()

	buff := proto.NewBuffer([]byte{})

	// read kv item from channel
	for {
		dataItem := <-ch
		item := dataItem.(*KVItem)

		if item.IsFinished() {
			return nil
		}

		// create new protobuf item
		protoKVItem := &ProtoKVItem{
			Key:   item.key,
			Value: item.value,
		}

		// encode message
		buff.EncodeMessage(protoKVItem)
		if _, err := sink.Write(buff.Bytes()); err != nil {
			return err
		}

	}
	return nil
}

func (f *fsmSnapshot) Release() {

}
