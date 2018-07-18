package config

type Config struct {
	Desc     string
	Listen   string
	RaftDir  string
	RaftBind string
	Join     string
	NodeID   string
}

func NewConfig(listen, raftdir, raftbind, nodeid, join string) *Config {
	return &Config{
		Listen:   listen,
		RaftDir:  raftdir,
		RaftBind: raftbind,
		NodeID:   nodeid,
		Join:     join,
		Desc:     "",
	}
}
