package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"golang.design/x/clipboard"
)

func main() {
	y := flag.Bool("y", false, "Skip prompt")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("❌ Missing filename"); return
	}
	fn := flag.Arg(0)

	// 1. Overwrite check
	if _, err := os.Stat(fn); err == nil && !*y {
		fmt.Printf("⚠️ Overwrite '%s'? [y/N]: ", fn)
		var res string
		fmt.Scanln(&res)
		if r := strings.ToLower(strings.TrimSpace(res)); r != "y" && r != "yes" {
			fmt.Println("🛑 Aborted"); return
		}
	}

	// 2. Read clipboard (Android vs Desktop)
	var txt string
	if runtime.GOOS == "android" {
		out, err := exec.Command("termux-clipboard-get").Output()
		if err != nil {
			fmt.Println("❌ Termux:API error. Run: pkg install termux-api"); return
		}
		txt = string(out)
	} else {
		if err := clipboard.Init(); err != nil {
			fmt.Println("❌ Clipboard API error. X11/Wayland required"); return
		}
		txt = string(clipboard.Read(clipboard.FmtText))
	}

	if len(strings.TrimSpace(txt)) == 0 {
		fmt.Println("⚠️ Clipboard empty"); return
	}

	// 3. Save to file
	if err := os.WriteFile(fn, []byte(txt), 0644); err != nil {
		fmt.Printf("❌ Write error: %v\n", err); return
	}
	fmt.Printf("📋 Captured into '%s' (%d bytes)!\n", fn, len(txt))
}