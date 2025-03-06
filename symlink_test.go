package treecut

import (
	"os"
	"path/filepath"
	"testing"
)

func createTestFile(t *testing.T, dir, name string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	f.Close()
	return path
}

func TestCreateSymlinkTree(t *testing.T) {
	tempDir := t.TempDir()
	var err error

	originalDir := filepath.Join(tempDir, "original")
	if err = os.Mkdir(originalDir, os.ModePerm); err != nil {
		t.Fatalf("error creating directory: %v", err)
	}

	createTestFiles(t, originalDir, []string{"file1.txt", "file2.txt"})

	// Create partitions
	partition1 := filepath.Join(tempDir, "partition1")
	partition2 := filepath.Join(tempDir, "partition2")
	if err = os.Mkdir(partition1, os.ModePerm); err != nil {
		t.Fatalf("error creating directory: %v", err)
	}

	if err = os.Mkdir(partition2, os.ModePerm); err != nil {
		t.Fatalf("error creating directory: %v", err)
	}

	// Partition files
	partitions := [][]string{
		{filepath.Join(originalDir, "file1.txt")},
		{filepath.Join(originalDir, "file2.txt")},
	}

	// Create symlinks
	outputDirs := []string{partition1, partition2}
	err = createSymlinkTree(partitions, outputDirs)
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

func TestCreateSymlinkTreeBySize(t *testing.T) {
	tempDir := t.TempDir()

	// Create original directory
	originalDir := filepath.Join(tempDir, "original")
	if err := os.Mkdir(originalDir, os.ModePerm); err != nil {
		t.Fatalf("error creating directory: %v", err)
	}

	// Create test files
	file1 := createTestFile(t, originalDir, "file1.txt")
	file2 := createTestFile(t, originalDir, "file2.txt")

	// Create partition output directories
	partition1 := filepath.Join(tempDir, "partition1")
	partition2 := filepath.Join(tempDir, "partition2")
	if err := os.Mkdir(partition1, os.ModePerm); err != nil {
		t.Fatalf("error creating directory: %v", err)
	}
	
	if err := os.Mkdir(partition2, os.ModePerm); err != nil {
		t.Fatalf("error creating directory: %v", err)
	}

	// Partition files
	partitions := [][]fileInfo{
		{{path: file1, size: 100}}, // Fake size
		{{path: file2, size: 200}},
	}

	// Create symlinks
	outputDirs := []string{partition1, partition2}
	err := createSymlinkTreeSize(partitions, outputDirs)
	if err != nil {
		t.Fatalf("createSymlinkTreeSize failed: %v", err)
	}

	// Check if symlinks exist
	for i, part := range outputDirs {
		for _, file := range partitions[i] {
			linkPath := filepath.Join(part, filepath.Base(file.path))
			info, err := os.Lstat(linkPath)
			if err != nil {
				t.Errorf("symlink not created: %s", linkPath)
			} else if info.Mode()&os.ModeSymlink == 0 {
				t.Errorf("expected symlink, but found regular file: %s", linkPath)
			}
		}
	}
}
