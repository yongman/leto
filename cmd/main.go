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
	listen    string
	nodeID    string
	bootstrap bool
)

func init() {
	flag.StringVar(&listen, "listen", ":5379", "server listen address")
	flag.StringVar(&nodeID, "id", "", "node id")
	flag.BoolVar(&bootstrap, "bootstrap", false, "bootstrap raft")
}

func main() {
	flag.Parse()

	var (
		c *config.Config
	)

	c = config.NewConfig(listen, nodeID, bootstrap)

	app := server.NewApp(c)

	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, os.Kill, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go app.Run()

	<-quitCh
}
