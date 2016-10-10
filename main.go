package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// TODO: Double printouts of short/long version arguments in helper (as well as double handling in the code)
// Long version arguments
var source = flag.String("source", os.Getenv("PWD"), "path to source of files")
var target = flag.String("target", os.Getenv("HOME"), "path to symlink target")
var user = flag.String("user", os.Getenv("USER"), "name of user for which home dir should be used as symlink target")
var dryrun = flag.Bool("dry-run", false, "Dry-run any operations to the file system")
var debug = flag.Bool("debug", false, "Output debugging information to the console")

func init() {
	// Environment variables

	// Short version arguments
	flag.StringVar(source, "s", os.Getenv("PWD"), "path to source of files")
	flag.StringVar(target, "t", os.Getenv("HOME"), "path to symlink target")
	flag.StringVar(user, "u", os.Getenv("USER"), "name of user for which home dir should be used as symlink target")
	flag.BoolVar(dryrun, "n", false, "Dry-run any operations to the file system")
	flag.BoolVar(debug, "d", false, "Output debugging information to the console")
}

func main() {
	flag.Parse()
	envDebug()
	flagDebug()

	filepath.Walk(*source, handleFile)
}

func envDebug() {
	logDebug("ENV $PWD:", os.Getenv("PWD"))
	logDebug("ENV $USER:", os.Getenv("USER"))
	logDebug("ENV $HOME:", os.Getenv("HOME"))
}

func flagDebug() {
	logDebug("ARG source:", *source)
	logDebug("ARG target:", *target)
	logDebug("ARG user:", *user)
	logDebug("ARG dryrun:", *dryrun)
	logDebug("ARG debug:", *debug)
}

// TODO: Integrate debug control into the logger instead of having to have ugly if statements directly in the code
func logDebug(args ...interface{}) {
	if *debug {
		log.Println("[DEBUG]", args)
	}
}

func handleFile(path string, info os.FileInfo, err error) error {
	logDebug("ENTER handleFile()")

	// Guard clause on directories, no need to handle them
	sourceInfo, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
		return err
	}

	if sourceInfo.IsDir() {
		logDebug("Ignore directory:", path)
		return nil
	}

	logDebug("Source file:", path)

	// Compute target path
	relativePath, err := filepath.Rel(*source, path)
	if err != nil {
		log.Panicln(err)
		return err
	}
	logDebug("Relative source path:", relativePath)

	targetPath := filepath.Join(*target, relativePath)
	logDebug("Target path:", targetPath)

	// Check if a target file exists or not
	if _, err := os.Lstat(targetPath); os.IsNotExist(err) {
		return handleNew(path, targetPath)
	} else {
		return handleExisting(path, targetPath)
	}

	return nil
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
	} else {
		// TODO: Wrap up or simplify response handling
		log.Println("[INFO] Existing Symlink points to ", evaluatedTargetPath, " ,replace with new symlink? [yN]")
		var response string
		_, err := fmt.Scanln(&response)
		if err != nil {
			log.Fatal(err)
		}
		okayResponses := []string{"y", "yes"}
		if contains(okayResponses, strings.ToLower(response)) {
			return symlink(sourcePath, targetPath)
		}
	}

	return nil
}

func handleExistingFile(sourcePath, targetPath string) error {
	logDebug("ENTER handleExistingFile()")

	equal, err := compareFiles(sourcePath, targetPath)
	if err != nil {
		log.Panicln(err)
		return err
	}

	if !equal {
		// TODO: Wrap up or simplify response handling
		log.Println("[INFO] Existing file differs, replace with symlink anyway? [Yn]")
		var response string
		_, err := fmt.Scanln(&response)
		if err != nil {
			log.Fatal(err)
		}
		okayResponses := []string{"n", "no"}
		if contains(okayResponses, strings.ToLower(response)) {
			return nil
		}
	}

	if err := os.Remove(targetPath); err != nil {
		log.Fatal(err)
		return err
	}

	return symlink(sourcePath, targetPath)

	return nil
}

func contains(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}

	return false
}

func compareFiles(lhs, rhs string) (bool, error) {
	logDebug("ENTER compareFiles()")

	return true, nil
}

func symlink(sourcePath, targetPath string) error {
	logDebug("ENTER symlink()")

	prepareDirectory(targetPath)

	if err := os.Symlink(sourcePath, targetPath); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func prepareDirectory(targetPath string) error {
	logDebug("ENTER prepareDirectory()")

	dirPath := filepath.Base(targetPath)

	// TODO: What is the correct FileMode to use instead of 0777?
	if err := os.MkdirAll(dirPath, 0777); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
