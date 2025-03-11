package trc

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

// Helper function to get map keys as slice
func getMapKeys(m map[string][]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func TestCollectFiles(t *testing.T) {
	tests := []struct {
		name          string
		files         []string
		expectedCount int
		expectError   bool
		setup         func(t *testing.T) string
	}{
		{
			name: "Three files",
			files: []string{
				"file1.txt",
				"file2.txt",
				"file3.txt",
			},
			expectedCount: 3,
			expectError:   false,
		},
		{
			name:          "Empty directory",
			files:         []string{},
			expectedCount: 0,
			expectError:   false,
		},
		{
			name: "Files with invalid names",
			files: []string{
				"valid1.txt",
				"valid2.log",
				"invalid..txt", // Should be ignored
				"invalid ",     // Should be ignored
			},
			expectedCount: 0,
			expectError:   true,
		},
		{
			name: "Non-existent directory",
			setup: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nonexistent")
			},
			expectError: true,
		},
		{
			name: "Permission denied",
			setup: func(t *testing.T) string {
				tempDir := t.TempDir()
				restrictedDir := filepath.Join(tempDir, "restricted")

				if err := os.Mkdir(restrictedDir, 0000); err != nil {
					t.Fatalf("Failed to create restricted directory: %v", err)
				}

				return restrictedDir
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testDir string

			if tt.setup != nil {
				testDir = tt.setup(t)
			} else {
				testDir = t.TempDir()

				// Create test files
				for _, file := range tt.files {
					path := filepath.Join(testDir, file)
					f, err := os.Create(path)
					if err != nil {
						t.Fatalf("Failed to create test file %s: %v", file, err)
					}

					f.Close()
				}
			}

			// Run collectFiles
			files, err := collectFiles(testDir)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error but got none")
				}

				return
			}

			// Validate file count
			if len(files) != tt.expectedCount {
				t.Errorf("Expected %d files, got %d", tt.expectedCount, len(files))
			}
		})
	}
}

func TestCollectFilesBySize(t *testing.T) {
	tests := []struct {
		name        string
		files       []string
		size        int64
		expectError bool
		setup       func(t *testing.T) string
	}{
		{
			name: "Small size file (10 bytes)",
			files: []string{
				"small1.txt",
				"small2.txt",
				"small3.txt",
			},
			size:        10,
			expectError: false,
		},
		{
			name: "Medium size file (512 bytes)",
			files: []string{
				"medium1.txt",
				"medium2.txt",
				"medium3.txt",
			},
			size:        512,
			expectError: false,
		},
		{
			name: "Big size file (1024 bytes / 1MB)",
			files: []string{
				"big1.txt",
				"big2.txt",
				"big3.txt",
			},
			size:        1024,
			expectError: false,
		},
		{
			name: "Files with invalid names",
			files: []string{
				"invalid ",
				"invalid..txt",
				"valid1.txt",
				"invalid...",
			},
			expectError: true,
		},
		{
			name: "Non-existent directory",
			setup: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "nonexistent")
			},
			expectError: true,
		},
		{
			name: "Permission denied",
			setup: func(t *testing.T) string {
				tempDir := t.TempDir()
				restrictedDir := filepath.Join(tempDir, "restricted")

				if err := os.Mkdir(restrictedDir, 0000); err != nil {
					t.Fatalf("Failed to create restricted directory: %v", err)
				}

				return restrictedDir
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var testDir string

			if tt.setup != nil {
				testDir = tt.setup(t)
			} else {
				testDir = t.TempDir()

				// Create test files with the specified size
				for _, file := range tt.files {
					path := filepath.Join(testDir, file)
					f, err := os.Create(path)
					if err != nil {
						t.Fatalf("Failed to create test file %s: %v", file, err)
					}

					// Write dummy data to matrch the expected file size
					if _, err := f.Write(make([]byte, tt.size)); err != nil {
						t.Fatalf("Failed to write to file %s: %v", file, err)
					}

					f.Close()
				}
			}

			// Run the function
			files, err := collectFilesWithSize(testDir)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("collectFilesWithSize failed: %v", err)
			}

			// Check the correct number of files
			if len(files) != len(tt.files) {
				t.Errorf("Expected %d files, got %d", len(tt.files), len(files))
			}

			// Verify file sizes
			expectedSizes := make(map[string]int64)
			for _, f := range tt.files {
				expectedSizes[filepath.Join(testDir, f)] = tt.size
			}

			for _, file := range files {
				expectedSize, exists := expectedSizes[file.path]
				if !exists {
					t.Errorf("Unexpected file collected: %s", file.path)
				} else if file.size != expectedSize {
					t.Errorf("File %s has size %d, expected %d", file.path, file.size, expectedSize)
				}
			}
		})
	}
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
		{"Has trailing dot", "file.txt.", false},
		{"Has trailing space", "file.txt ", false},
		{"Has multiple trailing dots", "file....", false},
		{"Has multiple trailing spaces", "file    ", false},
		{"Has mixed trailing dots and spaces", "file. . ", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got, _ := isValidFileName(test.input); got != test.expected {
				t.Errorf("isValidFileName(%q) = %v; want %v", test.input, got, test.expected)
			}
		})
	}
}

