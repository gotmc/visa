// Copyright (c) 2017-2023 The visa developers. All rights reserved.
// Project site: https://github.com/gotmc/visa
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package visa

import (
	"fmt"
	"log"
)

// A map of registered matchers for searching.
var drivers = make(map[InterfaceType]Driver)

// Driver defines the behavior required by types that want
// to implement a new search type.
type Driver interface {
	Open(address string) (Resource, error)
}

// Register is called to register a driver for use by the program.
func Register(interfaceType InterfaceType, driver Driver) {
	if _, exists := drivers[interfaceType]; exists {
		// TODO(mdr): Should we log.Fatalln, or should we just re-register the
		// newer driver?
		log.Fatalln(interfaceType, "Driver already registered")
	}
	drivers[interfaceType] = driver
}

// Resource is the interface that defines a VISA resource.
type Resource interface {
	Close() error
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	WriteString(s string) (n int, err error)
	Query(s string) (value string, err error)
}

// NewResource creates a new Resource using the given VISA address.
func NewResource(address string) (Resource, error) {
	interfaceType, err := determineInterfaceType(address)
	if err != nil {
		return nil, err
	}
	driver, exists := drivers[interfaceType]
	if !exists {
		return nil, fmt.Errorf("unregistered interface: %s", interfaceType)
	}
	return driver.Open(address)
}
