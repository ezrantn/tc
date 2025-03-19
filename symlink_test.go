package trc

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateSymlinkTree(t *testing.T) {
	tempDir := t.TempDir()

	// Define test cases
	tests := []struct {
		name          string
		files         []string
		partitions    [][]string
		expectedLinks map[string][]string // map of partition directory to expected symlinks
		expectedErr   bool
	}{
		{
			name:          "single partition",
			files:         []string{"file1.txt", "file2.txt"},
			partitions:    [][]string{{"file1.txt", "file2.txt"}},
			expectedLinks: map[string][]string{filepath.Join(tempDir, "partition1"): {"file1.txt", "file2.txt"}},
			expectedErr:   false,
		},
		{
			name:          "multiple partitions",
			files:         []string{"file1.txt", "file2.txt"},
			partitions:    [][]string{{"file1.txt"}, {"file2.txt"}},
			expectedLinks: map[string][]string{filepath.Join(tempDir, "partition1"): {"file1.txt"}, filepath.Join(tempDir, "partition2"): {"file2.txt"}},
			expectedErr:   false,
		},
		{
			name:          "no files",
			files:         []string{},
			partitions:    [][]string{},
			expectedLinks: map[string][]string{},
			expectedErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup original directory and partition directories
			originalDir := filepath.Join(tempDir, "original")
			if err := os.Mkdir(originalDir, os.ModePerm); err != nil {
				t.Fatalf("error creating directory: %v", err)
			}

			defer os.RemoveAll(originalDir)

			// Create the test files in the original directory
			for _, file := range tt.files {
				if err := os.WriteFile(filepath.Join(originalDir, file), []byte("content"), os.ModePerm); err != nil {
					t.Fatalf("error creating file: %v", err)
				}
			}

			// Create partition directories
			var outputDirs []string
			for i := range tt.partitions {
				partitionDir := filepath.Join(tempDir, fmt.Sprintf("partition%d", i+1))
				if err := os.Mkdir(partitionDir, os.ModePerm); err != nil {
					t.Fatalf("error creating directory: %v", err)
				}
				outputDirs = append(outputDirs, partitionDir)

				defer os.RemoveAll(partitionDir)
			}

			// Create symlinks
			err := createSymlinkTree(tt.partitions, outputDirs)

			// Check for expected error
			if (err != nil) != tt.expectedErr {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}

			// Check if symlinks were created as expected
			for part, expectedFiles := range tt.expectedLinks {
				for _, file := range expectedFiles {
					linkPath := filepath.Join(part, filepath.Base(file))
					if _, err := os.Lstat(linkPath); err != nil {
						t.Errorf("symlink not created: %s", linkPath)
					}
				}
			}
		})
	}
}

func TestCreateSymlinkTreeBySize(t *testing.T) {
	tempDir := t.TempDir()

	// Define test cases
	tests := []struct {
		name          string
		files         []string
		partitions    [][]fileInfo
		expectedLinks map[string][]string // map of partition directory to expected symlinks
		expectedErr   bool
	}{
		{
			name:          "single partition",
			files:         []string{"file1.txt", "file2.txt"},
			partitions:    [][]fileInfo{{{path: "file1.txt", size: 100}}, {{path: "file2.txt", size: 200}}},
			expectedLinks: map[string][]string{filepath.Join(tempDir, "partition1"): {"file1.txt"}, filepath.Join(tempDir, "partition2"): {"file2.txt"}},
			expectedErr:   false,
		},
		{
			name:          "multiple partitions with different sizes",
			files:         []string{"file1.txt", "file2.txt", "file3.txt"},
			partitions:    [][]fileInfo{{{path: "file1.txt", size: 100}}, {{path: "file2.txt", size: 200}}, {{path: "file3.txt", size: 50}}},
			expectedLinks: map[string][]string{filepath.Join(tempDir, "partition1"): {"file1.txt"}, filepath.Join(tempDir, "partition2"): {"file2.txt"}, filepath.Join(tempDir, "partition3"): {"file3.txt"}},
			expectedErr:   false,
		},
		{
			name:          "no files",
			files:         []string{},
			partitions:    [][]fileInfo{},
			expectedLinks: map[string][]string{},
			expectedErr:   false,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup original directory and partition directories
			originalDir := filepath.Join(tempDir, "original")
			if err := os.Mkdir(originalDir, os.ModePerm); err != nil {
				t.Fatalf("error creating directory: %v", err)
			}

			defer os.RemoveAll(originalDir)

			for _, file := range tt.files {
				filePath := filepath.Join(originalDir, file)
				if err := os.WriteFile(filePath, []byte("content"), os.ModePerm); err != nil {
					t.Fatalf("error creating file: %v", err)
				}
			}

			// Create partition directories
			var outputDirs []string
			for i := range tt.partitions {
				partitionDir := filepath.Join(tempDir, fmt.Sprintf("partition%d", i+1))
				if err := os.Mkdir(partitionDir, os.ModePerm); err != nil {
					t.Fatalf("error creating directory: %v", err)
				}

				outputDirs = append(outputDirs, partitionDir)

				defer os.RemoveAll(partitionDir)
			}

			// Create symlinks
			err := createSymlinkTreeBySize(tt.partitions, outputDirs)

			// Check for expected error
			if (err != nil) != tt.expectedErr {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}

			// Check if symlinks were created as expected
			for part, expectedFiles := range tt.expectedLinks {
				for _, file := range expectedFiles {
					linkPath := filepath.Join(part, filepath.Base(file))
					info, err := os.Lstat(linkPath)
					if err != nil {
						t.Errorf("symlink not created: %s", linkPath)
					} else if info.Mode()&os.ModeSymlink == 0 {
						t.Errorf("expected symlink, but found regular file: %s", linkPath)
					}
				}
			}
		})
	}
}

