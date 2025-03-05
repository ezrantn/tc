package treecut

import (
	"os"
	"path/filepath"
	"testing"
)

// Helper function to create temporary test files
func createTestFiles(t *testing.T, dir string, filenames []string) {
	t.Helper()
	for _, name := range filenames {
		path := filepath.Join(dir, name)
		f, err := os.Create(path)
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		f.Close()
	}
}

func TestCollectFiles(t *testing.T) {
	tempDir := t.TempDir()
	createTestFiles(t, tempDir, []string{"file1.txt", "file2.txt", "file3.txt"})

	files, err := collectFiles(tempDir)
	if err != nil {
		t.Fatalf("collect files failed: %v", err)
	}

	expected := 3
	if len(files) != expected {
		t.Errorf("expected %d files, got %d", expected, len(files))
	}
}
