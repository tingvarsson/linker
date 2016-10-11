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

	logDebugEnvironment()
	logDebugArguments()
}

func logDebugEnvironment() {
	logDebug("ENV $PWD:", os.Getenv("PWD"))
	logDebug("ENV $USER:", os.Getenv("USER"))
	logDebug("ENV $HOME:", os.Getenv("HOME"))
}

func logDebugArguments() {
	logDebug("ARG source:", *source)
	logDebug("ARG target:", *target)
	logDebug("ARG user:", *user)
	logDebug("ARG dryrun:", *dryrun)
	logDebug("ARG debug:", *debug)
}

// TODO: Integrate debug control into the logger instead of having to have ugly if statements directly in the code
// TODO: Extend the logger even further to also have ready generic functionality to log function ENTRY/EXIT
func logDebug(args ...interface{}) {
	if *debug {
		log.Println("[DEBUG]", args)
	}
}

func handleFile(path string, info os.FileInfo, err error) error {
	logDebug("ENTER handleFile()")

	if isDir(path) {
		logDebug("Ignore directory:", path)
		return nil
	}

	logDebug("Source path:", path)

	targetPath, err := extractTargetPath(path, filepath.Dir(*source), *target)
	if err != nil {
		return err
	}
	logDebug("Target path:", targetPath)

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
	logDebug("ENTER extractTargetPath()")

	relativePath, err := filepath.Rel(sourceDir, sourcePath)
	if err != nil {
		log.Panicln(err)
		return "", err
	}
	logDebug("Relative path:", relativePath)

	return filepath.Join(targetDir, relativePath), nil
}

func handleNew(sourcePath, targetPath string) error {
	logDebug("ENTER handleNew()")

	return symlink(sourcePath, targetPath)
}

func handleExisting(sourcePath, targetPath string) error {
	logDebug("ENTER handleExisting()")

	targetInfo, err := os.Lstat(targetPath)
	if err != nil {
		log.Panicln(err)
		return err
	}

	targetMode := targetInfo.Mode()
	logDebug("targetMode:", targetMode)

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
	logDebug("ENTER handleExistingSymlink()")

	evaluatedTargetPath, err := filepath.EvalSymlinks(targetPath)
	if err != nil {
		log.Panicln(err)
		return err
	}

	if evaluatedTargetPath == sourcePath {
		logDebug("Symlink points correctly")
		return nil
	}

	if !promptYesOrNo(fmt.Sprintf("[INFO] Existing Symlink points to %s ,replace with new symlink? [yN]", evaluatedTargetPath)) {
		logDebug("Symlink points incorrectly but is chosen by the user to not be replaced")
		return nil
	}

	return removeAndSymlink(sourcePath, targetPath)
}

func handleExistingFile(sourcePath, targetPath string) error {
	logDebug("ENTER handleExistingFile()")

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

// TODO: Fix so that the prompter accepts an empty string (just newline)
func promptYesOrNo(output string) bool {
	fmt.Print(output)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	re := regexp.MustCompile("(?i)^(y|yes)$")
	return re.MatchString(response)
}

func equalFiles(lhs, rhs string) (bool, error) {
	logDebug("ENTER compareFiles()")

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
	logDebug("ENTER removeAndSymlink()")

	if err := os.Remove(targetPath); err != nil {
		log.Fatal(err)
	}

	return symlink(sourcePath, targetPath)
}

func symlink(sourcePath, targetPath string) error {
	logDebug("ENTER symlink()")

	prepareDirectory(targetPath)

	if err := os.Symlink(sourcePath, targetPath); err != nil {
		log.Fatal(err)
	}

	return nil
}

func prepareDirectory(targetPath string) error {
	logDebug("ENTER prepareDirectory()")

	dirPath := filepath.Dir(targetPath)

	// TODO: What is the correct FileMode to use instead of 0755?
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		log.Fatal(err)
	}

	return nil
}
