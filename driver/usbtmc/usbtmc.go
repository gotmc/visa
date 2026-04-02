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

// WriteContext writes to the USBTMC connection with context support for
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

// Command sends a formatted SCPI command to the connected resource. The
// upstream usbtmc library does not natively support context, so cancellation
// is handled by racing the call against ctx.Done().
func (c *Connection) Command(ctx context.Context, format string, a ...any) error {
	ch := make(chan error, 1)
	go func() {
		ch <- c.dev.Command(format, a...)
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-ch:
		return err
	}
}

// Query writes the given string to the connected resource and then reads the
// return value from the VISA connection. The upstream usbtmc library does not
// natively support context, so cancellation is handled by racing the call
// against ctx.Done().
func (c *Connection) Query(ctx context.Context, s string) (string, error) {
	type result struct {
		value string
		err   error
	}
	ch := make(chan result, 1)
	go func() {
		v, err := c.dev.Query(s)
		ch <- result{v, err}
	}()
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case r := <-ch:
		return r.value, r.err
	}
}

// init registers the driver with the program.
func init() {
	var driver Driver
	visa.Register(visa.USBTMC, driver)
}
