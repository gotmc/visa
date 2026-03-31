# -*- Justfile -*-

coverage_file := "coverage.out"

# List the available justfile recipes.
[group('general')]
@default:
  just --list --unsorted

# List the lines of code in the project.
[group('general')]
loc:
  scc --remap-unknown "-*- Justfile -*-":"justfile"

# Format and vet Go code. Runs before tests.
[group('test')]
check:
	go fmt ./...
	go vet ./...

# Lint using golangci-lint
[group('test')]
lint:
  golangci-lint run --config .golangci.yaml

# Run the unit tests.
[group('test')]
unit *FLAGS: check
  go test ./... -cover -vet=off -race {{FLAGS}} -short

# HTML report for unit (default), int, e2e, or all tests.
[group('test')]
cover test='unit': check
  go test ./... -vet=off -coverprofile={{coverage_file}} \
  {{ if test == 'all' { '' } \
    else if test == 'int' { '-run Integration' } \
    else if test == 'e2e' { '-run E2E' } \
    else { '-short' } }}
  go tool cover -html={{coverage_file}}

# List the outdated direct dependencies (slow to run).
[group('dependencies')]
outdated:
  # (requires https://github.com/psampaz/go-mod-outdated).
  go list -u -m -json all | go-mod-outdated -update -direct

# Update the given module to the latest version.
[group('dependencies')]
update mod:
  go get -u {{mod}}
  go mod tidy

# Update all modules.
[group('dependencies')]
updateall:
  go get -u ./...
  go mod tidy

# Run go mod tidy and verify.
[group('dependencies')]
tidy:
  go mod tidy
  go mod verify

# Build and run the TCPIP Keysight 33220A example application.
[group('examples')]
k33220tcp ip:
  #!/usr/bin/env bash
  echo '# VISA TCPIP Keysight 33220A Example Application'
  cd {{justfile_directory()}}/examples/tcpip/key33220
  env go build -o key33220
  ./key33220 -ip={{ip}}

# Build and run the USBTMC Keysight 33220A example application.
[group('examples')]
k33220usb sn:
  #!/usr/bin/env bash
  echo '# VISA USBTMC Keysight 33220A Example Application'
  cd {{justfile_directory()}}/examples/usbtmc/key33220
  env go build -o key33220
  ./key33220 -sn={{sn}}

# Build and run the ASRL SRS DS345 example application.
[group('examples')]
ds345 port:
  #!/usr/bin/env bash
  echo '# VISA ASRL SRS DS345 Example Application'
  cd {{justfile_directory()}}/examples/asrl/ds345
  env go build -o ds345
  ./ds345 -ser={{port}}
