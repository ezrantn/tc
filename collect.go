package treecut

import (
	"errors"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
)

type fileInfo struct {
	path string
	size int64
}

const (
	characterFilter = `[\x00-\x1F\\/:*?"<>|@!]` // Some FAT systems don't allow @ and ! in filenames
	defaultName     = "file"
	maxLength       = 255
)

var (
	characterFilterRegex = regexp.MustCompile(characterFilter)

	dosReservedNames = map[string]struct{}{
		"CON": {}, "PRN": {}, "AUX": {}, "NUL": {}, "CLOCK$": {}, "CONFIG$": {}, "SCREEN$": {}, "$IDLE$": {},
		"COM0": {}, "COM1": {}, "COM2": {}, "COM3": {}, "COM4": {}, "COM5": {}, "COM6": {}, "COM7": {}, "COM8": {}, "COM9": {},
		"LPT0": {}, "LPT1": {}, "LPT2": {}, "LPT3": {}, "LPT4": {}, "LPT5": {}, "LPT6": {}, "LPT7": {}, "LPT8": {}, "LPT9": {},
	}

	windowsReservedNames = map[string]struct{}{
		"$Mft": {}, "$MftMirr": {}, "$LogFile": {}, "$Volume": {}, "$AttrDef": {}, "$Bitmap": {}, "$Boot": {}, "$BadClus": {},
		"$Secure": {}, "$Upcase": {}, "$Extend": {}, "$Quota": {}, "$ObjId": {}, "$Reparse": {},
	}
)

// isValidFileName ensures a file is not using reserved names and meets system constraints
func isValidFileName(name string) bool {
	name = strings.ToUpper(strings.TrimSpace(name))

	if len(name) > maxLength || characterFilterRegex.MatchString(name) {
		return false
	}

	if _, exists := dosReservedNames[name]; exists {
		return false
	}

	if _, exists := windowsReservedNames[name]; exists {
		return false
	}

	return true
}

// collectFiles walks through the source directory and returns all file paths
func collectFiles(sourceDir string) ([]string, error) {
	filesChan := make(chan string, 100) // Buffered channel to reduce contention
	errChan := make(chan error, 1)

	// Walk the directory in a separate goroutine
	go func() {
		defer close(filesChan)
		errChan <- filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				baseName := filepath.Base(path)
				if !isValidFileName(baseName) {
					return errors.New("invalid file name: " + baseName)
				}

				filesChan <- path
			}
			return nil
		})

		close(errChan)
	}()

	// Collect results
	var files []string
	for file := range filesChan {
		files = append(files, file)
	}

	// Return any errors encountered
	if err, ok := <-errChan; ok && err != nil {
		return nil, err
	}

	return files, nil
}

// collectFilesWithSize collects file paths and sizes
func collectFilesWithSize(sourceDir string) ([]fileInfo, error) {
	var files []fileInfo
	err := filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}

			// Validate filename before adding it
			if !isValidFileName(filepath.Base(path)) {
				return errors.New("invalid file name: " + path)
			}

			files = append(files, fileInfo{path: path, size: info.Size()})
		}

		return nil
	})

	return files, err
}
