# clipboard_to_file

This tool quickly grabs text from **your clipboard** and dumps it **into a file**.

## How to compile?

Before compilation, please initialize the `go.mod` file for the project [See `1.` in `Typical problems and fixes`]

To compile you can use standard `golang` compiler (Tested on ver. 1.18.x and above) using command:

```bash
go build -o "build/c2f" main.go
```

or with `ldflags` for version without addtional dev logging:

```bash
CGO_ENABLED=1 go build -ldflags="-s -w" -o "build/c2f" main.go
```

Note: `CGO_ENABLED=1` is required to compile the native X11/Wayland clipboard bindings.

## How to use?

To install the compiled binary system-wide, move it to your local bin directory:

```bash
sudo cp ./build/c2f /usr/local/bin/
```


Some great help (`-h`):

```text
Usage of c2f:

Quickly grabs text from your clipboard and dumps it into a file.

Usage:
  c2f [filename]
```

## Typical problems and fixes

1. Problem with `golang` mod file
* Fix: Run `go mod init [main or own module name]` and then `go mod tidy` to initialize the project context.

2. Problem with execute the compiled binary file
* Fix: Run `sudo chmod +x [binary filename]`

3. Missing X11 development headers during compilation
* Fix: If the compiler complains about missing X11 libraries, install them via your package manager (e.g., `sudo dnf install libX11-devel` on Fedora or `sudo apt install libx11-dev` on Ubuntu/Debian).

4. Problem with clipboard access on Android (Termux)
* Fix: Ensure you have the `Termux:API` package installed both in Android (app from F-Droid) and inside the terminal. Run: `pkg install termux-api`. If it still fails, check if Termux has permission to access the Android clipboard in your system settings.

## Dictionary

c2f = clipboard_to_file