func TestRemoveSymlinkTree(t *testing.T) {
	tempDir := t.TempDir()

	// Define test cases
	tests := []struct {
		name                string
		files               []string
		symlinks            []string
		directories         []string
		expectedErr         bool
		shouldSymlinkExist  bool
		shouldOriginalExist bool
	}{
		{
			name:                "single symlink",
			files:               []string{"testfile.txt"},
			symlinks:            []string{"symlink"},
			directories:         []string{tempDir},
			expectedErr:         false,
			shouldSymlinkExist:  false,
			shouldOriginalExist: true,
		},
		{
			name:                "multiple symlinks",
			files:               []string{"testfile.txt", "anotherfile.txt"},
			symlinks:            []string{"symlink1", "symlink2"},
			directories:         []string{tempDir},
			expectedErr:         false,
			shouldSymlinkExist:  false,
			shouldOriginalExist: true,
		},
		{
			name:                "no symlinks",
			files:               []string{"testfile.txt"},
			symlinks:            []string{},
			directories:         []string{tempDir},
			expectedErr:         false,
			shouldSymlinkExist:  false,
			shouldOriginalExist: true,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup original directory and test files
			originalDir := filepath.Join(tempDir, "original")
			if err := os.Mkdir(originalDir, os.ModePerm); err != nil {
				t.Fatalf("failed to create directory %s: %v", originalDir, err)
			}

			defer os.RemoveAll(originalDir)

			// Create files in the original directory
			var testFilePaths []string
			for _, file := range tt.files {
				filePath := filepath.Join(originalDir, file)
				if err := os.WriteFile(filePath, []byte("content"), os.ModePerm); err != nil {
					t.Fatalf("error creating file: %v", err)
				}
				testFilePaths = append(testFilePaths, filePath)
			}

			// Create symlinks in the temp directory
			for _, symlink := range tt.symlinks {
				symlinkPath := filepath.Join(tempDir, symlink)
				if err := os.Symlink(testFilePaths[0], symlinkPath); err != nil {
					t.Fatalf("failed to create symlink %s -> %s: %v", symlinkPath, testFilePaths[0], err)
				}
			}

			// Ensure symlinks exist before removal
			for _, symlink := range tt.symlinks {
				linkPath := filepath.Join(tempDir, symlink)
				if _, err := os.Lstat(linkPath); err != nil {
					t.Fatalf("symlink should exist before removal: %v", err)
				}
			}

			// Call removeSymlinkTree
			err := removeSymlinkTree(tt.directories)
			if (err != nil) != tt.expectedErr {
				t.Errorf("expected error: %v, got: %v", tt.expectedErr, err)
			}

			// Verify that symlinks have been removed
			for _, symlink := range tt.symlinks {
				linkPath := filepath.Join(tempDir, symlink)
				if _, err := os.Lstat(linkPath); !tt.shouldSymlinkExist && !os.IsNotExist(err) {
					t.Errorf("expected symlink %s to be removed, but it still exists", linkPath)
				}
			}

			// Verify that the original file(s) still exist
			for _, filePath := range testFilePaths {
				if _, err := os.Stat(filePath); err != nil && !os.IsNotExist(err) {
					t.Errorf("expected original file %s to exist, but got error: %v", filePath, err)
				}
			}
		})
	}
}
