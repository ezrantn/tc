package benches

// This file consists of benchmarks for the `trc` program.

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ezrantn/trc"
)

func BenchmarkMakePartitions(b *testing.B) {
	// Define test cases for the benchmark
	tests := []struct {
		name              string
		numFiles          int
		numPartitions     int
		filesPerPartition int
	}{
		{
			name:              "1000 files, 2 partitions",
			numFiles:          1000,
			numPartitions:     2,
			filesPerPartition: 500,
		},
		{
			name:              "5000 files, 5 partitions",
			numFiles:          5000,
			numPartitions:     5,
			filesPerPartition: 1000,
		},
		{
			name:              "10000 files, 10 partitions",
			numFiles:          10000,
			numPartitions:     10,
			filesPerPartition: 1000,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.StopTimer() // Stop timing while setting up test files

			tempDir := os.TempDir()
			originalDir := filepath.Join(tempDir, "original")
			if err := os.Mkdir(originalDir, os.ModePerm); err != nil {
				b.Errorf("error creating directory: %v", err)
			}

			// Create test files
			files := make([]string, tt.numFiles)
			for i := 0; i < tt.numFiles; i++ {
				files[i] = fmt.Sprintf("file%d.txt", i+1)
			}

			// Create the files directly
			for _, file := range files {
				filePath := filepath.Join(originalDir, file)
				if err := os.WriteFile(filePath, []byte("content"), os.ModePerm); err != nil {
					b.Fatalf("error creating file: %v", err)
				}
			}

			// Create partition directories
			outputDirs := make([]string, tt.numPartitions)
			for i := 0; i < tt.numPartitions; i++ {
				partitionDir := filepath.Join(tempDir, fmt.Sprintf("partition%d", i+1))
				if err := os.Mkdir(partitionDir, os.ModePerm); err != nil {
					b.Errorf("error creating partition directory: %v", err)
				}
				outputDirs[i] = partitionDir
			}

			// Set up the partition config
			config := trc.PartitionConfig{
				SourceDir:  originalDir,
				OutputDirs: outputDirs,
				BySize:     false,
			}

			b.StartTimer()

			// Run the benchmark
			for i := 0; i < b.N; i++ {
				if err := trc.MakePartitions(config); err != nil {
					b.Errorf("cannot create partition, something is wrong: %v", err)
				}
			}

			b.StopTimer()

			// Cleanup (but only once at the end)
			os.RemoveAll(originalDir)
			for _, d := range outputDirs {
				os.RemoveAll(d)
			}
		})
	}
}
