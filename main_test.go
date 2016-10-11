package main

import (
	"os"
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
func TestParseArguments(t *testing.T) {
	var cases = []struct {
		args           []string
		envPwd         string
		envHome        string
		envUser        string
		expectedSource string
		expectedTarget string
	}{
		{[]string{"linker"}, "/PWD/", "/HOME/", "username", "/PWD/", "/HOME/"},
	}

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	for _, c := range cases {
		os.Args = c.args
		if err := os.Setenv("PWD", c.envPwd); err != nil {
			t.Fatal(err)
		}
		if err := os.Setenv("HOME", c.envHome); err != nil {
			t.Fatal(err)
		}
		if err := os.Setenv("USER", c.envUser); err != nil {
			t.Fatal(err)
		}
		parseArguments()
		// TODO verify source, target & error
		if *source != c.expectedSource {
			t.Errorf("got %s expected %s", *source, c.expectedSource)
		}
		if *target != c.expectedTarget {
			t.Errorf("got %s expected %s", *target, c.expectedTarget)
		}
	}
}

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

// Bench: Setup a major complex scenario including all variations
func BenchmarkHandleFile(b *testing.B) {
}
