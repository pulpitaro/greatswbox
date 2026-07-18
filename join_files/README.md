# join_files

The tool for continuous text concatenation via F-Strings

## How to compile?

To compile you can use standard `golang` compiler (Tested on ver. 1.18.x and above) using command:

```bash
go mod tidy
go build -o "build/joinf" main.go
```

or with `ldflags` to strip debugging symbols and optimize the binary size:

```bash
go mod tidy
go build -ldflags="-s -w" -o "build/joinf" main.go
```

## How to use?

To install the compiled binary system-wide, move it to your local bin directory:

```bash
sudo cp ./build/joinf /usr/local/bin/
```


Some great help (`-h`):

```help
joinf - Continuous Text Concatenation Utility via F-Strings

Usage:
  joinf [flags] "literal text {file1.txt} more text {file2.txt}"

Flags:
  -h            Show this help manual
  -d            Process entire directory structure if a directory path is found inside {}
  -dev         Enable developer mode with dynamic tracking logs on stderr
  -ext          Extension filter for -d mode (default: ".txt", use "*" or "all" to merge everything)
  -prefix      String tokens to insert BEFORE each processed file block
  -suffix      String tokens to insert AFTER each processed file block

Macros allowed inside --prefix and --suffix (Triggered per-file context):
  \name         Outputs the file name (e.g., source.txt)
  \ext          Outputs the file extension (e.g., .txt)
  \path         Outputs the relative directory path from execution scope
  \full_path    Outputs the absolute computed path from filesystem root
```

## Typical problems and fixes

* Problem with `golang` mod file
* Fix: Remove `go.mod` and `go.sum` (if exist), then run `go mod init [module name]` and `go mod tidy` to initialize the project context.

* Problem with execute the compiled binary file
* Fix: Run `sudo chmod +x [binary filename]`


## Dictionary

join_files = joinf