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
	c.Source = DefaultPWD
	c.Target = DefaultHome
	c.Dryrun = false
	c.Force = false
	c.Debug = false
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
		if config.Source != c.expectedSource {
			t.Errorf("[CASE:%d] Source is %s expected %s", i, config.Source, c.expectedSource)
		}
		if config.Target != c.expectedTarget {
			t.Errorf("[CASE:%d] Target is %s expected %s", i, config.Target, c.expectedTarget)
		}
		if config.Dryrun != c.expectedDryrun {
			t.Errorf("[CASE:%d] Dryrun is %t expected %t", i, config.Dryrun, c.expectedDryrun)
		}
		if config.Force != c.expectedForce {
			t.Errorf("[CASE:%d] Force is %t expected %t", i, config.Force, c.expectedForce)
		}
		if config.Debug != c.expectedDebug {
			t.Errorf("[CASE:%d] Debug is %t expected %t", i, config.Debug, c.expectedDebug)
		}
	}
}
