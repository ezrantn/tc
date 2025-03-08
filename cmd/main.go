package main

import (
	"log"

	"github.com/ezrantn/treecut"
	"github.com/ezrantn/treecut/internal/cli"
)

func main() {
	config := cli.ParseCLI()

	if err := treecut.MakePartitions(config); err != nil {
		log.Fatalf("Error partitioning files: %v", err)
	}
}
