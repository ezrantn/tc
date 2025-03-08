package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/ezrantn/treecut"
)

func ParseCLI() treecut.PartitionConfig {
	sourceDir := flag.String("source", "", "Source directory to partition")
	outputDirs := flag.String("output", "", "Comma-separated list of output directories")
	bySize := flag.Bool("by-size", false, "Partition files by size")

	flag.Parse()

	if *sourceDir == "" || *outputDirs == "" {
		fmt.Println("Usage: treecut-cli --source=<source-dir> --output=<output-dir1,output-dir2> [--by-size]")
		os.Exit(1)
	}

	outputDirsList := splitOutputDirs(*outputDirs)

	return treecut.PartitionConfig{
		SourceDir:  *sourceDir,
		OutputDirs: outputDirsList,
		BySize:     *bySize,
	}
}

func splitOutputDirs(output string) []string {
	var dirs []string

	dirs = append(dirs, flag.Args()...)
	return dirs
}
