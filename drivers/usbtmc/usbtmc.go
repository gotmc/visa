// Copyright (c) 2017 The visa developers. All rights reserved.
// Project site: https://github.com/gotmc/visa
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package usbtmc

import (
	"github.com/gotmc/usbtmc"
	"github.com/gotmc/visa"
)

// Driver implements the ivi.Driver interface.
type Driver struct {
}

type Connection struct {
	ctx *usbtmc.Context
	dev *usbtmc.Device
}

func (d Driver) Open(address string) (visa.Resource, error) {
	var c Connection
	c.ctx = usbtmc.NewContext()
	dev, err := c.ctx.NewDevice(address)
	c.dev = dev
	return &c, err
}

func (c *Connection) Read(p []byte) (n int, err error) {
	return c.dev.Read(p)
}

func (c *Connection) Write(p []byte) (n int, err error) {
	return c.dev.Write(p)
}

func (c *Connection) Close() error {
	err := c.dev.Close()
	if err != nil {
		return err
	}
	return c.ctx.Close()
}

// init registers the driver with the program.
func init() {
	var driver Driver
	visa.Register(visa.Usbtmc, driver)
}
