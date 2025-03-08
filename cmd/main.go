package main

import (
	"fmt"
	"log"

	"github.com/ezrantn/treecut"
	"github.com/ezrantn/treecut/internal/cli"
)

func main() {
	config, err := cli.ParseCLI()
	if err != nil {
		log.Fatalf("cannot parse cli: %v", err)
	}

	if err := treecut.MakePartitions(config); err != nil {
		log.Fatalf("Error partitioning files: %v", err)
	} else {
		fmt.Println("Success creating symlinks...")
	}
}
