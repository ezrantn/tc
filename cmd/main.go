package main

import (
	"fmt"
	"os"

	"github.com/ezrantn/trc"
	"github.com/ezrantn/trc/internal/cli"
)

func main() {
	config, unlink, err := cli.ParseCLI()
	if err != nil {
		cli.PrintError(err)
		os.Exit(1)
	}

	if unlink {
		fmt.Println("Removing partitions and symlinks...")
		if err := trc.RemovePartitions(config.OutputDirs); err != nil {
			fmt.Println("Error removing partitions:", err)
			os.Exit(1)
		}

		fmt.Println("Partitions removed sucessfully")
	} else {
		fmt.Println("Creating partitions...")
		if err := trc.MakePartitions(config); err != nil {
			fmt.Println("Error creating partitions:", err)
			os.Exit(1)
		}

		fmt.Println("Partitions created sucessfully")
	}
}
