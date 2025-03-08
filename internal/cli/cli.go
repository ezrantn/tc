package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ezrantn/treecut"
)

func ParseCLI() (treecut.PartitionConfig, error) {
	sourceDir := flag.String("source", "", "Source directory to partition")
	outputDirs := flag.String("output", "", "Comma-separated list of output directories")
	bySize := flag.Bool("by-size", false, "Partition files by size")

	flag.Parse()

	if *sourceDir == "" || *outputDirs == "" {
		fmt.Println("Usage: treecut --source=<source-dir> --output=<output-dir1,output-dir2> [--by-size]")
		os.Exit(1)
	}

	outputDirsList, err := splitOutputDirs(*outputDirs)
	if err != nil {
		return treecut.PartitionConfig{}, errors.New("cannot split directory output")
	}

	return treecut.PartitionConfig{
		SourceDir:  *sourceDir,
		OutputDirs: outputDirsList,
		BySize:     *bySize,
	}, nil
}

func splitOutputDirs(output string) ([]string, error) {
	if output == "" {
		return nil, errors.New("output directories not found")
	}

	return strings.Split(output, ","), nil
}
