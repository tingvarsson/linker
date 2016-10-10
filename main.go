package main

import (
	"flag"
)

// Long version arguments
var source = flag.String("source", "./", "path to source of files")
var target = flag.String("target", "$HOME", "path to symlink target")
var user = flag.String("user", "$USER", "name of user for which home dir should be used as symlink target")
var dryrun = flag.Bool("dry-run", false, "Dry-run any operations to the file system")

func init() {
	// Short version arguments
	flag.StringVar(source, "s", "./", "path to source of files")
	flag.StringVar(target, "t", "$HOME", "path to symlink target")
	flag.StringVar(user, "u", "$USER", "name of user for which home dir should be used as symlink target")
	flag.BoolVar(dryrun, "n", false, "Dry-run any operations to the file system")
}

func main() {
	flag.Parse()
}
