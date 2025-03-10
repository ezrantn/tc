package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ezrantn/tc"
)

var (
	reset = "\033[0m"
	red   = "\033[31m"
)

// ParseCLI parses command-line arguments and returns a PartitionConfig.
func ParseCLI() (tc.PartitionConfig, bool, error) {
	sourceDir := flag.String("source", "", "Source directory to partition")
	outputDirs := flag.String("output", "", "Comma-separated list of output directories")
	bySize := flag.Bool("by-size", false, "Partition files by size")
	unlink := flag.Bool("unlink", false, "Unlink symlinks and remove the partition directories")

	flag.Parse()

	// Unlink mode (removing partitions)
	if *unlink {
		if *outputDirs == "" {
			return tc.PartitionConfig{}, false, errors.New("missing required --output flag for unlink mode")
		}

		outputDirsList, err := splitOutputDirs(*outputDirs)
		if err != nil {
			return tc.PartitionConfig{}, false, fmt.Errorf("invalid output directories: %w", err)
		}

		return tc.PartitionConfig{OutputDirs: outputDirsList}, true, nil
	}

	// Regular partitioning mode
	if *sourceDir == "" {
		return tc.PartitionConfig{}, false, errors.New("missing required --source flag")
	}

	if *outputDirs == "" {
		return tc.PartitionConfig{}, false, errors.New("missing required --output flag")
	}

	outputDirsList, err := splitOutputDirs(*outputDirs)
	if err != nil {
		return tc.PartitionConfig{}, false, fmt.Errorf("invalid output directories: %w", err)
	}

	return tc.PartitionConfig{
		SourceDir:  *sourceDir,
		OutputDirs: outputDirsList,
		BySize:     *bySize,
	}, false, nil
}

// splitOutputDirs splits output directories from a comma-separated string.
func splitOutputDirs(output string) ([]string, error) {
	if strings.TrimSpace(output) == "" {
		return nil, errors.New("output directories cannot be empty")
	}
	return strings.Split(output, ","), nil
}

// printError prints an error in red color
func PrintError(err error) {
	fmt.Fprintf(os.Stderr, "%sERROR:%s %v\n", red, reset, err)
}
