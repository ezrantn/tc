package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ezrantn/treecut"
)

func ParseCLI() (treecut.PartitionConfig, bool, error) {
	sourceDir := flag.String("source", "", "Source directory to partition")
	outputDirs := flag.String("output", "", "Comma-separated list of output directories")
	bySize := flag.Bool("by-size", false, "Partition files by size")
	unlink := flag.Bool("unlink", false, "Unlink symlinks and remove the partition directories")

	flag.Parse()

	// If the unlink is set, no need to check sourceDir
	if *unlink {
		if *outputDirs == "" {
			fmt.Println("Usage: treecut --source=<source-dir> --output=<output-dir1,output-dir2> [--by-size]")
			os.Exit(1)
		}

		outputDirsList, err := splitOutputDirs(*outputDirs)
		if err != nil {
			return treecut.PartitionConfig{}, false, errors.New("cannot split output directories")
		}

		return treecut.PartitionConfig{OutputDirs: outputDirsList}, true, nil
	}

	if *sourceDir == "" || *outputDirs == "" {
		fmt.Println("Usage: treecut --source=<source-dir> --output=<output-dir1,output-dir2> [--by-size]")
		os.Exit(1)
	}

	outputDirsList, err := splitOutputDirs(*outputDirs)
	if err != nil {
		return treecut.PartitionConfig{}, false, errors.New("cannot split directory output")
	}

	return treecut.PartitionConfig{
		SourceDir:  *sourceDir,
		OutputDirs: outputDirsList,
		BySize:     *bySize,
	}, false, nil
}

func splitOutputDirs(output string) ([]string, error) {
	if output == "" {
		return nil, errors.New("output directories not found")
	}

	return strings.Split(output, ","), nil
}
