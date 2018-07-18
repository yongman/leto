//
// client.go
// Copyright (C) 2018 YanMing <yming0221@gmail.com>
//
// Distributed under terms of the MIT license.
//

package server

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/yongman/go/goredis"
	"github.com/yongman/leto/store"
)

var (
	ErrParams        = errors.New("ERR params invalid")
	ErrRespType      = errors.New("ERR resp type invalid")
	ErrCmdNotSupport = errors.New("ERR command not supported")
)

type Command struct {
	cmd  string
	args [][]byte
}

type Client struct {
	app   *App
	store *store.Store

	// request is processing
	cmd  string
	args [][]byte

	buf bytes.Buffer

	conn net.Conn

	rReader *goredis.RespReader
	rWriter *goredis.RespWriter

	logger *log.Logger
}

func newClient(app *App) *Client {
	client := &Client{
		app:    app,
		store:  app.store,
		logger: log.New(os.Stderr, "[client] ", log.LstdFlags),
	}
	return client
}

func ClientHandler(conn net.Conn, app *App) {
	c := newClient(app)

	c.conn = conn
	// connection buffer setting

	br := bufio.NewReader(conn)
	c.rReader = goredis.NewRespReader(br)

	bw := bufio.NewWriter(conn)
	c.rWriter = goredis.NewRespWriter(bw)

	go c.connHandler()
}

func (c *Client) Resp(resp interface{}) error {
	var err error = nil

	switch v := resp.(type) {
	case []interface{}:
		err = c.rWriter.WriteArray(v)
	case []byte:
		err = c.rWriter.WriteBulk(v)
	case nil:
		err = c.rWriter.WriteBulk(nil)
	case int64:
		err = c.rWriter.WriteInteger(v)
	case string:
		err = c.rWriter.WriteString(v)
	case error:
		err = c.rWriter.WriteError(v)
	default:
		err = ErrRespType
	}

	return err
}

func (c *Client) FlushResp(resp interface{}) error {
	err := c.Resp(resp)
	if err != nil {
		return err
	}
	return c.rWriter.Flush()
}

// treat string as bulk array
func (c *Client) Resp1(resp interface{}) error {
	var err error = nil

	switch v := resp.(type) {
	case []interface{}:
		err = c.rWriter.WriteArray(v)
	case []byte:
		err = c.rWriter.WriteBulk(v)
	case nil:
		err = c.rWriter.WriteBulk(nil)
	case int64:
		err = c.rWriter.WriteInteger(v)
	case string:
		err = c.rWriter.WriteBulk([]byte(v))
	case error:
		err = c.rWriter.WriteError(v)
	default:
		err = ErrRespType
	}

	return err
}
func (c *Client) connHandler() {

	defer func(c *Client) {
		c.conn.Close()
	}(c)

	for {
		c.cmd = ""
		c.args = nil

		req, err := c.rReader.ParseRequest()
		if err != nil && err != io.EOF {
			c.logger.Println(err.Error())
			return
		} else if err != nil {
			return
		}
		err = c.handleRequest(req)
		if err != nil && err != io.EOF {
			c.logger.Println(err.Error())
			return
		}
	}
}

func (c *Client) handleRequest(req [][]byte) error {
	if len(req) == 0 {
		c.cmd = ""
		c.args = nil
	} else {
		c.cmd = strings.ToLower(string(req[0]))
		c.args = req[1:]
	}

	var (
		err error
		v   string
	)

	c.logger.Printf("process %s command", c.cmd)

	switch c.cmd {
	case "get":
		if v, err = c.handleGet(); err == nil {
			c.FlushResp(v)
		}
	case "set":
		if err = c.handleSet(); err == nil {
			c.FlushResp("OK")
		}
	case "del":
		if err = c.handleDel(); err == nil {
			c.FlushResp("OK")
		}
	case "join":
		if err = c.handleJoin(); err == nil {
			c.FlushResp("OK")
		}
	case "leave":
		if err = c.handleLeave(); err == nil {
			c.FlushResp("OK")
		}
	case "ping":
		if len(c.args) != 0 {
			err = ErrParams
		}
		c.FlushResp("PONG")
		err = nil

	default:
		err = ErrCmdNotSupport
	}
	if err != nil {
		c.FlushResp(err)
	}

	return err
}
