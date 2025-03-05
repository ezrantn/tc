package treecut

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
		t.Errorf("mismatch in total files: expected %d, got %d", len(files), totalFiles)
	}
}

func TestMakePartitions(t *testing.T) {
	tempDir := t.TempDir()
	originalDir := filepath.Join(tempDir, "original")
	if err := os.Mkdir(originalDir, os.ModePerm); err != nil {
		t.Fatalf("error creating directory: %v", err)
	}

	createTestFiles(t, originalDir, []string{"file1.txt", "file2.txt", "file3.txt", "file4.txt"})

	outputDirs := []string{
		filepath.Join(tempDir, "partition1"),
		filepath.Join(tempDir, "partition2"),
	}

	for _, d := range outputDirs {
		if err := os.Mkdir(d, os.ModePerm); err != nil {
			t.Fatalf("error creating directory: %v", err)
		}
	}

	config := PartitionConfig{
		SourceDir:  originalDir,
		OutputDirs: outputDirs,
	}

	err := MakePartitions(config)
	if err != nil {
		t.Fatalf("partition failed: %v", err)
	}

	// Ensure files were partitioned
	counts := make(map[string]int)
	for _, d := range outputDirs {
		files, _ := collectFiles(d)
		counts[d] = len(files)
	}

	expectedPerPartition := 2
	for _, count := range counts {
		if count != expectedPerPartition {
			t.Errorf("expected %d files per partition, got %d", expectedPerPartition, count)
		}
	}
}

func TestPartitionFilesBySize(t *testing.T) {
	files := []FileInfo{
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
		size1 += f.Size
	}

	for _, f := range result[1] {
		size2 += f.Size
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
