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
	reset     = "\033[0m"
	red       = "\033[31m"
	version   = "v0.0.1"
	asciiText = `
╱╭╮╱╱╱╱╱╱╱╱╱╱╱╱╱╱╭╮
╭╯╰╮╱╱╱╱╱╱╱╱╱╱╱╱╭╯╰╮
╰╮╭╋━┳━━┳━━┳━━┳╮┣╮╭╯
╱┃┃┃╭┫┃━┫┃━┫╭━┫┃┃┃┃
╱┃╰┫┃┃┃━┫┃━┫╰━┫╰╯┃╰╮
╱╰━┻╯╰━━┻━━┻━━┻━━┻━╯`
)

// ParseCLI parses command-line arguments and returns a PartitionConfig.
func ParseCLI() (tc.PartitionConfig, bool, error) {
	if len(os.Args) == 1 {
		printHelp()
		os.Exit(0)
	}

	var versionFlag bool
	flag.BoolVar(&versionFlag, "version", false, "Print tc version")
	flag.BoolVar(&versionFlag, "v", false, "Shorthand for --version")

	sourceDir := flag.String("source", "", "Source directory to partition")
	flag.StringVar(sourceDir, "s", "", "Shorthand for --source")

	outputDirs := flag.String("output", "", "Comma-separated list of output directories")
	flag.StringVar(outputDirs, "o", "", "Shorthand for --output")

	bySize := flag.Bool("by-size", false, "Partition files by size")
	flag.BoolVar(bySize, "b", false, "Shorthand for --by-size")

	unlink := flag.Bool("unlink", false, "Unlink symlinks and remove the partition directories")
	flag.BoolVar(unlink, "u", false, "Shorthand for --unlink")

	flag.Parse()

	if versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

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

// printError prints an error in red color
func PrintError(err error) {
	fmt.Fprintf(os.Stderr, "%sERROR:%s %v\n", red, reset, err)
}

// splitOutputDirs splits output directories from a comma-separated string.
func splitOutputDirs(output string) ([]string, error) {
	if strings.TrimSpace(output) == "" {
		return nil, errors.New("output directories cannot be empty")
	}
	return strings.Split(output, ","), nil
}

func printHelp() {
	fmt.Println(asciiText)
	fmt.Println()
	fmt.Println("tc (treecut) - A Fast and Efficient File Tree Partitioning Tool")
	fmt.Println("\033[3mDeveloped by: Ezra Natanael\033[0m")
	fmt.Println()
	fmt.Println("tc is a Go library and CLI tool for splitting large file trees into smaller, more manageable subtrees using symbolic links.")
	fmt.Println("It helps organize massive datasets, optimize storage, and enable parallel processing by partitioning files efficiently—without duplication.")
	fmt.Println()
	fmt.Println("Partitioning Methods:")
	fmt.Println("  - By file count → Each partition contains approximately the same number of files.")
	fmt.Println("  - By file size  → Each partition holds a roughly equal total file size.")
	fmt.Println()
	fmt.Println("Why Use tc?")
	fmt.Println("  - Prevent large directories from slowing down file operations.")
	fmt.Println("  - Improve load balancing across multiple storage devices.")
	fmt.Println("  - Enable parallel processing by distributing files into smaller subsets.")
	fmt.Println("  - Optimize backups and transfers by organizing files more efficiently.")
	fmt.Println("  - Save disk space by using symlinks instead of copying files.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  tc --source <dir> --output <dir1,dir2,...> [--by-size]")
	fmt.Println("  tc --unlink --output <dir1,dir2,...>")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -s, --source <dir>   Source directory to partition")
	fmt.Println("  -o, --output <dirs>  Comma-separated list of output directories")
	fmt.Println("  -b, --by-size        Partition files by size instead of default method")
	fmt.Println("  -u, --unlink         Remove symlinks and partition directories")
	fmt.Println("  -v, --version        Print tc (treecut) version")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  tc --source /data --output /part1,/part2")
	fmt.Println("  tc -s /data -o /part1,/part2")
	fmt.Println("  tc --unlink --output /part1,/part2")
	fmt.Println("  tc -u -o /part1,/part2")
	fmt.Println()
	fmt.Println("For more details, visit: https://github.com/ezrantn/tc")
}
