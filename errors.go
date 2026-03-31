// Copyright (c) 2017-2023 The visa developers. All rights reserved.
// Project site: https://github.com/gotmc/visa
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package visa

import "errors"

var (
	ErrInvalidAddress       = errors.New("address does not match VISA format")
	ErrUnknownInterfaceType = errors.New("unknown interface type")
	ErrDriverNotRegistered  = errors.New("unregistered interface")
)
