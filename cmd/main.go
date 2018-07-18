//
// main.go
// Copyright (C) 2018 YanMing <yming0221@gmail.com>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/yongman/leto/config"
	"github.com/yongman/leto/server"
)

var (
	listen   string
	raftdir  string
	raftbind string
	nodeID   string
	join     string
)

func init() {
	flag.StringVar(&listen, "listen", ":5379", "server listen address")
	flag.StringVar(&raftdir, "raftdir", "./", "raft data directory")
	flag.StringVar(&raftbind, "raftbind", ":15379", "raft bus transport bind address")
	flag.StringVar(&nodeID, "id", "", "node id")
	flag.StringVar(&join, "join", "", "join to already exist cluster")
}

func main() {
	flag.Parse()

	var (
		c *config.Config
	)

	c = config.NewConfig(listen, raftdir, raftbind, nodeID, join)

	app := server.NewApp(c)

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Kill, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go app.Run()

	<-quitCh
}
