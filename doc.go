// Copyright (c) 2017-2026 The visa developers. All rights reserved.
// Project site: https://github.com/gotmc/visa
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

// Package visa implements a Virtual Instrument Software Architecture (VISA)
// resource manager for sending Standard Commands for Programmable Instruments
// (SCPI) commands or providing an interface for Interchangeable Virtual
// Instrument (IVI) drivers. It parses VISA resource address strings to create
// resources that abstract the hardware interface type (USBTMC, TCPIP, ASRL),
// using a driver registration model similar to database/sql.
//
// This package is part of the gotmc ecosystem. The visa package
// (github.com/gotmc/visa) defines a common interface for instrument
// communication across different transports (GPIB, USB, TCP/IP, serial). The
// asrl package provides the serial transport implementation. The ivi package
// (github.com/gotmc/ivi) builds on top of visa to provide standardized,
// instrument-class-specific APIs following the IVI Foundation specifications.
//
// Devices are addressed using VISA resource strings of the form:
//
//	TCPIP<boardIndex>::<hostAddress>::<port>::SOCKET
//	ASRL::<port>::<baud>::<dataflow>::INSTR
//
// For example:
//
//	TCPIP0::192.168.1.101::5025::SOCKET
//	ASRL::/dev/tty.usbserial-PX484GRU::9600::8N2::INSTR
package visa
