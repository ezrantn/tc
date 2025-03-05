package treecut

import (
	"io/fs"
	"path/filepath"
)

type FileInfo struct {
	Path string
	Size int64
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
				filesChan <- path
			}

			return nil
		})
	}()

	// Collect results
	var files []string
	for file := range filesChan {
		files = append(files, file)
	}

	// Return any errors encountered
	if err := <-errChan; err != nil {
		return nil, err
	}

	return files, nil
}

// collectFilesWithSize collects file paths and sizes
func collectFilesWithSize(sourceDir string) ([]FileInfo, error) {
	var files []FileInfo
	err := filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}

			files = append(files, FileInfo{Path: path, Size: info.Size()})
		}

		return nil
	})

	return files, err
}
