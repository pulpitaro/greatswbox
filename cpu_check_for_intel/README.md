# intel_cpu_check

This tool detects **Intel Hybrid Architecture (P/E Cores)** and tracks temperatures via sysfs.

## How to compile?

To compile you can use standard `golang` compiler (Tested on ver. 1.18.x and above) using command:

```bash
go mod tidy
go build -o "build/intel_cpu_check" main.go
```

or with `ldflags` to strip debugging symbols and optimize the binary size:

```bash
go mod tidy
go build -ldflags="-s -w" -o "build/intel_cpu_check" main.go
```

## How to use?

To install the compiled binary system-wide, move it to your local bin directory:

```bash
sudo cp ./build/intel_cpu_check /usr/local/bin/
```


Some great help (`-h`):

```help
Usage of intel_cpu_check:

This tool detects Intel Hybrid Architecture (P/E Cores) and tracks temperatures.

Options:
  -c		Run once and exit immediately (snapshot mode)
  -h, --help	Show this help message

Examples:
  sudo intel_cpu_check         # Runs as a live monitor (10s refresh)
  sudo intel_cpu_check -c      # Prints the CPU map once and exits
```

## Typical problems and fixes

* Problem with `golang` mod file
* Fix: Remove `go.mod` and `go.sum` (if exist), then run `go mod init [module name]` and `go mod tidy` to initialize the project context.

* Problem with execute the compiled binary file
* Fix: Run `sudo chmod +x [binary filename]`


## Dictionary

intel_cpu_check = cpu_check_for_intel