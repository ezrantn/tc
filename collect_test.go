package treecut

import (
	"os"
	"path/filepath"
	"strings"
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

func TestCollectFilesBySize(t *testing.T) {

}

func TestIsValidName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid filename", "myfile.txt", true},
		{"Empty string", "", false},
		{"Exceeds max length", strings.Repeat("a", maxLength+1), false},
		{"Contains invalid character", "file:name.txt", false},
		{"DOS reserved name", "CON", false},
		{"DOS reserved name case insensitive", "con", false},
		{"Windows reserved name", "$MFT", false},
		{"Windows reserved name case insensitive", "$mft", false},
		{"Valid long filename", strings.Repeat("a", maxLength), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, _ := isValidFileName(test.input); got != test.expected {
				t.Errorf("isValidFileName(%q) = %v; want %v", test.input, got, test.expected)
			}
		})
	}
}
