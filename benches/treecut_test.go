package benches

// This file consists of benchmarks for the `trc` program.

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ezrantn/trc"
)

func createTestFiles(t *testing.T, dir string, filenames []string) {
	if t != nil {
		t.Helper()
	}

	for _, name := range filenames {
		path := filepath.Join(dir, name)
		f, err := os.Create(path)
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		f.Close()
	}
}

func BenchmarkMakePartitions(b *testing.B) {
	b.StopTimer() // Stop timing while setting up test files

	tempDir := os.TempDir()
	originalDir := filepath.Join(tempDir, "original")
	if err := os.Mkdir(originalDir, os.ModePerm); err != nil {
		b.Errorf("error creating directory: %v", err)
	}

	// Create 1000 test files *once* before running the benchmark loop
	files := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		files[i] = fmt.Sprintf("file%d.txt", i+1)
	}

	createTestFiles(nil, originalDir, files)

	outputDirs := []string{
		filepath.Join(tempDir, "partition1"),
		filepath.Join(tempDir, "partition2"),
	}

	for _, d := range outputDirs {
		if err := os.Mkdir(d, os.ModePerm); err != nil {
			b.Errorf("error creating directory: %v", err)
		}
	}

	config := trc.PartitionConfig{
		SourceDir:  originalDir,
		OutputDirs: outputDirs,
		BySize:     false,
	}

	b.StartTimer() // Start timing after setup

	for i := 0; i < b.N; i++ {
		if err := trc.MakePartitions(config); err != nil {
			b.Errorf("cannot create partition, something is wrong: %v", err)
		}
	}

	b.StopTimer() // Stop timing before cleanup

	// Cleanup (but only once at the end)
	os.RemoveAll(originalDir)
	for _, d := range outputDirs {
		os.RemoveAll(d)
	}
}
