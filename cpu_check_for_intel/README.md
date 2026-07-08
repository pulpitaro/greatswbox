# intel_cpu_check

## How to compile?

To compile you can use standard `golang` compiler (Tested on ver. 1.18.x and 1.22.5) using command:

```bash
go build -o "build/intel_cpu_check" main.go
```

or with `ldflags` for version without addtional dev logging:

```bash
go build -ldflags="-s -w" -o "build/intel_cpu_check" main.go
```

## How to use?

You can `Go` interpreter or run the compiled binary file! You can run this command if you want make the compiled binary the OS binary:

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
* Fix: Run the command `go mod init [main or own module name]` and run `go mod tidy`

* Problem with execute the compiled binary file
* Fix: Run the command `sudo chmod +x [binary filename]`