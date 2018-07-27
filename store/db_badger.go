//
// db_badger.go
// Copyright (C) 2018 YanMing <yming0221@gmail.com>
//
// Distributed under terms of the MIT license.
//

package store

import (
	"log"
	"os"

	"github.com/dgraph-io/badger"
)

type BadgerDB struct {
	dir      string
	valueDir string
	db       *badger.DB
	logger   *log.Logger
}

type KVItem struct {
	key   []byte
	value []byte
	err   error
}

func (i *KVItem) IsFinished() bool {
	return i.err == ErrIterFinished
}

func NewBadgerDB(dir, valueDir string) (*BadgerDB, error) {
	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = valueDir
	opts.SyncWrites = false
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &BadgerDB{
		dir:      dir,
		valueDir: valueDir,
		db:       db,
		logger:   log.New(os.Stderr, "[db_badger] ", log.LstdFlags),
	}, nil
}

func (b *BadgerDB) Get(key []byte) ([]byte, error) {
	value := []byte{}

	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil
			}
			return err
		}
		v, err := item.Value()
		if err != nil {
			return err
		}
		value = append([]byte{}, v...)
		return nil
	})
	if err != nil {
		return nil, err
	} else {
		return value, nil
	}
}

func (b *BadgerDB) Set(key, value []byte) error {
	err := b.db.Update(func(txn *badger.Txn) error {
		if err1 := txn.Set(key, value); err1 != nil {
			return err1
		}
		return nil
	})
	return err
}

func (b *BadgerDB) Delete(key []byte) (bool, error) {
	err := b.db.Update(func(txn *badger.Txn) error {
		if err1 := txn.Delete(key); err1 != nil {
			return err1
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (b *BadgerDB) SnapshotItems() <-chan DataItem {
	// create a new channel
	ch := make(chan DataItem, 1024)

	// generate items from snapshot to channel
	go b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		keyCount := 0
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			v, err := item.Value()

			kvi := &KVItem{
				key:   append([]byte{}, k...),
				value: append([]byte{}, v...),
				err:   err,
			}

			// write kvitem to channel with last error
			ch <- kvi
			keyCount = keyCount + 1

			if err != nil {
				return err
			}
		}

		// just use nil kvitem to mark the end
		kvi := &KVItem{
			key:   nil,
			value: nil,
			err:   ErrIterFinished,
		}
		ch <- kvi

		b.logger.Printf("Snapshot total %d keys", keyCount)

		return nil
	})

	// return channel to persist
	return ch
}

func (b *BadgerDB) Close() {
	b.db.Close()
}
