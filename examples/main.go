package main

import (
	"log"

	"github.com/ezrantn/trc"
)

// This is a code example of how to use this library. In the `examples` directory, I created a `data` folder
// to store our example files, in this case `.txt`, but you could store anything.
//
// Here we are not partitioning by size or file, which means we are using the default value,
// partitioning by file type.
func main() {
	// Create partition example:
	//
	config := trc.PartitionConfig{
		SourceDir:  "examples/data",
		OutputDirs: []string{"examples/partition1", "examples/partition2"},
	}

	if err := trc.MakePartitions(config); err != nil {
		log.Fatal(err)
	}

	// Unlink example:
	//
	// outputDirs := []string{"examples/partition1", "examples/partition2"}
	// // outputDirsFalse := []string{"examples/false_dir"}
	// if err := trc.RemovePartitions(outputDirs); err != nil {
	// 	log.Fatal(err)
	// }
}
