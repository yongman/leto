// client.go
// Copyright (C) 2018 YanMing <yming0221@gmail.com>
//
// Distributed under terms of the MIT license.
//

package store

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

type Store struct {
	RaftDir  string
	RaftBind string

	raft *raft.Raft // The consensus mechanism
	fsm  *fsm

	logger *log.Logger
}

var (
	ErrNotLeader error = errors.New("not leader")
)

func NewStore() *Store {
	return &Store{
		logger: log.New(os.Stderr, "[store] ", log.LstdFlags),
		fsm: &fsm{
			m: make(map[string]string),
		},
	}
}

func (s *Store) Open(bootstrap bool, localID string) error {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(localID)

	addr, err := net.ResolveTCPAddr("tcp", s.RaftBind)
	if err != nil {
		return err
	}

	transport, err := raft.NewTCPTransport(s.RaftBind, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return err
	}

	ss, err := raft.NewFileSnapshotStore(s.RaftDir, 2, os.Stderr)
	if err != nil {
		return err
	}

	// boltDB implement log store and stable store interface
	boltDB, err := raftboltdb.NewBoltStore(filepath.Join(s.RaftDir, "raft.db"))
	if err != nil {
		return err
	}

	// raft system
	r, err := raft.NewRaft(config, s.fsm, boltDB, boltDB, ss, transport)
	if err != nil {
		return err
	}
	s.raft = r

	if bootstrap {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		s.raft.BootstrapCluster(configuration)
	}
	return nil
}

func (s *Store) Get(key string) (string, error) {
	return s.fsm.Get(key)
}

func (s *Store) Set(key, value string) error {
	if s.raft.State() != raft.Leader {
		return ErrNotLeader
	}

	c := NewSetCommand(key, value)

	msg, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f := s.raft.Apply(msg, 10*time.Second)

	return f.Error()
}

func (s *Store) Delete(key string) error {
	if s.raft.State() != raft.Leader {
		return ErrNotLeader
	}

	c := NewDeleteCommand(key)

	msg, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f := s.raft.Apply(msg, 10*time.Second)

	return f.Error()
}

func (s *Store) Join(nodeID, addr string) error {
	s.logger.Printf("received join request for remote node %s, addr %s", nodeID, addr)

	cf := s.raft.GetConfiguration()

	if err := cf.Error(); err != nil {
		s.logger.Printf("failed to get raft configuration")
		return err
	}

	for _, server := range cf.Configuration().Servers {
		if server.ID == raft.ServerID(nodeID) {
			s.logger.Printf("node %s already joined raft cluster", nodeID)
			return nil
		}
	}

	f := s.raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, 0)
	if err := f.Error(); err != nil {
		return err
	}

	s.logger.Printf("node %s at %s joined successfully", nodeID, addr)

	return nil
}

func (s *Store) Leave(nodeID string) error {
	s.logger.Printf("received leave request for remote node %s", nodeID)

	cf := s.raft.GetConfiguration()

	if err := cf.Error(); err != nil {
		s.logger.Printf("failed to get raft configuration")
		return err
	}

	for _, server := range cf.Configuration().Servers {
		if server.ID == raft.ServerID(nodeID) {
			f := s.raft.RemoveServer(server.ID, 0, 0)
			if err := f.Error(); err != nil {
				s.logger.Printf("failed to remove server %s", nodeID)
				return err
			}

			s.logger.Printf("node %s leaved successfully", nodeID)
			return nil
		}
	}

	s.logger.Printf("node %s not exists in raft group", nodeID)

	return nil
}
