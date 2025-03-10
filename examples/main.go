package main

import (
	"log"

	"github.com/ezrantn/tc"
)

// This is a code example of how to use this library. In the `examples` directory, I created a `data` folder
// to store our example files, in this case `.txt`, but you could store anything.
//
// Here we are not partitioning by size, which means we are using the default value,
// partitioning by file count.
func main() {
	// Create partition example:
	//
	config := tc.PartitionConfig{
		SourceDir:  "examples/data",
		OutputDirs: []string{"examples/partition1", "examples/partition2"},
		BySize:     false,
	}

	if err := tc.MakePartitions(config); err != nil {
		log.Fatal(err)
	}

	// Unlink example:
	//
	outputDirs := []string{"examples/partition1", "examples/partition2"}
	// outputDirsFalse := []string{"examples/false_dir"}
	if err := tc.RemovePartitions(outputDirs); err != nil {
		log.Fatal(err)
	}
}
