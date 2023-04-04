package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	// check if dummy file path and directory prefix are provided as arguments
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: go run main.go dummy_file_path directory_prefix")
		os.Exit(1)
	}

	dummyFilePath := os.Args[1]
	directoryPrefix := os.Args[2]

	// open the dummy file
	dummyFile, err := os.Open(dummyFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error opening dummy file:", err)
		os.Exit(1)
	}
	defer dummyFile.Close()

	// read file names from stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fileName := scanner.Text()

		// construct the full file path using the directory prefix
		fullPath := filepath.Join(directoryPrefix, fileName)

		// open the file
		file, err := os.Create(fullPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating file:", err)
			continue
		}

		// copy the contents of the dummy file to the file
		_, err = io.Copy(file, dummyFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error copying file contents:", err)
			continue
		}

		// set the permissions of the new file to match the original file
		fileInfo, err := os.Stat(dummyFilePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting file info:", err)
			continue
		}
		err = os.Chmod(fullPath, fileInfo.Mode())
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error setting file permissions:", err)
			continue
		}

		// fmt.Println("Created", fullPath)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading from stdin:", err)
	}

	// reset dummy file offset to the beginning
	if _, err := dummyFile.Seek(0, io.SeekStart); err != nil {
		fmt.Fprintln(os.Stderr, "Error resetting dummy file offset:", err)
		os.Exit(1)
	}
}
