// Copyright (c) 2017-2022 The visa developers. All rights reserved.
// Project site: https://github.com/gotmc/visa
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

// Package asrl implements a VISA driver for serial (ASRL) connected
// instruments.
package asrl

import (
	"context"

	"github.com/gotmc/asrl"
	"github.com/gotmc/visa"
)

// Driver implements the visa.Driver interface for a serial HW interface driver.
type Driver struct{}

// Open takes a VISA address string and returns a VISA resource.
func (d Driver) Open(_ context.Context, address string) (visa.Resource, error) {
	dev, err := asrl.NewDevice(address)
	if err != nil {
		return nil, err
	}
	return &Connection{dev: dev}, nil
}

// Connection wraps an asrl.Device to implement the visa.Resource interface.
type Connection struct {
	dev *asrl.Device
}

// Close closes the serial connection.
func (c *Connection) Close() error {
	return c.dev.Close()
}

// Read implements the Reader interface for Connection.
func (c *Connection) Read(p []byte) (n int, err error) {
	return c.dev.Read(p)
}

// Write implements the Writer interface for Connection.
func (c *Connection) Write(p []byte) (n int, err error) {
	return c.dev.Write(p)
}

// WriteString implements the StringWriter interface for Connection.
func (c *Connection) WriteString(s string) (int, error) {
	return c.dev.WriteString(s)
}

// ReadContext reads from the serial connection with context support for
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

// WriteContext writes to the serial connection with context support for
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
	visa.Register(visa.ASRL, driver)
}
