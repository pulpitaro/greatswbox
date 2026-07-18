# cpu_throttle_watch

This tool continuously monitors **CPU temperatures and frequencies** to catch thermal throttling events in real-time.

## How to compile?

Before compilation, please initialize the `go.mod` file for the project [See `1.` in `Typical problems and fixes`]

To compile you can use standard `golang` compiler (Tested on ver. 1.18.x and above) using command:

```bash
go mod tidy
go build -o "build/ctw" main.go

```

or with `ldflags` for a production-ready, stripped minimal version:

```bash
go mod tidy
go build -ldflags="-s -w" -o "build/ctw" main.go

```

## How to use?

To install the compiled binary system-wide, move it to your local bin directory:

```bash
sudo cp ./build/ctw /usr/local/bin/

```

### Examples of use

1. **Run with default threshold (80.0°C):**

```bash
ctw

```

2. **Set a custom alert threshold to 85.5°C:**

```bash
ctw -t 85.5

```

3. **Run without clearing the screen (useful for log redirection or background history):**

```bash
ctw -cn

```

Some great help (`-h`):

```text
Usage of ./ctw:

Flags:
  -cn
    	Disable clear (clear none)
  -t float
    	Threshold (default 80)

```

## Typical problems and fixes

1. Problem with `golang` mod file

* Fix: Remove `go.mod` and `go.sum` (if exist), then run `go mod init [module name]` and `go mod tidy` to initialize the project context.

2. Problem with execute the compiled binary file

* Fix: Run `sudo chmod +x [binary filename]`

3. Binary immediately exits with "No thermal zones" error

* Fix: The tool relies on standard Linux sysfs paths (`/sys/class/thermal/...`). If you are running this inside a heavily restricted container or a non-standard environment, ensure that the `/sys` filesystem is correctly mounted and accessible.

4. The status line looks broken or overwrites incorrectly

* Fix: By default, the tool uses standard ANSI escape codes (`\033[H\033[2J`) and carriage returns (`\r`) to refresh the terminal view. If your terminal emulator does not support ANSI sequences, pass the `-cn` flag to disable automatic screen clearing.

## Dictionary

ctw = cpu_throttle_watch
cn = clear_none
t = threshold