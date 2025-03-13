package trc

import (
	"errors"
	"fmt"
	"os"
	"sort"
)

// PartitionConfig holds the configuration for partitioning files
type PartitionConfig struct {
	SourceDir  string   // Original directory
	OutputDirs []string // Partition directories
	BySize     bool     // Set to true to activate partition by size (largest -> smallest)
	ByFile     bool     // Partition by MIME type
}

// MakePartitions partitions the files in the source directory according to the configuration.
func MakePartitions(config PartitionConfig) error {
	if len(config.OutputDirs) == 0 {
		return errors.New("at least one output directory is required")
	}

	partitionFn, err := getPartitionFunction(config.ByFile, config.BySize)
	if err != nil {
		return err
	}

	return partitionFn(config.SourceDir, config.OutputDirs)
}

// getPartitionFunction returns the appropriate partition function based on the flags.
func getPartitionFunction(byFile, bySize bool) (func(string, []string) error, error) {
	switch {
	case byFile:
		return partitionByFile, nil
	case bySize:
		return partitionBySize, nil
	default:
		return partitionByType, nil
	}
}

// RemovePartitions removes the partition directories and their symlinks.
func RemovePartitions(outputDirs []string) error {
	if err := removeSymlinkTree(outputDirs); err != nil {
		return fmt.Errorf("failed to remove symlink tree: %w", err)
	}

	for _, dir := range outputDirs {
		if err := os.RemoveAll(dir); err != nil {
			return fmt.Errorf("failed to remove partition directory %s: %w", dir, err)
		}
	}

	return nil
}

// partitionFiles splits a list of file paths into equal-sized groups.
func partitionFiles(files []string, partitions int) [][]string {
	if partitions <= 0 || len(files) == 0 {
		return nil
	}

	result := make([][]string, partitions)
	avgSize := (len(files) + partitions - 1) / partitions

	// Preallocate capacity to avoid reallocations
	for i := range result {
		result[i] = make([]string, 0, avgSize)
	}

	// Distribute files across partitions
	for i, file := range files {
		result[i%partitions] = append(result[i%partitions], file)
	}

	return result
}

// partitionFilesBySize splits files into partitions while balancing their sizes.
func partitionFilesBySize(files []fileInfo, partitions int) [][]fileInfo {
	if partitions <= 0 || len(files) == 0 {
		return nil
	}

	// Sort files by size (largest first)
	sort.Slice(files, func(i, j int) bool {
		return files[i].size > files[j].size
	})

	result := make([][]fileInfo, partitions)
	sizes := make([]int64, partitions)

	// Distribute files across partitions to balance the size
	for _, file := range files {
		minIndex := findMinPartitionIndex(sizes)
		result[minIndex] = append(result[minIndex], file)
		sizes[minIndex] += file.size
	}

	return result
}

// findMinPartitionIndex returns the index of the partition with the smallest size.
func findMinPartitionIndex(sizes []int64) int {
	minIndex := 0
	for i := 1; i < len(sizes); i++ {
		if sizes[i] < sizes[minIndex] {
			minIndex = i
		}
	}
	return minIndex
}

// partitionByFile partitions files by their MIME type.
func partitionByFile(sourceDir string, outputDirs []string) error {
	files, err := collectFiles(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to collect files from %s: %w", sourceDir, err)
	}

	partitions := partitionFiles(files, len(outputDirs))
	if err := createSymlinkTree(partitions, outputDirs); err != nil {
		return fmt.Errorf("failed to create symlink tree: %w", err)
	}

	return nil
}

// partitionBySize partitions files based on their size, attempting to balance partition sizes.
func partitionBySize(sourceDir string, outputDirs []string) error {
	files, err := collectFilesWithSize(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to collect files with size from %s: %w", sourceDir, err)
	}

	partitions := partitionFilesBySize(files, len(outputDirs))
	if err := createSymlinkTreeBySize(partitions, outputDirs); err != nil {
		return fmt.Errorf("failed to create symlink tree by size: %w", err)
	}

	return nil
}

// partitionByType partitions files by their MIME type using round-robin distribution.
func partitionByType(sourceDir string, destDirs []string) error {
	mimeMap, err := collectFilesWithMimeType(sourceDir)
	if err != nil {
		return err
	}

	if len(destDirs) == 0 {
		return errors.New("no destination directories provided")
	}

	// Round-robin distribution of files across directories
	i := 0
	for category, files := range mimeMap {
		destDir := destDirs[i%len(destDirs)]
		if err := createSymlinkWithMimeType(map[string][]string{category: files}, destDir); err != nil {
			return err
		}

		i++
	}

	return nil
}
