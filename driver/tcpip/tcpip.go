// Project site: https://github.com/gotmc/visa
// Copyright (c) 2017 The visa developers. All rights reserved.
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package tcpip

import (
	"fmt"
	"net"

	"github.com/gotmc/visa"
)

// Driver implements the ivi.Driver interface for a TCPIP HW interface driver.
type Driver struct{}

// Open takes a VISA address string and returns a VISA resource.
func (d Driver) Open(address string) (visa.Resource, error) {
	var c Connection
	tcpAddress, err := getTCPAddress(address)
	if err != nil {
		return nil, fmt.Errorf("%s is not a TCPIP VISA resource address: %s", address, err)
	}
	c.conn, err = net.Dial("tcp", tcpAddress)
	if err != nil {
		return nil, fmt.Errorf("Problem connecting to TCP instrument at %s: %s", tcpAddress, err)
	}
	return &c, nil
}

// Connection models a network connection.
type Connection struct {
	conn net.Conn
}

// Read implements the reader interface for Connection.
func (c *Connection) Read(p []byte) (n int, err error) {
	return 0, nil
}

// Write implements the writer interface for Connection.
func (c *Connection) Write(p []byte) (n int, err error) {
	return 0, nil
}

// Close implements the closer interface for Connection.
func (c *Connection) Close() error {
	return c.conn.Close()
}

// WriteString implements the StringWriter interface for Connection.
func (c *Connection) WriteString(s string) (int, error) {
	return 0, nil
}

// Query implements the Querier interface for Connection.
func (c *Connection) Query(s string) (value string, err error) {
	return "foo", nil
}

func getTCPAddress(address string) (string, error) {
	return "127.0.0.1:5025", nil
}

// init registers the driver with the program.
func init() {
	var driver Driver
	visa.Register(visa.TCPIP, driver)
}
