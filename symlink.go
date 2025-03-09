package treecut

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// CreateSymlinkTree creates symlinks in partition directories
func createSymlinkTree(files [][]string, outputDirs []string) error {
	for i, partition := range files {
		for _, file := range partition {
			// Symlink path inside partition
			linkPath := filepath.Join(outputDirs[i], filepath.Base(file))

			// Remove existing symlink or file before creating a new one
			if _, err := os.Lstat(linkPath); err == nil {
				if err := os.Remove(linkPath); err != nil {
					return fmt.Errorf("failed to remove existing symlink %s: %w", linkPath, err)
				}
			} else if !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("failed to check symlink existence %s: %w", linkPath, err)
			}

			// Ensure the partition directory exists
			if err := os.MkdirAll(filepath.Dir(linkPath), os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", filepath.Dir(linkPath), err)
			}

			// Create a symlink
			if err := os.Symlink(file, linkPath); err != nil {
				return fmt.Errorf("failed to create symlink from %s to %s: %w", file, linkPath, err)
			}
		}
	}

	return nil
}

func createSymlinkTreeBySize(files [][]fileInfo, outputDirs []string) error {
	for i, partition := range files {
		for _, file := range partition {
			linkPath := filepath.Join(outputDirs[i], filepath.Base(file.path))

			// Check if the symlink already exists and remove it
			if _, err := os.Lstat(linkPath); err == nil {
				if err := os.Remove(linkPath); err != nil {
					return fmt.Errorf("failed to remove existing symlink %s: %w", linkPath, err)
				}
			} else if !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("failed to check symlink existence %s: %w", linkPath, err)
			}

			// Ensure the partition directory exists
			if err := os.MkdirAll(filepath.Dir(linkPath), os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", filepath.Dir(linkPath), err)
			}

			// Create a symlink
			if err := os.Symlink(file.path, linkPath); err != nil {
				return fmt.Errorf("failed to create symlink from %s to %s: %w", file.path, linkPath, err)
			}
		}
	}

	return nil
}

func removeSymlinkTree(outputDirs []string) error {
	for _, dir := range outputDirs {
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
	}

	return nil
}
