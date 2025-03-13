package trc

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPartitionFiles(t *testing.T) {
	files := []string{"a.txt", "b.txt", "c.txt", "d.txt", "e.txt"}
	partitions := 2

	result := partitionFiles(files, partitions)

	if len(result) != partitions {
		t.Errorf("expected %d partitions, got %d", partitions, len(result))
	}

	// Check if files are roughly evenly distributed
	totalFiles := 0
	for _, p := range result {
		totalFiles += len(p)
	}

	if totalFiles != len(files) {
		t.Errorf("mismatrch in total files: expected %d, got %d", len(files), totalFiles)
	}
}

func TestMakePartitions(t *testing.T) {
	tempDir := t.TempDir()
	originalDir := filepath.Join(tempDir, "original")
	if err := os.Mkdir(originalDir, os.ModePerm); err != nil {
		t.Fatalf("error creating directory: %v", err)
	}

	testFiles := []struct {
		path    string
		content []byte
		mime    string
	}{
		{
			path:    filepath.Join(originalDir, "text.txt"),
			content: []byte("This is a text file with enough content to be detected correctly."),
			mime:    "text", // Text file
		},
		{
			path: filepath.Join(originalDir, "image.png"),
			content: []byte{
				0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
				0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
				0x08, 0x02, 0x00, 0x00, 0x00,
			},
			mime: "image", // Image file with PNG magic number
		},
		{
			path:    filepath.Join(originalDir, "empty.txt"),
			content: []byte{},
			mime:    "", // Empty file (will be skipped)
		},
		{
			path: filepath.Join(originalDir, "document.pdf"),
			content: []byte{
				0x25, 0x50, 0x44, 0x46, 0x2D, 0x31, 0x2E, 0x35, 0x0A, 0x25, 0xE2, 0xE3,
				0xCF, 0xD3, 0x0A, 0x0A, 0x31, 0x20, 0x30, 0x20, 0x6F, 0x62, 0x6A, 0x0A,
			},
			mime: "application", // PDF file (application/pdf)
		},
		{
			path: filepath.Join(originalDir, "audio.mp3"),
			content: []byte{
				0xFF, 0xFB, 0x90, 0x44, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			mime: "audio", // MP3 file (audio/mpeg)
		},
	}

	// Create files with the specified headers and content
	for _, file := range testFiles {
		err := os.WriteFile(file.path, file.content, 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file.path, err)
		}
	}

	outputDirs := []string{
		filepath.Join(tempDir, "partition1"),
		filepath.Join(tempDir, "partition2"),
	}

	// Create partition directories
	for _, d := range outputDirs {
		if err := os.Mkdir(d, os.ModePerm); err != nil {
			t.Fatalf("Error creating directory: %v", err)
		}
	}

	config := PartitionConfig{
		SourceDir:  originalDir,
		OutputDirs: outputDirs,
		BySize:     false,
		ByFile:     false,
	}

	err := MakePartitions(config)
	if err != nil {
		t.Fatalf("Partitioning failed: %v", err)
	}

	// Ensure files were partitioned
	counts := make(map[string]int)
	for _, d := range outputDirs {
		files, err := collectFilesWithMimeType(d)
		if err != nil {
			t.Fatalf("Failed to collect files from partition: %v", err)
		}

		counts[d] = len(files)
	}

	expectedTextFiles := 1  // text.txt
	expectedImageFiles := 1 // image.png
	expectedAudioFiles := 1 // audio.mp3

	if counts[outputDirs[0]] != expectedTextFiles && counts[outputDirs[1]] != expectedImageFiles+expectedAudioFiles {
		t.Errorf("Expected 1 text file in one partition and 2 image/audio files in the other. Got counts: %v", counts)
	}
	if counts[outputDirs[1]] != expectedTextFiles && counts[outputDirs[0]] != expectedImageFiles+expectedAudioFiles {
		t.Errorf("Expected 1 text file in one partition and 2 image/audio files in the other. Got counts: %v", counts)
	}
}

func TestPartitionFilesBySize(t *testing.T) {
	files := []fileInfo{
		{"a.txt", 100},
		{"b.txt", 200},
		{"c.txt", 300},
		{"d.txt", 400},
		{"e.txt", 500},
	}

	partitions := 2

	result := partitionFilesBySize(files, partitions)

	if len(result) != partitions {
		t.Errorf("expected %d partitions, got %d", partitions, len(result))
	}

	// Validate total size balancing
	size1, size2 := int64(0), int64(0)
	for _, f := range result[0] {
		size1 += f.size
	}

	for _, f := range result[1] {
		size2 += f.size
	}

	if abs(size1-size2) > 100 {
		t.Errorf("partitions not balanced: %d vs %d", size1, size2)
	}
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}

	return n
}
