package server

import (
	"fmt"
	"net"

	"github.com/yongman/leto/config"
	"github.com/yongman/leto/store"
)

type App struct {
	listener net.Listener

	// wrapper and manager for db instance
	store *store.Store
}

// initialize an app
func NewApp(conf *config.Config) *App {
	var err error
	app := &App{}

	app.store = store.NewStore()

	err = app.store.Open(conf.Bootstrap, conf.NodeID)
	if err != nil {
		fmt.Println(err.Error())
	}

	app.listener, err = net.Listen("tcp", conf.Listen)
	fmt.Println("server listen in %s", conf.Listen)
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
