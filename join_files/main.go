package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Define CLI flags
	helpFlag := flag.Bool("h", false, "Display help and usage information")
	dirFlag := flag.Bool("d", false, "Process the entire directory if a folder path is provided")
	devFlag := flag.Bool("dev", false, "Enable developer mode with detailed console logs on stderr")
	extFlag := flag.String("ext", ".txt", "File extension filter for directory merging (e.g., .md, .json, or * for all files)")
	prefixFlag := flag.String("prefix", "", "Text to prepend BEFORE each file block")
	suffixFlag := flag.String("suffix", "", "Text to append AFTER each file block")

	flag.Parse()

	if *helpFlag {
		printHelp()
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: No input template provided. Use -h for help.")
		os.Exit(1)
	}

	// Join all remaining arguments to rebuild the full template string
	templateStr := strings.Join(args, " ")

	// Smart fallback: If -d is active and user forgot brackets around a path, wrap it automatically
	if *dirFlag && !strings.Contains(templateStr, "{") && !strings.Contains(templateStr, "}") {
		info, err := os.Stat(strings.TrimSpace(templateStr))
		if err == nil && info.IsDir() {
			if *devFlag {
				fmt.Fprintln(os.Stderr, "[DEV] Smart fallback: Directory path detected without brackets. Wrapping automatically.")
			}
			templateStr = "{" + strings.TrimSpace(templateStr) + "}"
		}
	}

	if *devFlag {
		fmt.Fprintf(os.Stderr, "[DEV] Processing template: %s\n", templateStr)
		fmt.Fprintf(os.Stderr, "[DEV] Target extension filter: %s\n", *extFlag)
	}

	parsedPrefix := parseEscapeSequences(*prefixFlag)
	parsedSuffix := parseEscapeSequences(*suffixFlag)

	var finalBuilder strings.Builder
	runes := []rune(templateStr)
	length := len(runes)
	
	var currentBlock strings.Builder
	inBlock := false

	// Single-pass parser to extract literal text and placeholder targets {}
	for i := 0; i < length; i++ {
		r := runes[i]

		if r == '{' && (i == 0 || runes[i-1] != '\\') {
			if inBlock {
				fmt.Fprintln(os.Stderr, "Error: Nested brackets '{' are not allowed.")
				os.Exit(1)
			}
			finalBuilder.WriteString(parseEscapeSequences(currentBlock.String()))
			currentBlock.Reset()
			inBlock = true
		} else if r == '}' && (i == 0 || runes[i-1] != '\\') {
			if !inBlock {
				fmt.Fprintln(os.Stderr, "Error: Unmatched closing bracket '}'.")
				os.Exit(1)
			}
			
			// targetPath processed out of the block context
			targetPath := strings.TrimSpace(currentBlock.String())
			currentBlock.Reset()
			inBlock = false

			if targetPath == "" {
				continue
			}

			// Clean any internal escapes that could reside in target path names
			targetPath = parseEscapeSequences(targetPath)

			fileData, err := processTarget(targetPath, *dirFlag, *devFlag, *extFlag, parsedPrefix, parsedSuffix)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error processing target {%s}: %v\n", targetPath, err)
				os.Exit(1)
			}
			finalBuilder.WriteString(fileData)
		} else {
			// FIX: Safe multi-byte rune backslash removal for literal printing of \{ or \}
			if (r == '{' || r == '}') && i > 0 && runes[i-1] == '\\' {
				blockRunes := []rune(currentBlock.String())
				if len(blockRunes) > 0 {
					currentBlock.Reset()
					currentBlock.WriteString(string(blockRunes[:len(blockRunes)-1]))
				}
			}
			currentBlock.WriteRune(r)
		}
	}

	if inBlock {
		fmt.Fprintln(os.Stderr, "Error: Unclosed bracket '{' at the end of the expression.")
		os.Exit(1)
	}
	if currentBlock.Len() > 0 {
		finalBuilder.WriteString(parseEscapeSequences(currentBlock.String()))
	}

	fmt.Println(finalBuilder.String())
}

// processTarget checks if the target path is a file or folder and processes it accordingly
func processTarget(path string, allowDir, devMode bool, extFilter, prefix, suffix string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if info.IsDir() {
		if !allowDir {
			return "", fmt.Errorf("'%s' is a directory. Enable folder merging with the -d flag", path)
		}

		if devMode {
			fmt.Fprintf(os.Stderr, "[DEV] Merging directory target: %s\n", path)
		}

		// Normalize filter format (ensure it starts with a dot unless it's a wildcard)
		filter := strings.ToLower(extFilter)
		if filter != "*" && filter != "all" && filter != "" && !strings.HasPrefix(filter, ".") {
			filter = "." + filter
		}

		var dirBuilder strings.Builder
		err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}

			// Match extension filter
			fileExt := strings.ToLower(filepath.Ext(d.Name()))
			shouldProcess := filter == "*" || filter == "all" || filter == "" || fileExt == filter

			if shouldProcess {
				content, err := readFileWithMetadata(p, devMode, prefix, suffix)
				if err != nil {
					return err
				}
				dirBuilder.WriteString(content)
			}
			return nil
		})
		
		return dirBuilder.String(), err
	}

	return readFileWithMetadata(path, devMode, prefix, suffix)
}

// readFileWithMetadata loads a file and injects contextual macros into prefixes/suffixes
func readFileWithMetadata(path string, devMode bool, prefix, suffix string) (string, error) {
	if devMode {
		fmt.Fprintf(os.Stderr, "[DEV] Reading file stream: %s\n", path)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}
	relPath := path
	name := filepath.Base(path)
	ext := filepath.Ext(path)

	localPrefix := applyMacros(prefix, name, ext, relPath, absPath)
	localSuffix := applyMacros(suffix, name, ext, relPath, absPath)

	return localPrefix + string(data) + localSuffix, nil
}

// applyMacros replaces custom structural strings with actual metadata values
func applyMacros(target, name, ext, relPath, absPath string) string {
	// FIX: Longest matching macro (\full_path) must be handled BEFORE the shorter substring macro (\path)
	res := strings.ReplaceAll(target, `\full_path`, absPath)
	res = strings.ReplaceAll(res, `\path`, relPath)
	res = strings.ReplaceAll(res, `\name`, name)
	res = strings.ReplaceAll(res, `\ext`, ext)
	return res
}

// parseEscapeSequences manually decodes standard control characters like \n and \t
func parseEscapeSequences(s string) string {
	s = strings.ReplaceAll(s, `\n`, "\n")
	s = strings.ReplaceAll(s, `\t`, "\t")
	s = strings.ReplaceAll(s, `\r`, "\r")
	return s
}

// printHelp displays the dynamic formatting instructions for the user terminal
func printHelp() {
	fmt.Println(`joinf - Continuous Text Concatenation Utility via F-Strings

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
`)
}