// Copyright (c) 2017-2022 The visa developers. All rights reserved.
// Project site: https://github.com/gotmc/visa
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

// Package usbtmc implements a VISA driver for USB Test & Measurement Class
// (USBTMC) connected instruments.
package usbtmc

import (
	"context"

	"github.com/gotmc/usbtmc"
	"github.com/gotmc/visa"
)

// Driver implements the visa.Driver interface.
type Driver struct {
}

// Open opens a VISA resource given a VISA address string.
func (d Driver) Open(_ context.Context, address string) (visa.Resource, error) {
	ctx, err := usbtmc.NewContext()
	if err != nil {
		return nil, err
	}
	dev, err := ctx.NewDevice(address)
	if err != nil {
		ctx.Close()
		return nil, err
	}
	return &Connection{ctx: ctx, dev: dev}, nil
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

// ReadContext reads from the USBTMC connection with context support for
// cancellation and deadlines.
func (c *Connection) ReadContext(ctx context.Context, p []byte) (n int, err error) {
	done := make(chan struct{})
	go func() {
		n, err = c.dev.Read(p)
		close(done)
	}()
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case <-done:
		return n, err
	}
}

// Write implements the Writer interface for Connection.
func (c *Connection) Write(p []byte) (n int, err error) {
	return c.dev.Write(p)
}

// WriteContext writes to the USBTMC connection with context support for
// cancellation and deadlines.
func (c *Connection) WriteContext(ctx context.Context, p []byte) (n int, err error) {
	done := make(chan struct{})
	go func() {
		n, err = c.dev.Write(p)
		close(done)
	}()
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case <-done:
		return n, err
	}
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

// Command sends a formatted SCPI command to the connected resource.
func (c *Connection) Command(ctx context.Context, format string, a ...any) error {
	return c.dev.Command(format, a...)
}

// Query writes the given string to the connected resource and then reads the
// return value from the VISA connection.
func (c *Connection) Query(ctx context.Context, s string) (value string, err error) {
	return c.dev.Query(s)
}

// init registers the driver with the program.
func init() {
	var driver Driver
	visa.Register(visa.USBTMC, driver)
}
