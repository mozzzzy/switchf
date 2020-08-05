package main

/*
 * Module Dependencies
 */

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mozzzzy/arguments"
	"github.com/mozzzzy/arguments/argumentOption"
	"github.com/mozzzzy/cui/v2"
	"github.com/mozzzzy/switchf/fileUtil"
)

/*
 * Types
 */

/*
 * Constants and Package Scope Variables
 */

var DEBUG bool = false

/*
 * Functions
 */

func debug(msg string) {
	if DEBUG {
		cui.Debug(msg)
	}
}

func configArgOptions() (arguments.Args, error) {
	var args arguments.Args
	err := args.AddOptions([]argumentOption.Option{
		{
			LongKey:     "help",
			ShortKey:    "h",
			Description: "Show help message and exit.",
		},
		{
			LongKey:     "file",
			ShortKey:    "f",
			Description: "Target file path.",
			ValueType:   "string",
		},
		{
			LongKey:     "verbose",
			ShortKey:    "v",
			Description: "Print debug message.",
		},
	})

	return args, err
}

func parseArgs() (arguments.Args, error) {
	args, err := configArgOptions()
	if err != nil {
		return args, err
	}
	err = args.Parse()
	return args, err
}

func getSwitchFilePaths(targetPath string) ([]string, error) {
	targetDir := filepath.Dir(targetPath)
	targetFilePrefix := filepath.Base(targetPath) + "_"

	var targetFiles []string
	files, err := ioutil.ReadDir(targetDir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if !strings.HasPrefix(file.Name(), targetFilePrefix) {
			continue
		}
		if file.IsDir() {
			continue
		}
		targetFiles = append(targetFiles, targetDir+"/"+file.Name())
	}
	return targetFiles, nil
}

func switchFile(targetPath, switchFilePath string) error {
	tmpCopyFilePath, err := fileUtil.CreateTmpCopyFile(switchFilePath)
	if err != nil {
		return err
	}

	if err := os.Rename(tmpCopyFilePath, targetPath); err != nil {
		os.Remove(tmpCopyFilePath)
		return err
	}
	return nil
}

func main() {
	debug("switchf start.")

	args, err := parseArgs()
	if err != nil {
		cui.Error("Failed to parse arguments.")
		cui.Error(err.Error())
		return
	}

	if args.IsSet("help") {
		cui.Message(args.String(), []string{})
		return
	}

	if !args.IsSet("file") {
		cui.Error("Please specify --file, -f TARGET_FILE_PATH")
		return
	}

	if args.IsSet("verbose") {
		DEBUG = true
		debug("--verbose, -v option is set.")
	}

	targetPath, err := args.GetString("file")
	if err != nil {
		cui.Error("Failed to get file path from --file -f option.")
		cui.Error(err.Error())
		return
	}
	debug("Target file path: \"" + targetPath + "\".")

	switchFilePaths, err := getSwitchFilePaths(targetPath)
	if err != nil {
		cui.Error("Failed to get switch file paths.")
		cui.Error(err.Error())
	}
	if len(switchFilePaths) == 0 {
		cui.Error("No switch candidates are found.")
		return
	}

	switchPathIndex, canceled := cui.List("Switch candidates.", switchFilePaths)
	if canceled {
		cui.Warn("Canceled.")
		return
	}
	debug("Switch to \"" + switchFilePaths[switchPathIndex] + "\".")

	if err := switchFile(targetPath, switchFilePaths[switchPathIndex]); err != nil {
		cui.Error("Failed to switch file.")
		cui.Error(err.Error())
		return
	}

	cui.Info("File switched successfully.")
	debug("switchf finish successfully.")
}
