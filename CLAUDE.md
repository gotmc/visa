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

Prefer `just` (no Makefile in this repo). Requires Go 1.25+.

## Architecture

**Driver registration pattern** (like `database/sql`): drivers self-register via `init()` calling `visa.Register()`. Users blank-import the drivers they need:

```go
import (
    "github.com/gotmc/visa"
    _ "github.com/gotmc/visa/driver/tcpip"
)
```

**Key interfaces:**
- `Driver` — has `Open(ctx, address) (Resource, error)`, implemented by each driver package
- `Resource` — `io.ReadWriteCloser` + `WriteString`, `ReadContext`, `WriteContext`, `Command(ctx, format, a...)`, `Query(ctx, s)`

**Core flow:** `visa.NewResource(ctx, address)` -> `determineInterfaceType()` parses the address prefix (USB/TCPIP/ASRL) -> looks up registered driver -> calls `driver.Open(ctx, address)`.

**Drivers** (`driver/` subdirectories) are thin `Connection` wrappers delegating to external libraries:
- `tcpip` -> `github.com/gotmc/lxi`
- `usbtmc` -> `github.com/gotmc/usbtmc`
- `asrl` -> `github.com/gotmc/asrl`

Each driver wraps the upstream device type in a `Connection` struct that satisfies the `visa.Resource` interface. Most drivers delegate context directly to upstream libraries; the `usbtmc` driver still uses a select/channel goroutine pattern to race `Open()` against `ctx.Done()` because the upstream `usbtmc.NewContext`/`NewDevice` calls don't accept a context.

**Error handling:** Sentinel errors in `errors.go` (`ErrInvalidAddress`, `ErrUnknownInterfaceType`, `ErrDriverNotRegistered`) are wrapped with `%w` for programmatic checking via `errors.Is`. `Register()` panics on duplicate driver registration (fail-fast by design).

**Testing conventions:** Table-driven tests with `t.Run()` subtests. Tests are co-located with source files (e.g., `helpers_test.go`).

**`InterfaceType` enum:** `UNKNOWN` starts at `-1` (iota - 1), so `USBTMC=0`, `TCPIP=1`, `ASRL=2`. The VISA address regex is compiled once at package level in `helpers.go`.

## VISA Address Format

Addresses follow the VPP-4.3 spec pattern: `<InterfaceType><BoardIndex>::<details>::INSTR`
- `TCPIP0::192.168.1.100::inst0::INSTR`
- `USB0::0x0957::0x0407::MY44012345::INSTR`
- `ASRL1::/dev/ttyUSB0::INSTR`
