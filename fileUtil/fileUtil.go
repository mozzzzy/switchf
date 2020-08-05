package fileUtil

/*
 * Module Dependencies
 */

import (
	"io/ioutil"
	"os"
)


/*
 * Types
 */

/*
 * Constants and Package Scope Variables
 */

/*
 * Functions
 */

func createTmpFile(tmpDirPath string) (*os.File, error) {
	tmpFilePrefix := "switch_file."
	tmpFile, err := ioutil.TempFile(tmpDirPath, tmpFilePrefix)
	return tmpFile, err
}

func readFile(file *os.File) ([]byte, error) {
	bufSize := 1024
	buf := make([]byte, bufSize)
	var totalReadData []byte
	var totalReadSize int
	for {
		readSize, err := file.Read(buf)
		if readSize == 0 {
			break
		}
		if err != nil {
			return nil, err
		}
		totalReadData = append(totalReadData, buf...)
		totalReadSize += readSize
	}
	totalReadData = totalReadData[:totalReadSize]
	return totalReadData, nil
}

func writeFile(file *os.File, data []byte) error {
	_, err := file.Write(data)
	return err
}

func CreateTmpCopyFile(originalPath string) (string, error) {
	// Read original file data
	originalFile, err := os.Open(originalPath)
	if err != nil {
		return "", err
	}
	defer originalFile.Close()
	data, err := readFile(originalFile)
	if err != nil {
		return "", err
	}

	// Create new tmp file
	tmpDirPath := "/tmp"
	tmpFile, err := createTmpFile(tmpDirPath)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	// Write original data to tmp file
	if err = writeFile(tmpFile, data); err != nil {
		return "", err
	}

	// tmpFile.Name() returns absolute path.
	return tmpFile.Name(), nil
}
