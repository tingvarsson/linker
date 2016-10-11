package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/tingvarsson/dlog"
)

// TODO: Double printouts of short/long version arguments in helper (as well as double handling in the code)
// Long version arguments
var source = flag.String("source", os.Getenv("PWD"), "path to source of files")
var target = flag.String("target", os.Getenv("HOME"), "path to symlink target")
var user = flag.String("user", os.Getenv("USER"), "name of user for which home dir should be used as symlink target")
var dryrun = flag.Bool("dry-run", false, "Dry-run any operations to the file system")
var force = flag.Bool("force", false, "Force all, default to yes, operations to the file system")
var debug = flag.Bool("debug", false, "Output debugging information to the console")

func init() {
	// Environment variables

	// Short version arguments
	flag.StringVar(source, "s", os.Getenv("PWD"), "path to source of files")
	flag.StringVar(target, "t", os.Getenv("HOME"), "path to symlink target")
	flag.StringVar(user, "u", os.Getenv("USER"), "name of user for which home dir should be used as symlink target")
	flag.BoolVar(dryrun, "n", false, "Dry-run any operations to the file system")
	flag.BoolVar(force, "f", false, "Force all, default to yes, operations to the file system")
	flag.BoolVar(debug, "d", false, "Output debugging information to the console")
}

func main() {
	parseArguments()

	filepath.Walk(*source, handleFile)
}

func parseArguments() {
	flag.Parse()

	if !isDir(*target) {
		log.Fatal("Target is not a path to a directory")
	}

	dlog.SetDebug(*debug)

	logDebugEnvironment()
	logDebugArguments()
}

func logDebugEnvironment() {
	dlog.Debug("ENV $PWD:", os.Getenv("PWD"))
	dlog.Debug("ENV $USER:", os.Getenv("USER"))
	dlog.Debug("ENV $HOME:", os.Getenv("HOME"))
}

func logDebugArguments() {
	dlog.Debug("ARG source:", *source)
	dlog.Debug("ARG target:", *target)
	dlog.Debug("ARG user:", *user)
	dlog.Debug("ARG dryrun:", *dryrun)
	dlog.Debug("ARG debug:", *debug)
}

func handleFile(path string, info os.FileInfo, err error) error {
	dlog.Enter("handleFile()")

	if isDir(path) {
		dlog.Debug("Ignore directory:", path)
		return nil
	}

	dlog.Debug("Source path:", path)

	targetPath, err := extractTargetPath(path, filepath.Dir(*source), *target)
	if err != nil {
		return err
	}
	dlog.Debug("Target path:", targetPath)

	// Check if a target file exists or not
	if _, err := os.Lstat(targetPath); os.IsNotExist(err) {
		return handleNew(path, targetPath)
	} else {
		return handleExisting(path, targetPath)
	}

	return nil
}

func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}

	return fileInfo.IsDir()
}

func extractTargetPath(sourcePath, sourceDir, targetDir string) (string, error) {
	dlog.Enter("extractTargetPath()")

	relativePath, err := filepath.Rel(sourceDir, sourcePath)
	if err != nil {
		log.Panicln(err)
		return "", err
	}
	dlog.Debug("Relative path:", relativePath)

	return filepath.Join(targetDir, relativePath), nil
}

func handleNew(sourcePath, targetPath string) error {
	dlog.Enter("handleNew()")

	return symlink(sourcePath, targetPath)
}

func handleExisting(sourcePath, targetPath string) error {
	dlog.Enter("handleExisting()")

	targetInfo, err := os.Lstat(targetPath)
	if err != nil {
		log.Panicln(err)
		return err
	}

	targetMode := targetInfo.Mode()
	dlog.Debug("targetMode:", targetMode)

	if targetMode&os.ModeSymlink == os.ModeSymlink {
		return handleExistingSymlink(sourcePath, targetPath)
	} else if targetMode.IsRegular() {
		return handleExistingFile(sourcePath, targetPath)
	} else {
		err := errors.New("What the fuck is this? A directory? A bird? A plane? NO IT IS NOT SUPERMAN!")
		log.Panicln(err)
		return err
	}

	return nil
}

func handleExistingSymlink(sourcePath, targetPath string) error {
	dlog.Enter("handleExistingSymlink()")

	evaluatedTargetPath, err := filepath.EvalSymlinks(targetPath)
	if err != nil {
		log.Panicln(err)
		return err
	}

	if evaluatedTargetPath == sourcePath {
		dlog.Debug("Symlink points correctly")
		return nil
	}

	if !promptYesOrNo("[INFO] Existing Symlink points to %s ,replace with new symlink? [yN]", evaluatedTargetPath) {
		dlog.Debug("Symlink points incorrectly but is chosen by the user to not be replaced")
		return nil
	}

	return removeAndSymlink(sourcePath, targetPath)
}

func handleExistingFile(sourcePath, targetPath string) error {
	dlog.Enter("handleExistingFile()")

	equal, err := equalFiles(sourcePath, targetPath)
	if err != nil {
		log.Panicln(err)
		return err
	}

	if !equal && promptYesOrNo("[INFO] Existing file differs, replace with symlink anyway? [yN]") {
		return nil
	}

	return removeAndSymlink(sourcePath, targetPath)
}

// TODO: Fix so that the prompter (fmt.Scanln) accepts an empty string (just newline)
// TODO: Extract prompter to interface/package to enable mocking for test purpose
func promptYesOrNo(format string, args ...interface{}) bool {
	fmt.Printf(format, args)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	// (?i) = case insensitive
	// ^ = match start (exact match)
	// (y|yes) = accept either y or yes
	// $ = match end (exact match)
	re := regexp.MustCompile("(?i)^(y|yes)$")
	return re.MatchString(response)
}

func equalFiles(lhs, rhs string) (bool, error) {
	dlog.Enter("compareFiles()")

	lhsBytes, err := ioutil.ReadFile(lhs)
	if err != nil {
		log.Panicln(err)
		return false, err
	}
	rhsBytes, err := ioutil.ReadFile(rhs)
	if err != nil {
		log.Panicln(err)
		return false, err
	}

	return bytes.Equal(lhsBytes, rhsBytes), nil
}

func removeAndSymlink(sourcePath, targetPath string) error {
	dlog.Enter("removeAndSymlink()")

	if err := os.Remove(targetPath); err != nil {
		log.Fatal(err)
	}

	return symlink(sourcePath, targetPath)
}

func symlink(sourcePath, targetPath string) error {
	dlog.Enter("symlink()")

	prepareDirectory(targetPath)

	if err := os.Symlink(sourcePath, targetPath); err != nil {
		log.Fatal(err)
	}

	return nil
}

func prepareDirectory(targetPath string) error {
	dlog.Enter("prepareDirectory()")

	dirPath := filepath.Dir(targetPath)

	// TODO: What is the correct FileMode to use instead of 0755?
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Fatal(err)
	}

	return nil
}
