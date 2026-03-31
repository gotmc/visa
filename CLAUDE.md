# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go implementation of the Virtual Instrument Software Architecture (VISA) resource manager. Parses VISA address strings (e.g., `TCPIP0::192.168.1.100::INSTR`) to create resources that abstract hardware interface types (USBTMC, TCPIP, ASRL/serial). Used standalone for SCPI commands or as a foundation for IVI drivers.

## Build & Test Commands

```bash
just check       # fmt, vet, test with coverage
just checkv      # same but verbose
just lint        # staticcheck
just cover       # HTML coverage report
go test ./...    # run all tests
go test -run TestName ./...  # run a single test
```

Both `Makefile` and `Justfile` exist; prefer `just`.

## Architecture

**Driver registration pattern** (like `database/sql`): drivers self-register via `init()` calling `visa.Register()`. Users blank-import the drivers they need:

```go
import (
    "github.com/gotmc/visa"
    _ "github.com/gotmc/visa/driver/tcpip"
)
```

**Key interfaces:**
- `Driver` — has `Open(address string) (Resource, error)`, implemented by each driver package
- `Resource` — `io.ReadWriteCloser` + `WriteString`, `Command`, `Query`

**Core flow:** `visa.NewResource(address)` → `determineInterfaceType()` parses the address prefix (USB/TCPIP/ASRL) → looks up registered driver → calls `driver.Open()`.

**Drivers** (`driver/` subdirectories) are thin wrappers delegating to external libraries:
- `tcpip` → `github.com/gotmc/lxi`
- `usbtmc` → `github.com/gotmc/usbtmc`
- `asrl` → `github.com/gotmc/asrl`

## VISA Address Format

Addresses follow the VPP-4.3 spec pattern: `<InterfaceType><BoardIndex>::<details>::INSTR`
- `TCPIP0::192.168.1.100::inst0::INSTR`
- `USB0::0x0957::0x0407::MY44012345::INSTR`
- `ASRL1::/dev/ttyUSB0::INSTR`
