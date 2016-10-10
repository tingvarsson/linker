package main

import (
	"testing"
)

// General test ideas
// either mock fs or create temp files!
// temp files = test directory with source & target

// Verif: ENV & ARGS (a whole bunch of variations)
func testInitArguments(t *testing.T) {
}

// Verif: debug mode

// Verif: new file
// Verif: new file without directories
// Verif: existing symlink
// Verif: existing file
func testHandleFile(t *testing.T) {
}

// Verif: dry-run

// Bench: Setup a major complex scenario including all variations
func benchHandleFile(t *testing.T) {
}
