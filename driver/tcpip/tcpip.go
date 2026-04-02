// Copyright (c) 2017-2026 The visa developers. All rights reserved.
// Project site: https://github.com/gotmc/visa
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

// Package tcpip implements a VISA driver for TCP/IP connected instruments
// using the LXI protocol.
package tcpip

import (
	"context"

	"github.com/gotmc/lxi"
	"github.com/gotmc/visa"
)

var _ visa.Resource = (*Connection)(nil)

// Driver implements the visa.Driver interface for a TCPIP HW interface driver.
type Driver struct{}

// Open takes a VISA address string and returns a VISA resource.
func (d Driver) Open(ctx context.Context, address string) (visa.Resource, error) {
	dev, err := lxi.NewDevice(ctx, address)
	if err != nil {
		return nil, err
	}
	return &Connection{dev: dev}, nil
}

// Connection wraps an lxi.Device to implement the visa.Resource interface.
type Connection struct {
	dev *lxi.Device
}

// Close closes the TCPIP connection.
func (c *Connection) Close() error {
	return c.dev.Close()
}

// Read implements the Reader interface for Connection.
func (c *Connection) Read(p []byte) (n int, err error) {
	return c.dev.Read(p)
}

// ReadContext reads from the TCPIP connection with context support for
// cancellation and deadlines. If the context is cancelled, the underlying
// read may still complete in the background.
func (c *Connection) ReadContext(ctx context.Context, p []byte) (int, error) {
	type result struct {
		n   int
		err error
	}
	ch := make(chan result, 1)
	go func() {
		n, err := c.dev.Read(p)
		ch <- result{n, err}
	}()
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case r := <-ch:
		return r.n, r.err
	}
}

// Write implements the Writer interface for Connection.
func (c *Connection) Write(p []byte) (n int, err error) {
	return c.dev.Write(p)
}

// WriteContext writes to the TCPIP connection with context support for
// cancellation and deadlines. If the context is cancelled, the underlying
// write may still complete in the background.
func (c *Connection) WriteContext(ctx context.Context, p []byte) (int, error) {
	type result struct {
		n   int
		err error
	}
	ch := make(chan result, 1)
	go func() {
		n, err := c.dev.Write(p)
		ch <- result{n, err}
	}()
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case r := <-ch:
		return r.n, r.err
	}
}

// WriteString implements the StringWriter interface for Connection.
func (c *Connection) WriteString(s string) (int, error) {
	return c.dev.WriteString(s)
}

// Command sends a formatted SCPI command to the connected resource.
func (c *Connection) Command(ctx context.Context, format string, a ...any) error {
	return c.dev.Command(ctx, format, a...)
}

// Query writes the given string to the connected resource and then reads the
// return value from the VISA connection.
func (c *Connection) Query(ctx context.Context, s string) (string, error) {
	return c.dev.Query(ctx, s)
}

// init registers the driver with the program.
func init() {
	var driver Driver
	visa.Register(visa.TCPIP, driver)
}
