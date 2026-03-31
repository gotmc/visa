# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go implementation of the Virtual Instrument Software Architecture (VISA) resource manager. Parses VISA address strings (e.g., `TCPIP0::192.168.1.100::INSTR`) to create resources that abstract hardware interface types (USBTMC, TCPIP, ASRL/serial). Used standalone for SCPI commands or as a foundation for IVI drivers. Based on the VPP-4.3 VISA Library spec.

## Build & Test Commands

```bash
just check          # fmt and vet
just unit           # fmt, vet, then run unit tests with -race -short
just unit -v        # verbose unit tests (extra flags passed through)
just lint           # golangci-lint (config in .golangci.yaml)
just cover          # HTML coverage report (unit tests by default)
just cover all      # HTML coverage for all tests
go test -run TestName ./...  # run a single test
```

Prefer `just` over the `Makefile` (Makefile is outdated).

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
- `Resource` — `io.ReadWriteCloser` + `WriteString`, `Command(ctx, format, a...)`, `Query(ctx, s)`

**Core flow:** `visa.NewResource(address)` -> `determineInterfaceType()` parses the address prefix (USB/TCPIP/ASRL) -> looks up registered driver -> calls `driver.Open()`.

**Drivers** (`driver/` subdirectories) are thin `Connection` wrappers delegating to external libraries:
- `tcpip` -> `github.com/gotmc/lxi`
- `usbtmc` -> `github.com/gotmc/usbtmc`
- `asrl` -> `github.com/gotmc/asrl`

Each driver wraps the upstream device type in a `Connection` struct that satisfies the `visa.Resource` interface, adding `context.Context` to `Command` and `Query`.

## VISA Address Format

Addresses follow the VPP-4.3 spec pattern: `<InterfaceType><BoardIndex>::<details>::INSTR`
- `TCPIP0::192.168.1.100::inst0::INSTR`
- `USB0::0x0957::0x0407::MY44012345::INSTR`
- `ASRL1::/dev/ttyUSB0::INSTR`
