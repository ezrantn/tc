package trc

import (
	"errors"
	"fmt"
	"os"
	"sort"
)

type PartitionConfig struct {
	SourceDir  string   // Original directory
	OutputDirs []string // Partition directories
	BySize     bool     // Set to true to activate partition by size (largest -> smallest)
}

func MakePartitions(config PartitionConfig) error {
	if len(config.OutputDirs) == 0 {
		return errors.New("at least one output directory is required")
	}

	// Collect files
	if config.BySize {
		files, err := collectFilesWithSize(config.SourceDir)
		if err != nil {
			return fmt.Errorf("failed to collect files with size from %s: %w", config.SourceDir, err)
		}

		partitions := partitionFilesBySize(files, len(config.OutputDirs))
		if err := createSymlinkTreeBySize(partitions, config.OutputDirs); err != nil {
			return fmt.Errorf("failed to create symlink tree by size: %w", err)
		}
	} else {
		files, err := collectFiles(config.SourceDir)
		if err != nil {
			return fmt.Errorf("failed to collect files from %s: %w", config.SourceDir, err)
		}

		partitions := partitionFiles(files, len(config.OutputDirs))
		if err := createSymlinkTree(partitions, config.OutputDirs); err != nil {
			return fmt.Errorf("failed to create symlink tree: %w", err)
		}
	}

	return nil
}

func RemovePartitions(outputDirs []string) error {
	// Remove the symlink first
	if err := removeSymlinkTree(outputDirs); err != nil {
		return fmt.Errorf("failed to remove symlink tree: %w", err)
	}

	// Remove the partition directories
	for _, dir := range outputDirs {
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("failed to remove partition directory %s: %w", dir, err)
		}
	}

	return nil
}

// partitionFiles splits a list of file paths into `partitions` equal groups
func partitionFiles(files []string, partitions int) [][]string {
	if partitions <= 0 {
		return nil
	}

	// Preallocate slices with estimated capacity
	result := make([][]string, partitions)
	avgSize := (len(files) + partitions - 1) / partitions // Ceiling division

	for i := range result {
		result[i] = make([]string, 0, avgSize) // Preallocate capacity
	}

	// Distribute files into partitions
	for i, file := range files {
		result[i%partitions] = append(result[i%partitions], file)
	}

	return result
}

func partitionFilesBySize(files []fileInfo, partitions int) [][]fileInfo {
	if partitions <= 0 {
		return nil
	}

	// Sort files by size (largest first)
	sort.Slice(files, func(i, j int) bool {
		return files[i].size > files[j].size
	})

	// Create partition buckets
	result := make([][]fileInfo, partitions)
	sizes := make([]int64, partitions) // Track partition sizes

	for _, file := range files {
		// Find the partition with the smallest current size
		minIndex := 0
		for i := 1; i < partitions; i++ {
			if sizes[i] < sizes[minIndex] {
				minIndex = i
			}
		}

		// Assign file to the partition
		result[minIndex] = append(result[minIndex], file)
		sizes[minIndex] += file.size
	}

	return result
}
