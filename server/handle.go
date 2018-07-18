//
// handle.go
// Copyright (C) 2018 YanMing <yming0221@gmail.com>
//
// Distributed under terms of the MIT license.
//

package server

func (c *Client) handleJoin() error {
	if len(c.args) != 2 {
		return ErrParams
	}
	raftAddr := string(c.args[0])
	nodeID := string(c.args[1])

	return c.store.Join(nodeID, raftAddr)
}

func (c *Client) handleLeave() error {
	if len(c.args) != 1 {
		return ErrParams
	}

	nodeID := string(c.args[0])

	return c.store.Leave(nodeID)
}

func (c *Client) handleSet() error {
	if len(c.args) != 2 {
		return ErrParams
	}
	key := string(c.args[0])
	value := string(c.args[1])
	err := c.store.Set(key, value)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) handleGet() (string, error) {
	if len(c.args) != 1 {
		return "", ErrParams
	}
	key := string(c.args[0])
	v, err := c.store.Get(key)
	if err != nil {
		return "", err
	}

	return v, nil
}

func (c *Client) handleDel() error {
	if len(c.args) != 1 {
		return ErrParams
	}
	key := string(c.args[0])
	err := c.store.Delete(key)
	if err != nil {
		return err
	}
	return nil
}
