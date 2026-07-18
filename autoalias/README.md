# autoalias

This tool provides a lightweight, automated way to **manage, view, and safely backup your shell aliases** directly into your configuration files (`.bashrc` or `.zshrc`).

## How to compile?

Before compilation, please initialize the `go.mod` file for the project [See `1.` in `Typical problems and fixes`]

To compile you can use standard `golang` compiler (Tested on ver. 1.18.x and above) using command:

```bash
go mod tidy
go build -o "build/autoalias" main.go

```

or with `ldflags` for a production-ready, stripped minimal version:

```bash
go mod tidy
go build -ldflags="-s -w" -o "build/autoalias" main.go

```

## How to use?

To install the compiled binary system-wide, move it to your local bin directory:

```bash
sudo cp ./build/autoalias /usr/local/bin/

```

### Examples of use

1. **List all managed aliases:**

```bash
autoalias -list

```

2. **Add a new alias:**

```bash
autoalias -add "ll=ls -laG"

```

3. **Remove an existing alias:**

```bash
autoalias -rm "ll"

```

4. **Add an alias and limit total backup files to max 3:**

```bash
autoalias -add "gs=git status" -cb 3

```

Some great help (`-h`):

```text
Usage of ./autoalias:

Flags:
  -add string
    	Add 'name=cmd'
  -cb int
    	Keep backups (default 1)
  -list
    	List
  -rm string
    	Remove

```

## Typical problems and fixes

1. Problem with `golang` mod file

* Fix: Remove `go.mod` and `go.sum` (if exist), then run `go mod init [module name]` and `go mod tidy` to initialize the project context.

2. Problem with execute the compiled binary file

* Fix: Run `sudo chmod +x [binary filename]`

3. Changes are not visible immediately after running the command

* Fix: The script automatically appends or modifies the files, but you need to reload your current shell instance to apply them. Run `source ~/.bashrc` or `source ~/.zshrc` depending on your environment.

4. Backups are accumulating too fast or taking too much space

* Fix: Use the `-cb [number]` flag when performing adjustments. By default, the tool sets a conservative value of `1` and cleans up old historical snapshot iterations automatically. All backups are safely stored in the `~/bashrc_bak/` directory.

## Dictionary

cb = clean_backups
rm = remove
rc = run_commands configuration file (.bashrc / .zshrc)