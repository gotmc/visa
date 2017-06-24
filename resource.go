// Copyright (c) 2017 The visa developers. All rights reserved.
// Project site: https://github.com/gotmc/ivi
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package visa

import (
	"errors"
	"fmt"
	"io"
)

type Resource interface {
	io.ReadWriteCloser
	WriterString(s string) (n int, err error)
	Query(s string) (value string, err error)
}

// NewResource creates a new Resource using the given VISA address.
func NewResource(address string) (Resource, error) {
	interfaceType, err := determineInterfaceType(address)
	if err != nil {
		return nil, errors.New("Problem determining interface type in address.")
	}
	driver, exists := drivers[interfaceType]
	if !exists {
		return nil, fmt.Errorf("The %s interface hasn't been registered.", interfaceType)
	}
	return driver.Open(address)
}
