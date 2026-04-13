// Copyright (c) 2017-2026 The visa developers. All rights reserved.
// Project site: https://github.com/gotmc/visa
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package visa

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"testing"
)

// mockResource is a minimal Resource implementation for testing.
type mockResource struct{}

func (m *mockResource) Close() error                                          { return nil }
func (m *mockResource) Read(p []byte) (int, error)                            { return 0, nil }
func (m *mockResource) Write(p []byte) (int, error)                           { return 0, nil }
func (m *mockResource) WriteString(s string) (int, error)                     { return 0, nil }
func (m *mockResource) ReadBinary(_ context.Context, p []byte) (int, error)  { return 0, nil }
func (m *mockResource) WriteBinary(_ context.Context, p []byte) (int, error) { return 0, nil }
func (m *mockResource) Command(_ context.Context, _ string, _ ...any) error   { return nil }
func (m *mockResource) Query(_ context.Context, _ string) (string, error)     { return "", nil }

// mockDriver implements Driver for testing.
type mockDriver struct {
	openErr error
}

func (d *mockDriver) Open(_ context.Context, _ string) (Resource, error) {
	if d.openErr != nil {
		return nil, d.openErr
	}
	return &mockResource{}, nil
}

// saveAndClearDrivers saves the current driver registry and returns a restore
// function. Use with t.Cleanup to isolate tests from each other.
func saveAndClearDrivers(t *testing.T) {
	t.Helper()
	saved := make(map[InterfaceType]Driver)
	maps.Copy(saved, drivers)
	drivers = make(map[InterfaceType]Driver)
	t.Cleanup(func() {
		drivers = saved
	})
}

func TestRegister(t *testing.T) {
	t.Run("register new driver", func(t *testing.T) {
		saveAndClearDrivers(t)
		Register(TCPIP, &mockDriver{})
		driversMu.RLock()
		_, exists := drivers[TCPIP]
		driversMu.RUnlock()
		if !exists {
			t.Error("expected TCPIP driver to be registered")
		}
	})

	t.Run("register multiple drivers", func(t *testing.T) {
		saveAndClearDrivers(t)
		Register(TCPIP, &mockDriver{})
		Register(USBTMC, &mockDriver{})
		Register(ASRL, &mockDriver{})
		driversMu.RLock()
		count := len(drivers)
		driversMu.RUnlock()
		if count != 3 {
			t.Errorf("expected 3 registered drivers, got %d", count)
		}
	})

	t.Run("panics on duplicate registration", func(t *testing.T) {
		saveAndClearDrivers(t)
		Register(TCPIP, &mockDriver{})
		defer func() {
			r := recover()
			if r == nil {
				t.Error("expected panic on duplicate registration")
			}
			want := "visa: TCP-IP driver already registered"
			if got, ok := r.(string); !ok || got != want {
				t.Errorf("panic message = %q, want %q", got, want)
			}
		}()
		Register(TCPIP, &mockDriver{})
	})
}

func TestNewResource(t *testing.T) {
	ctx := context.Background()

	t.Run("successful open", func(t *testing.T) {
		saveAndClearDrivers(t)
		Register(TCPIP, &mockDriver{})
		res, err := NewResource(ctx, "TCPIP0::192.168.1.100::inst0::INSTR")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res == nil {
			t.Fatal("expected non-nil resource")
		}
	})

	t.Run("invalid address", func(t *testing.T) {
		saveAndClearDrivers(t)
		_, err := NewResource(ctx, "bogus")
		if err == nil {
			t.Fatal("expected error for invalid address")
		}
		if !errors.Is(err, ErrInvalidAddress) {
			t.Errorf("expected ErrInvalidAddress, got: %v", err)
		}
	})

	t.Run("unregistered driver", func(t *testing.T) {
		saveAndClearDrivers(t)
		_, err := NewResource(ctx, "TCPIP0::192.168.1.100::inst0::INSTR")
		if err == nil {
			t.Fatal("expected error for unregistered driver")
		}
		if !errors.Is(err, ErrDriverNotRegistered) {
			t.Errorf("expected ErrDriverNotRegistered, got: %v", err)
		}
	})

	t.Run("driver open error propagates", func(t *testing.T) {
		saveAndClearDrivers(t)
		openErr := fmt.Errorf("connection refused")
		Register(TCPIP, &mockDriver{openErr: openErr})
		_, err := NewResource(ctx, "TCPIP0::192.168.1.100::inst0::INSTR")
		if err == nil {
			t.Fatal("expected error from driver.Open")
		}
		if err.Error() != "connection refused" {
			t.Errorf("expected 'connection refused', got: %v", err)
		}
	})

	t.Run("USB address routes to USBTMC driver", func(t *testing.T) {
		saveAndClearDrivers(t)
		Register(USBTMC, &mockDriver{})
		res, err := NewResource(ctx, "USB0::0x0957::0x0407::MY44012345::INSTR")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res == nil {
			t.Fatal("expected non-nil resource")
		}
	})

	t.Run("ASRL address routes to ASRL driver", func(t *testing.T) {
		saveAndClearDrivers(t)
		Register(ASRL, &mockDriver{})
		res, err := NewResource(ctx, "ASRL1::/dev/ttyUSB0::INSTR")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res == nil {
			t.Fatal("expected non-nil resource")
		}
	})
}

func TestInterfaceTypeString(t *testing.T) {
	tests := []struct {
		ifaceType InterfaceType
		want      string
	}{
		{UNKNOWN, "Unknown"},
		{USBTMC, "USBTMC"},
		{TCPIP, "TCP-IP"},
		{ASRL, "Serial"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.ifaceType.String(); got != tt.want {
				t.Errorf("InterfaceType(%d).String() = %q, want %q", tt.ifaceType, got, tt.want)
			}
		})
	}
}
