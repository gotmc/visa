// Copyright (c) 2017-2022 The visa developers. All rights reserved.
// Project site: https://github.com/gotmc/visa
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

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
