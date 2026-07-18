package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func runMainWithArgs(t *testing.T, tmpHome string, args []string) {
	t.Helper()
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	
	flag.CommandLine = flag.NewFlagSet(args[0], flag.ExitOnError)
	os.Args = args

	t.Setenv("HOME", tmpHome)
	t.Setenv("SHELL", "/bin/bash")

	main()
}

func TestAliasManagerMinimal(t *testing.T) {
	tmpHome := t.TempDir()
	rcPath := filepath.Join(tmpHome, ".bashrc")

	initialContent := []string{
		"my_func() {",
		"    echo \"hello (world)\"",
		"}",
		"alias ll='ls -la'",
		"alias old='echo test'",
	}
	_ = os.WriteFile(rcPath, []byte(strings.Join(initialContent, "\n")), 0644)

	// 1. Test adding new alias
	t.Run("Add new alias", func(t *testing.T) {
		runMainWithArgs(t, tmpHome, []string{"cmd", "-add", "g=git status"})

		bytes, _ := os.ReadFile(rcPath)
		content := string(bytes)

		if !strings.Contains(content, "alias g='git status'") {
			t.Error("❌ Alias was not added")
		}
		if !strings.Contains(content, "my_func() {") {
			t.Error("❌ Brackets or structural code corrupted")
		}
	})

	// 2. Test duplicate protection
	t.Run("Skip existing duplicate alias", func(t *testing.T) {
		runMainWithArgs(t, tmpHome, []string{"cmd", "-add", "ll=ls -lah"})

		bytes, _ := os.ReadFile(rcPath)
		content := string(bytes)

		if strings.Contains(content, "alias ll='ls -lah'") {
			t.Error("❌ Duplicate alias was incorrectly added")
		}
	})

	// 3. Test removal and text safety
	t.Run("Remove alias safely", func(t *testing.T) {
		runMainWithArgs(t, tmpHome, []string{"cmd", "-rm", "old"})

		bytes, _ := os.ReadFile(rcPath)
		content := string(bytes)

		if strings.Contains(content, "alias old=") {
			t.Error("❌ Alias 'old' still exists")
		}
		if !strings.Contains(content, "alias ll='ls -la'") {
			t.Error("❌ Unrelated alias lost")
		}
		if !strings.Contains(content, "my_func() {") {
			t.Error("❌ Brackets broken during removal")
		}
	})

	// 4. Test backup directory rotation
	t.Run("Verify backup rotation", func(t *testing.T) {
		bakDir := filepath.Join(tmpHome, "bashrc_bak")
		
		runMainWithArgs(t, tmpHome, []string{"cmd", "-add", "b1=echo 1", "-cb", "2"})
		runMainWithArgs(t, tmpHome, []string{"cmd", "-add", "b2=echo 2", "-cb", "2"})
		runMainWithArgs(t, tmpHome, []string{"cmd", "-add", "b3=echo 3", "-cb", "2"})

		files, _ := filepath.Glob(bakDir + "/*.bak")
		if len(files) > 2 {
			t.Errorf("❌ Expected max 2 backups, found %d", len(files))
		}
	})
}