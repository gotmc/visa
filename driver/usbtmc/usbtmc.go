// Copyright (c) 2017-2020 The visa developers. All rights reserved.
// Project site: https://github.com/gotmc/visa
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package usbtmc

import (
	"github.com/gotmc/usbtmc"
	"github.com/gotmc/visa"
)

// Driver implements the visa.Driver interface.
type Driver struct {
}

// Open opens a VISA resource given a VISA address string.
func (d Driver) Open(address string) (visa.Resource, error) {
	var c Connection
	c.ctx, _ = usbtmc.NewContext()
	dev, err := c.ctx.NewDevice(address)
	c.dev = dev
	return &c, err
}

// Connection models a USBTMC connection.
type Connection struct {
	ctx *usbtmc.Context
	dev *usbtmc.Device
}

// Read implements the Reader interface for Connection.
func (c *Connection) Read(p []byte) (n int, err error) {
	return c.dev.Read(p)
}

// Write implements the Writer interface for Connection.
func (c *Connection) Write(p []byte) (n int, err error) {
	return c.dev.Write(p)
}

// Close closes the USBTMC connection.
func (c *Connection) Close() error {
	err := c.dev.Close()
	if err != nil {
		return err
	}
	return c.ctx.Close()
}

// WriteString implements the StringWriter interface for Connection.
func (c *Connection) WriteString(s string) (int, error) {
	return c.dev.WriteString(s)
}

// Query writes the given string to the connected resource and then reads the
// return value from the VISA connection.
func (c *Connection) Query(s string) (value string, err error) {
	return c.dev.Query(s)
}

// init registers the driver with the program.
func init() {
	var driver Driver
	visa.Register(visa.USBTMC, driver)
}