func TestCollectFilesWithMimeType(t *testing.T) {
	testDir := t.TempDir()

	// Create subdirectory for audio file
	subDir := filepath.Join(testDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	// Create test files with more complete headers/content for accurate MIME detection
	testFiles := []struct {
		path    string
		content []byte
		mime    string // Expected MIME type
	}{
		{
			path:    filepath.Join(testDir, "text.txt"),
			content: []byte("This is a text file with enough content to be detected correctly."),
			mime:    "text",
		},
		{
			path: filepath.Join(testDir, "image.png"),
			content: []byte{
				0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
				0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
				0x08, 0x02, 0x00, 0x00, 0x00,
			},
			mime: "image",
		},
		{
			path:    filepath.Join(testDir, "empty.txt"),
			content: []byte{},
			mime:    "",
		},
		{
			path: filepath.Join(testDir, "document.pdf"),
			content: []byte{
				0x25, 0x50, 0x44, 0x46, 0x2D, 0x31, 0x2E, 0x35, 0x0A, 0x25, 0xE2, 0xE3,
				0xCF, 0xD3, 0x0A, 0x0A, 0x31, 0x20, 0x30, 0x20, 0x6F, 0x62, 0x6A, 0x0A,
			},
			mime: "application",
		},
		{
			path: filepath.Join(subDir, "audio.mp3"),
			content: []byte{
				0xFF, 0xFB, 0x90, 0x44, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			mime: "audio",
		},
	}

	for _, tf := range testFiles {
		if err := os.WriteFile(tf.path, tf.content, 0644); err != nil {
			t.Fatalf("failed to create test file %s: %v", tf.path, err)
		}
	}

	result, err := collectFilesWithMimeType(testDir)
	if err != nil {
		t.Fatalf("collectFilesWithMimeType returned an error: %v", err)
	}

	expectedCategories := []string{"text", "image", "application", "audio"}
	missingCategories := []string{}

	for _, category := range expectedCategories {
		if _, exists := result[category]; !exists {
			missingCategories = append(missingCategories, category)
		}
	}

	if len(missingCategories) > 0 {
		t.Errorf("Missing expected categories: %v", missingCategories)
	}

	// Verify empty file was skipped
	for _, files := range result {
		for _, file := range files {
			if filepath.Base(file) == "empty.txt" {
				t.Errorf("empty.txt should be skipped but was found in results")
			}
		}
	}
}

func TestFileInCategory(t *testing.T) {
	testDir := t.TempDir()

	// Create a text file
	textFile := filepath.Join(testDir, "sample.txt")
	if err := os.WriteFile(textFile, []byte("This is a sample text file"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	result, err := collectFilesWithMimeType(testDir)
	if err != nil {
		t.Fatalf("collectFilesWithMimeType returned an error: %v", err)
	}

	// Check if the file is in the text category
	if files, ok := result["text"]; ok {
		found := false
		for _, file := range files {
			if strings.HasSuffix(file, "sample.txt") {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("sample.txt should be in the 'text' category")
		}
	} else {
		t.Errorf("Expected 'text' category not found in results")
	}
}
