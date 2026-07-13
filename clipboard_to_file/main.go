package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"golang.design/x/clipboard"
)

// Fetches clipboard content when running inside Termux on Android
func getAndroidClipboard() (string, error) {
	// Check if termux-clipboard-get is available in PATH
	_, err := exec.LookPath("termux-clipboard-get")
	if err != nil {
		return "", fmt.Errorf("Termux:API is missing. Please run: pkg install termux-api")
	}

	cmd := exec.Command("termux-clipboard-get")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to read Android clipboard via Termux API")
	}

	return string(out), nil
}

func main() {
	// Define the -y flag for automatic overwriting
	skipPrompt := flag.Bool("y", false, "Assume 'yes' and skip overwrite confirmation prompt")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Quickly grabs text from your clipboard and dumps it into a file.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  c2f [flags] [filename]\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	// Ensure user provided a target filename
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("❌ Error: Missing target filename.")
		flag.Usage()
		return
	}
	filename := args[0]

	// Foolproof check: Verify if file exists before overwriting
	if _, err := os.Stat(filename); err == nil {
		// If file exists and -y flag was NOT provided, prompt the user
		if !*skipPrompt {
			fmt.Printf("⚠️  Warning: File '%s' already exists. Overwrite? [y/N]: ", filename)
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("❌ Error reading response. Aborting.")
				return
			}
			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" && response != "yes" {
				fmt.Println("🛑 Aborted. File was not overwritten.")
				return
			}
		}
	}

	var content string

	// Smart routing based on current OS runtime environment
	if runtime.GOOS == "android" {
		androidContent, err := getAndroidClipboard()
		if err != nil {
			fmt.Println("❌ Error:", err)
			return
		}
		content = androidContent
	} else {
		// Native desktop X11/Wayland routine
		err := clipboard.Init()
		if err != nil {
			fmt.Printf("❌ Error initializing clipboard API: %v\n", err)
			fmt.Println("👉 Make sure your display server (X11/Wayland) is running active.")
			return
		}
		rawContent := clipboard.Read(clipboard.FmtText)
		content = string(rawContent)
	}

	if len(strings.TrimSpace(content)) == 0 {
		fmt.Println("⚠️  Clipboard content is empty or not valid text. Skipping write.")
		return
	}

	// Write the file atomics to current working directory
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		fmt.Printf("❌ Error writing to file '%s': %v\n", filename, err)
		return
	}

	fmt.Printf("📋 Successfully captured clipboard into '%s' (%d bytes)!\n", filename, len(content))
}