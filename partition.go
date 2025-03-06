package treecut

import (
	"errors"
	"sort"
)

type PartitionConfig struct {
	SourceDir  string   // Original directory
	OutputDirs []string // Partition directories
	BySize     bool     // Set to true to activate partition by size (largest -> smallest)
}

// partitionFiles splits a list of file paths into `partitions` equal groups
func partitionFiles(files []string, partitions int) [][]string {
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

func MakePartitions(config PartitionConfig) error {
	if len(config.OutputDirs) == 0 {
		return errors.New("at least one output directory is required")
	}

	// Collect files
	if config.BySize {
		files, err := collectFilesWithSize(config.SourceDir)
		if err != nil {
			return err
		}

		partitions := partitionFilesBySize(files, len(config.OutputDirs))
		return createSymlinkTreeSize(partitions, config.OutputDirs)
	} else {
		files, err := collectFiles(config.SourceDir)
		if err != nil {
			return err
		}

		partitions := partitionFiles(files, len(config.OutputDirs))
		return createSymlinkTree(partitions, config.OutputDirs)
	}
}

func partitionFilesBySize(files []fileInfo, partitions int) [][]fileInfo {
	// Sort files by size (largest first)
	sort.Slice(files, func(i, j int) bool {
		return files[i].size > files[j].size
	})

	// Create partitions buckets
	result := make([][]fileInfo, partitions)
	sizes := make([]int64, partitions) // Track partitions sizes

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
