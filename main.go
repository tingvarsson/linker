package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
)

// Long version arguments
var source = flag.String("source", "./", "path to source of files")
var target = flag.String("target", "$HOME", "path to symlink target")
var user = flag.String("user", "$USER", "name of user for which home dir should be used as symlink target")
var dryrun = flag.Bool("dry-run", false, "Dry-run any operations to the file system")
var debug = flag.Bool("debug", false, "Output debugging information to the console")

func init() {
	// Short version arguments
	flag.StringVar(source, "s", "./", "path to source of files")
	flag.StringVar(target, "t", "$HOME", "path to symlink target")
	flag.StringVar(user, "u", "$USER", "name of user for which home dir should be used as symlink target")
	flag.BoolVar(dryrun, "n", false, "Dry-run any operations to the file system")
	flag.BoolVar(debug, "d", false, "Output debugging information to the console")
}

func main() {
	flag.Parse()

	filepath.Walk(*source, handleFile)
}

func handleFile(path string, info os.FileInfo, err error) error {
	if *debug {
		log.Println("[DEBUG]", path)
	}

	return nil
}
