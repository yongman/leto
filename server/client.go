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
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/yongman/go/goredis"
	"github.com/yongman/leto/store"
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
}

func newClient(app *App) *Client {
	client := &Client{
		app:   app,
		store: app.store,
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
		err = fmt.Errorf("unknown resp type")
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
		err = fmt.Errorf("unknown resp type")
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
			fmt.Println(err.Error())
			return
		} else if err != nil {
			return
		}
		err = c.handleRequest(req)
		if err != nil && err != io.EOF {
			fmt.Println(err.Error())
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

	switch c.cmd {
	case "get":
		if len(c.args) != 1 {
			c.FlushResp(fmt.Errorf("params error"))
			return nil
		}
		key := string(c.args[0])
		v, err := c.store.Get(key)
		if err != nil {
			c.FlushResp(err)
			return err
		}
		c.FlushResp(v)
	case "set":
		if len(c.args) != 2 {
			c.FlushResp(fmt.Errorf("params error"))
			return nil
		}
		key := string(c.args[0])
		value := string(c.args[1])
		err := c.store.Set(key, value)
		if err != nil {
			c.FlushResp(err)
			return err
		}
	case "del":
		if len(c.args) != 1 {
			c.FlushResp(fmt.Errorf("params error"))
			return nil
		}
		key := string(c.args[0])
		err := c.store.Delete(key)
		if err != nil {
			c.FlushResp(err)
			return err
		}
	case "ping":
		if len(c.args) != 0 {
			c.FlushResp(fmt.Errorf("params error"))
			return nil
		}
		c.FlushResp("PONG")

	default:
		c.FlushResp(fmt.Errorf("not support command"))
	}

	return nil
}
