package treecut

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateSymlinkTree(t *testing.T) {
	tempDir := t.TempDir()

	originalDir := filepath.Join(tempDir, "original")
	os.Mkdir(originalDir, os.ModePerm)
	createTestFiles(t, originalDir, []string{"file1.txt", "file2.txt"})

	// Create partitions
	partition1 := filepath.Join(tempDir, "partition1")
	partition2 := filepath.Join(tempDir, "partition2")
	os.Mkdir(partition1, os.ModePerm)
	os.Mkdir(partition2, os.ModePerm)

	// Partition files
	partitions := [][]string{
		{filepath.Join(originalDir, "file1.txt")},
		{filepath.Join(originalDir, "file2.txt")},
	}

	// Create symlinks
	outputDirs := []string{partition1, partition2}
	err := createSymlinkTree(partitions, outputDirs)
	if err != nil {
		t.Fatalf("create symlink tree failed: %v", err)
	}

	// Check if symlinks exists
	for i, part := range outputDirs {
		for _, file := range partitions[i] {
			linkPath := filepath.Join(part, filepath.Base(file))
			if _, err := os.Lstat(linkPath); err != nil {
				t.Errorf("symlink not created: %s", linkPath)
			}
		}
	}
}
