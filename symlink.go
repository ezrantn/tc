package trc

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// createSymlinks handles the creation of symlinks for the provided files and output directories.
// The `getPath` function is used to extract the file path from each element of the files slice.
func createSymlinks[T any](files [][]T, outputDirs []string, getPath func(T) string) error {
	for i, partition := range files {
		for _, file := range partition {
			filePath := getPath(file)
			linkPath := filepath.Join(outputDirs[i], filepath.Base(filePath))

			// Remove existing symlink or file before creating a new one
			if err := removeExistingSymlink(linkPath); err != nil {
				return err
			}

			// Ensure the partition directory exists
			if err := ensureDirectory(filepath.Dir(linkPath)); err != nil {
				return err
			}

			// Create a symlink
			if err := os.Symlink(filePath, linkPath); err != nil {
				return fmt.Errorf("failed to create symlink from %s to %s: %w", filePath, linkPath, err)
			}
		}
	}
	return nil
}

// removeExistingSymlink removes an existing symlink or file, if it exists.
func removeExistingSymlink(linkPath string) error {
	if _, err := os.Lstat(linkPath); err == nil {
		if err := os.Remove(linkPath); err != nil {
			return fmt.Errorf("failed to remove existing symlink %s: %w", linkPath, err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to check symlink existence %s: %w", linkPath, err)
	}
	return nil
}

// ensureDirectory ensures that the directory exists, creating it if necessary.
func ensureDirectory(dir string) error {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	return nil
}

// createSymlinkTree creates symlinks in partition directories.
func createSymlinkTree(files [][]string, outputDirs []string) error {
	return createSymlinks(files, outputDirs, func(f string) string {
		return f
	})
}

// createSymlinkTreeBySize creates symlinks in partition directories, based on file sizes.
func createSymlinkTreeBySize(files [][]fileInfo, outputDirs []string) error {
	return createSymlinks(files, outputDirs, func(f fileInfo) string {
		return f.path
	})
}

// createSymlinkWithMimeType creates symlinks for files based on their MIME type, organizing them into categories.
func createSymlinkWithMimeType(mimeMap map[string][]string, destDir string) error {
	for category, files := range mimeMap {
		categoryFolder := filepath.Join(destDir, category)
		if err := ensureDirectory(categoryFolder); err != nil {
			return err
		}

		for _, file := range files {
			linkPath := filepath.Join(categoryFolder, filepath.Base(file))
			if err := os.Symlink(file, linkPath); err != nil {
				return fmt.Errorf("failed to create symlink for %s: %w", file, err)
			}
		}
	}

	return nil
}

// removeSymlinkTree removes all symlinks within the provided directories.
func removeSymlinkTree(outputDirs []string) error {
	for _, dir := range outputDirs {
		if err := walkAndRemoveSymlinks(dir); err != nil {
			return err
		}
	}
	return nil
}

// walkAndRemoveSymlinks walks through a directory tree and removes all symlinks found.
func walkAndRemoveSymlinks(dir string) error {
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %s: %w", path, err)
		}

		// Check if the file is a symlink
		if info.Mode()&os.ModeSymlink != 0 {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove symlink %s: %w", path, err)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to remove symlinks in directory %s: %w", dir, err)
	}
	return nil
}
