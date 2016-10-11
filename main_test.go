package main

import (
	"testing"
)

// General test ideas
// either mock fs or create temp files!
// temp files = test directory with source & target

// Verif: ENV & ARGS (a whole bunch of variations)
// no args => S:$PWD T:$HOME
// -u => S:$PWD T:/home/$USER
// -t => S:$PWD T:target
// -u -t => S:$PWD T:target (target over user)
// -s => S:source T:$HOME
// -u nonUser => error
// -s nonPath => error
// -t nonDir => error
// -t nonPath => error
// no args, $PWD not set => error
// no args, $HOME not set => error
// -u, $USER not set => error
func testInitArguments(t *testing.T) {
}

// Verif: S:file w. any scenario
// Verif: directory => ignore
// Verif: new file
// Verif: new file without directories
// Verif: existing symlink (points correctly)
// Verif: existing symlink (points to other)
// Verif: existing file (same content)
// Verif: existing file (different content)
func testHandleFile(t *testing.T) {
}

// Verif: dry-run

// Veriy: force

// Verif: debug mode (verify printouts?)

// Bench: Setup a major complex scenario including all variations
func benchHandleFile(t *testing.T) {
}
