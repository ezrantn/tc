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
	ByFile     bool     // Partition by MIME type
}

func MakePartitions(config PartitionConfig) error {
	if len(config.OutputDirs) == 0 {
		return errors.New("at least one output directory is required")
	}

	var err error

	if config.ByFile {
		err = partitionByFile(config.SourceDir, config.OutputDirs)
	} else if config.BySize {
		err = partitionBySize(config.SourceDir, config.OutputDirs)
	} else {
		err = partitionByType(config.SourceDir, config.OutputDirs)
	}

	if err != nil {
		return fmt.Errorf("partitioning failed: %w", err)
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

func partitionByType(sourceDir string, destDirs []string) error {
	mimeMap, err := collectFilesWithMimeType(sourceDir)

	if err != nil {
		return err
	}

	if len(destDirs) == 0 {
		return errors.New("no destination directories provided")
	}

	// Round-robin distribution of files across multiple partitions
	i := 0
	for file := range mimeMap {
		destDir := destDirs[i%len(destDirs)]
		err := createSymlinkWithMimeType(map[string][]string{file: mimeMap[file]}, destDir)
		if err != nil {
			return err
		}

		i++
	}

	return nil
}
