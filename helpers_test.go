// Copyright (c) 2017-2026 The visa developers. All rights reserved.
// Project site: https://github.com/gotmc/visa
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package visa

import (
	"testing"
)

func TestDetermineInterfaceType(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    InterfaceType
		wantErr bool
	}{
		{
			name:    "TCPIP with board index",
			address: "TCPIP0::192.168.1.100::inst0::INSTR",
			want:    TCPIP,
		},
		{
			name:    "TCPIP without board index",
			address: "TCPIP::192.168.1.100::inst0::INSTR",
			want:    TCPIP,
		},
		{
			name:    "TCPIP socket",
			address: "TCPIP0::192.168.1.100::5025::SOCKET",
			want:    TCPIP,
		},
		{
			name:    "USB with board index",
			address: "USB0::0x0957::0x0407::MY44012345::INSTR",
			want:    USBTMC,
		},
		{
			name:    "USB without board index",
			address: "USB::0x0957::0x0407::MY44012345::INSTR",
			want:    USBTMC,
		},
		{
			name:    "ASRL with board index",
			address: "ASRL1::/dev/ttyUSB0::INSTR",
			want:    ASRL,
		},
		{
			name:    "ASRL without board index",
			address: "ASRL::/dev/ttyUSB0::9600::8N2::INSTR",
			want:    ASRL,
		},
		{
			name:    "lowercase tcpip",
			address: "tcpip0::192.168.1.100::inst0::INSTR",
			want:    TCPIP,
		},
		{
			name:    "mixed case Usb",
			address: "Usb0::0x0957::0x0407::MY44012345::INSTR",
			want:    USBTMC,
		},
		{
			name:    "mixed case Asrl",
			address: "Asrl::/dev/ttyUSB0::INSTR",
			want:    ASRL,
		},
		{
			name:    "unknown interface type",
			address: "GPIB0::1::INSTR",
			wantErr: true,
		},
		{
			name:    "empty address",
			address: "",
			wantErr: true,
		},
		{
			name:    "no separator",
			address: "TCPIP0",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := determineInterfaceType(tt.address)
			if tt.wantErr {
				if err == nil {
					t.Errorf("determineInterfaceType(%q) expected error, got %v", tt.address, got)
				}
				return
			}
			if err != nil {
				t.Errorf("determineInterfaceType(%q) unexpected error: %v", tt.address, err)
				return
			}
			if got != tt.want {
				t.Errorf("determineInterfaceType(%q) = %v, want %v", tt.address, got, tt.want)
			}
		})
	}
}
