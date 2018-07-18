package server

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/siddontang/goredis"
	"github.com/yongman/leto/config"
	"github.com/yongman/leto/store"
)

type App struct {
	listener net.Listener

	// wrapper and manager for db instance
	store *store.Store

	logger *log.Logger
}

// initialize an app
func NewApp(conf *config.Config) *App {
	var err error
	app := &App{}

	app.logger = log.New(os.Stderr, "[server] ", log.LstdFlags)

	app.store = store.NewStore(conf.RaftDir, conf.RaftBind)

	bootstrap := conf.Join == ""
	err = app.store.Open(bootstrap, conf.NodeID)
	if err != nil {
		app.logger.Println(err.Error())
	}

	if !bootstrap {
		// send join request to node already exists
		rc := goredis.NewClient(conf.Join, "")
		app.logger.Printf("join request send to %s", conf.Join)
		_, err := rc.Do("join", conf.RaftBind, conf.NodeID)
		if err != nil {
			app.logger.Println(err)
		}
		rc.Close()
	}

	app.listener, err = net.Listen("tcp", conf.Listen)
	app.logger.Printf("server listen in %s", conf.Listen)
	if err != nil {
		fmt.Println(err.Error())
	}

	return app
}

func (app *App) Run() {
	// accept connections
	for {
		select {
		default:
			// accept new client connect and perform
			conn, err := app.listener.Accept()
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			// handle conn
			ClientHandler(conn, app)
		}
	}
}
