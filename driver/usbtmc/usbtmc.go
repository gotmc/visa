// Copyright (c) 2017-2026 The visa developers. All rights reserved.
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

var _ visa.Resource = (*Connection)(nil)

// Driver implements the visa.Driver interface.
type Driver struct {
}

// Open opens a VISA resource given a VISA address string. The upstream usbtmc
// library does not natively support context for device creation, so
// cancellation is handled by racing the call against ctx.Done().
func (d Driver) Open(ctx context.Context, address string) (visa.Resource, error) {
	type result struct {
		conn *Connection
		err  error
	}
	ch := make(chan result, 1)
	go func() {
		usbCtx, err := usbtmc.NewContext()
		if err != nil {
			ch <- result{nil, err}
			return
		}
		dev, err := usbCtx.NewDevice(address)
		if err != nil {
			usbCtx.Close()
			ch <- result{nil, err}
			return
		}
		ch <- result{&Connection{ctx: usbCtx, dev: dev}, nil}
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r := <-ch:
		return r.conn, r.err
	}
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
func (c *Connection) ReadContext(ctx context.Context, p []byte) (int, error) {
	return c.dev.ReadContext(ctx, p)
}

// Write implements the Writer interface for Connection.
func (c *Connection) Write(p []byte) (n int, err error) {
	return c.dev.Write(p)
}

// WriteContext writes to the USBTMC connection with context support for
// cancellation and deadlines.
func (c *Connection) WriteContext(ctx context.Context, p []byte) (int, error) {
	return c.dev.WriteContext(ctx, p)
}

// Close closes the USBTMC connection.
func (c *Connection) Close() error {
	devErr := c.dev.Close()
	ctxErr := c.ctx.Close()
	if devErr != nil {
		return devErr
	}
	return ctxErr
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
	visa.Register(visa.USBTMC, driver)
}
