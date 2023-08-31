// Copyright (c) 2017-2023 The visa developers. All rights reserved.
// Project site: https://github.com/gotmc/visa
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package visa

// InterfaceType identifies the type of the HW interface: USBTMC, TCPIP, or
// ASRL (serial).
type InterfaceType int

// The available HW interface types for the VISA resource.
const (
	USBTMC InterfaceType = iota
	TCPIP
	ASRL
)

var interfaceDescription = map[InterfaceType]string{
	USBTMC: "USBTMC",
	TCPIP:  "TCP-IP",
	ASRL:   "Serial",
}

// String implements the Stringer interface for the InterfaceType.
func (i InterfaceType) String() string {
	return interfaceDescription[i]
}

type interfaceType uint16
type viVersion uint32
type viRsrc string
type viAccessMode uint32

const (
	noLock        viAccessMode = 0 // VI_NO_LOCK
	exclusiveLock viAccessMode = 1 // VI_EXCLUSIVE_LOCK
	sharedLock    viAccessMode = 2 // VI_SHARED_LOCK
	loadConfig    viAccessMode = 4 // VI_LOAD_CONFIG
)

// Hardware interface types given as VI_ATTR_INTF_TYPE in the VISA standard.
const (
	gpibInterface    interfaceType = 1
	vxiInterface     interfaceType = 2
	gpibVxiInterface interfaceType = 3
	asrlInterface    interfaceType = 4
	pxiInterface     interfaceType = 5
	tcpipInterface   interfaceType = 6
	usbInterface     interfaceType = 7
)

type resourceClass string

const (
	instrResource     resourceClass = "INSTR"
	memaccResource    resourceClass = "MEMACC"
	intfcResource     resourceClass = "INTFC"
	backplaneResource resourceClass = "BACKPLANE"
	servantResource   resourceClass = "SERVANT"
	socketResource    resourceClass = "SOCKET"
)
