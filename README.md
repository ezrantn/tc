# treecut 🌲

treecut is a Go library and CLI tool for splitting large file trees into smaller, more manageable subtrees using symbolic links. Whether you're organizing massive datasets, optimizing storage, or enabling parallel processing, treecut helps you partition files efficiently—without creating duplicates.

It supports two partitioning methods:

- By file count → Each partition contains approximately the same number of files.
- By file size → Each partition holds a roughly equal total file size.

If you’ve ever struggled with thousands (or millions) of files cluttering a single directory, you know how frustrating it can be. Large directories can slow down file operations, complicate backups, and overwhelm your storage system. treecut provides a simple way to reorganize and distribute files efficiently.

Partitioning a file tree can improve load balancing by spreading files across multiple storage devices, enable parallel processing by allowing batch jobs to run on subsets of files, and make dataset management easier by breaking large datasets into smaller chunks. It’s also useful for organizing files before migration, ensuring backups and transfers are faster and more reliable.

By using symlinks instead of copying files, treecut prevents duplication and saves disk space while maintaining full access to your original files.

> [!NOTE]
> If you're not sure what a symbolic link is, read [this article](https://en.wikipedia.org/wiki/Symbolic_link) on Wikipedia.

## Features

- Partition files by count or size
- Creates batch symlinks instead of duplicating files
- Optimized for fast file traversal
- Simple API for integration into Go applications
- Command-line tool for quick use
- Cross-platform support (Windows, macOS, Linux)

## Installation

### Go Library

```bash
go get github.com/ezrantn/treecut
```

### CLI

```bash
go install github.com/ezrantn/treecut@latest
```

## Usage

Using treecut as a library is simple:

```go
// Here we are not partitioning by size, which means we are using the default value, partitioning by file count.
config := treecut.PartitionConfig{
    SourceDir:  "examples/data",
    OutputDirs: []string{"examples/partition1", "examples/partition2"},
    BySize:     false,
}

if err := treecut.MakePartitions(config); err != nil {
    slog.Error(err.Error())
}
```

To use the CLI tool, run:

```bash
./bin/treecut --source=examples/data --output=examples/partition1,examples/partition2
```

- `--source=<path>`: Specifies the source directory containing the files you want to partition.
- `--output=<paths>`: Defines one or more output directories where symlinks to the partitioned files will be created. Multiple directories should be separated by commas (`,`).

**Important Notes**

- All flags are required—you must specify both --source and --output.
- Multiple output directories allow for distributing files across partitions. The more output directories you provide, the more partitions treecut will create.
- Files will not be copied, only symbolic links will be created in the output directories, saving disk space.

## Why Treecut?

Symbolic links can be created manually using the native `ln` command in Linux/macOS:

```bash
ln -s /path/to/source /path/to/link
```

While `ln` works well for individual files and directories, batch creating symlinks for an entire directory tree while ensuring balanced partitioning is difficult. This is where treecut comes in.

If you want to partition a large directory into multiple subdirectories using `ln`, you would need a custom script like this:

```sh
mkdir -p partition1 partition2
count=0
for file in source/*; do
    if [ $((count % 2)) -eq 0 ]; then
        ln -s "$(realpath "$file")" partition1/
    else
        ln -s "$(realpath "$file")" partition2/
    fi
    count=$((count + 1))
done
```

That looks tedious..

- You must manually decide how to distribute files.
- There's no built-in way to balance by file size.
- The script can get complex when dealing with nested directories.

With treecut, the same task is effortless:

```bash
./bin/treecut --source=examples/data --output=examples/partition1,examples/partition2
```

- Automatically partitions files by count or size.
- Works recursively on nested directories.
- No scripting required—just run a single command.
- Cross-platform support: Unlike `ln`, which is not available on Windows, treecut works on Windows, macOS, and Linux.

Since treecut wraps around the same core idea as `ln` but adds automation, cross-platform support, and partitioning logic, it serves as a better alternative when dealing with large file trees.
