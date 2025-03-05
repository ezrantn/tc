package treecut

import (
	"os"
	"path/filepath"
)

// CreateSymlinkTree creates symlinks in partition directories
func createSymlinkTree(files [][]string, outputDirs []string) error {
	for i, partition := range files {
		for _, file := range partition {
			// Symlink path inside partition
			linkPath := filepath.Join(outputDirs[i], filepath.Base(file))

			// Ensure the partition directory exists
			if err := os.MkdirAll(filepath.Dir(linkPath), os.ModePerm); err != nil {
				return err
			}

			// Create a symlink
			err := os.Symlink(file, linkPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func createSymlinkTreeSize(files [][]fileInfo, outputDirs []string) error {
	for i, partition := range files {
		for _, file := range partition {
			linkPath := filepath.Join(outputDirs[i], filepath.Base(file.path))
			if err := os.MkdirAll(filepath.Dir(linkPath), os.ModePerm); err != nil {
				return err
			}
			if err := os.Symlink(file.path, linkPath); err != nil {
				return err
			}
		}
	}
	return nil
}
