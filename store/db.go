//
// db.go
// Copyright (C) 2018 YanMing <yming0221@gmail.com>
//
// Distributed under terms of the MIT license.
//

package store

type DB interface {
	Get(key []byte) ([]byte, error)
	Set(key, value []byte) error
	Delete(key []byte) (bool, error)

	SnapshotItems() <-chan DataItem

	Close()
}

type DataItem interface{}
