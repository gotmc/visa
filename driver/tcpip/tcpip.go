// Copyright (c) 2017-2022 The visa developers. All rights reserved.
// Project site: https://github.com/gotmc/visa
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package tcpip

import (
	"github.com/gotmc/lxi"
	"github.com/gotmc/visa"
)

// Driver implements the visa.Driver interface for a TCPIP HW interface driver.
type Driver struct{}

// Open takes a VISA address string and returns a VISA resource.
func (d Driver) Open(address string) (visa.Resource, error) {
	return lxi.NewDevice(address)
}

// init registers the driver with the program.
func init() {
	var driver Driver
	visa.Register(visa.TCPIP, driver)
}
