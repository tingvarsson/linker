package main

import (
	"testing"
)

// Verif: S:file w. any scenario
// Verif: directory => ignore
// Verif: new file
// Verif: new file without directories
// Verif: existing symlink (points correctly)
// Verif: existing symlink (points to other) (both yes/no scenario)
// Verif: existing file (same content)
// Verif: existing file (different content) (both yes/no scenario)
func TestHandleFile(t *testing.T) {
}

// Verif: dry-run

// Veriy: force

// Verif: debug mode (verify printouts?)
// oldLogger := logger; defer func() {logger = oldLogger}()

// Bench: Setup a major complex scenario including all variations
func BenchmarkHandleFile(b *testing.B) {
}
