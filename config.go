package main

import (
	"flag"
	"os" // TODO: Consider getting rid of dependency to os.Getenv()

	"github.com/tingvarsson/dlog"
)

const (
	FlagUsageSource = "path to source of files"
	FlagUsageTarget = "path to symlink target"
	FlagUsageDryrun = "dry-run any operations to the file system"
	FlagUsageForce  = "force all, default to yes, operations to the file system"
	FlagUsageDebug  = "output debugging information to the console"
	FlagUsageShort  = " (short version)"
)

const (
	EnvPWD  = "PWD"
	EnvHome = "HOME"
)

type Config struct {
	logger *dlog.DebugLogger
	Source string
	Target string
	Dryrun bool
	Force  bool
	Debug  bool
}

func NewConfig(d *dlog.DebugLogger) *Config {
	return &Config{logger: d}
}

func (c *Config) Init() {
	// TODO: Double printouts of short/long version arguments in helper (as well as double handling in the code)
	flag.StringVar(&c.Source, "source", os.Getenv(EnvPWD), FlagUsageSource)
	flag.StringVar(&c.Target, "target", os.Getenv(EnvHome), FlagUsageTarget)
	flag.BoolVar(&c.Dryrun, "dry-run", false, FlagUsageDryrun)
	flag.BoolVar(&c.Force, "force", false, FlagUsageForce)
	flag.BoolVar(&c.Debug, "debug", false, FlagUsageDebug)
	flag.StringVar(&c.Source, "s", os.Getenv(EnvPWD), FlagUsageSource+FlagUsageShort)
	flag.StringVar(&c.Target, "t", os.Getenv(EnvHome), FlagUsageTarget+FlagUsageShort)
	flag.BoolVar(&c.Dryrun, "n", false, FlagUsageDryrun+FlagUsageShort)
	flag.BoolVar(&c.Force, "f", false, FlagUsageForce+FlagUsageShort)
	flag.BoolVar(&c.Debug, "d", false, FlagUsageDebug+FlagUsageShort)
}

func (c *Config) ParseFlags() {
	flag.Parse()

	if c.Debug {
		c.logger.EnableDebug()
	}

	c.logDebugEnvironment()
	c.logDebugArguments()
}

func (c *Config) logDebugEnvironment() {
	c.logger.Debug("ENV $PWD: ", os.Getenv(EnvPWD))
	c.logger.Debug("ENV $HOME: ", os.Getenv(EnvHome))
}

func (c *Config) logDebugArguments() {
	c.logger.Debug("ARG source: ", c.Source)
	c.logger.Debug("ARG target: ", c.Target)
	c.logger.Debug("ARG dryrun: ", c.Dryrun)
	c.logger.Debug("ARG force: ", c.Force)
	c.logger.Debug("ARG debug: ", c.Debug)
}

func (c *Config) Verify() {
	// TODO: Add sanity check of source to be a path

	if !isDir(c.Target) {
		c.logger.Fatal("Target is not a path to a directory")
	}
}
