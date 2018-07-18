//
// db_badger.go
// Copyright (C) 2018 YanMing <yming0221@gmail.com>
//
// Distributed under terms of the MIT license.
//

package store

import "github.com/dgraph-io/badger"

type BadgerDB struct {
	dir      string
	valueDir string
	db       *badger.DB
}

func NewBadgerDB(dir, valueDir string) (*BadgerDB, error) {
	opts := badger.DefaultOptions
	opts.Dir = dir
	opts.ValueDir = valueDir
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &BadgerDB{
		dir:      dir,
		valueDir: valueDir,
		db:       db,
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

func (b *BadgerDB) Close() {
	b.db.Close()
}
