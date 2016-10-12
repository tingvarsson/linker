package main

import (
	"os"
	"testing"
)

const (
	Command     = "linker"
	DefaultPWD  = "/PWD/"
	DefaultHome = "/HOME/"
	ShortTarget = "/T/"
	LongTarget  = "/TARGET/"
	ShortSource = "/S/"
	LongSource  = "/SOURCE/"
)

func (c *Config) Reset() {
	c.source = DefaultPWD
	c.target = DefaultHome
	c.dryrun = false
	c.force = false
	c.debug = false
}

// TODO: Prints a bunch when enabling debug mode, should be handled with a testLogger (that later can be used for log verification)
func TestParseArguments(t *testing.T) {
	var cases = []struct {
		args           []string
		expectedSource string
		expectedTarget string
		expectedDryrun bool
		expectedForce  bool
		expectedDebug  bool
	}{
		{[]string{Command}, DefaultPWD, DefaultHome, false, false, false},
		{[]string{Command, "-t", ShortTarget}, DefaultPWD, ShortTarget, false, false, false},
		{[]string{Command, "-target", LongTarget}, DefaultPWD, LongTarget, false, false, false},
		{[]string{Command, "-s", ShortSource}, ShortSource, DefaultHome, false, false, false},
		{[]string{Command, "-source", LongSource}, LongSource, DefaultHome, false, false, false},
		{[]string{Command, "-n"}, DefaultPWD, DefaultHome, true, false, false},
		{[]string{Command, "-dry-run"}, DefaultPWD, DefaultHome, true, false, false},
		{[]string{Command, "-f"}, DefaultPWD, DefaultHome, false, true, false},
		{[]string{Command, "-force"}, DefaultPWD, DefaultHome, false, true, false},
		{[]string{Command, "-d"}, DefaultPWD, DefaultHome, false, false, true},
		{[]string{Command, "-debug"}, DefaultPWD, DefaultHome, false, false, true},
		{[]string{Command, "-s", LongSource, "-t", LongTarget, "-n", "-f", "-d"}, LongSource, LongTarget, true, true, true},
	}

	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	for i, c := range cases {
		// Reset & Setup
		config.Reset()
		os.Args = c.args

		// Execute
		config.ParseFlags()

		// Verify
		if config.source != c.expectedSource {
			t.Errorf("[CASE:%d] Source is %s expected %s", i, config.source, c.expectedSource)
		}
		if config.target != c.expectedTarget {
			t.Errorf("[CASE:%d] Target is %s expected %s", i, config.target, c.expectedTarget)
		}
		if config.dryrun != c.expectedDryrun {
			t.Errorf("[CASE:%d] Dryrun is %t expected %t", i, config.dryrun, c.expectedDryrun)
		}
		if config.force != c.expectedForce {
			t.Errorf("[CASE:%d] Force is %t expected %t", i, config.force, c.expectedForce)
		}
		if config.debug != c.expectedDebug {
			t.Errorf("[CASE:%d] Debug is %t expected %t", i, config.debug, c.expectedDebug)
		}
	}
}
