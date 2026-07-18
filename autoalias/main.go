package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	add := flag.String("add", "", "Add 'name=cmd'")
	list := flag.Bool("list", false, "List")
	rm := flag.String("rm", "", "Remove")
	cb := flag.Int("cb", 1, "Keep backups")
	flag.Parse()

	home, _ := os.UserHomeDir()
	rc := home + "/.bashrc"
	if strings.Contains(os.Getenv("SHELL"), "zsh") { rc = home + "/.zshrc" }
	bakDir := home + "/bashrc_bak"
	_ = os.MkdirAll(bakDir, 0755)

	b, err := os.ReadFile(rc)
	if err != nil { fmt.Println("❌ Missing file"); return }
	lines := strings.Split(string(b), "\n")

	if *list {
		for _, l := range lines {
			if strings.HasPrefix(strings.TrimSpace(l), "alias ") { fmt.Println(l) }
		}
		return
	}

	bak := func() {
		_ = os.WriteFile(fmt.Sprintf("%s/%s_%s.bak", bakDir, filepath.Base(rc), time.Now().Format("20060102_150405")), b, 0644)
		if ff, _ := filepath.Glob(bakDir + "/*.bak"); len(ff) > *cb {
			for i := 0; i < len(ff)-*cb; i++ { _ = os.Remove(ff[i]) } // Glob returns sorted by default
		}
	}

	if *add != "" {
		p := strings.SplitN(*add, "=", 2)
		if len(p) != 2 { fmt.Println("❌ Format: name=command"); return }
		n, c := strings.TrimSpace(p[0]), strings.Trim(p[1], "'")
		
		for _, l := range lines {
			if strings.HasPrefix(strings.TrimSpace(l), "alias "+n+"=") { fmt.Printf("⚠️ Exists\n"); return }
		}
		bak()
		f, _ := os.OpenFile(rc, os.O_APPEND|os.O_WRONLY, 0644)
		defer f.Close()
		f.WriteString(fmt.Sprintf("alias %s='%s'\n", n, c))
		fmt.Printf("✅ Added '%s'. Run 'source %s'\n", n, rc)
		return
	}

	if *rm != "" {
		var out []string
		found := false
		for _, l := range lines {
			if strings.HasPrefix(strings.TrimSpace(l), "alias "+*rm+"=") { found = true; continue }
			out = append(out, l)
		}
		if !found { fmt.Printf("⚠️ Not found\n"); return }
		bak()
		os.WriteFile(rc, []byte(strings.Join(out, "\n")), 0644)
		fmt.Printf("🗑️ Removed '%s'. Run 'source %s'\n", *rm, rc)
		return
	}
	flag.Usage()
}