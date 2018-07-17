// client.go
// Copyright (C) 2018 YanMing <yming0221@gmail.com>
//
// Distributed under terms of the MIT license.
//

package store

type command struct {
	Op    string `json:"op,omitempty"`
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

func NewSetCommand(key, value string) *command {
	return &command{
		Op:    "set",
		Key:   key,
		Value: value,
	}
}

func NewDeleteCommand(key string) *command {
	return &command{
		Op:  "delete",
		Key: key,
	}
}
