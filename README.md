# treecut

treecut is a Go library and CLI tool for splitting large file trees into smaller subtrees using symlinks. It efficiently partitions files either by count or size, making it useful for organizing large datasets, parallel processing, and storage management.

It supports two partitioning methods:

- By file count → Each partition has an approximately equal number of files.
- By file size → Each partition has an approximately equal total file size.

## Motivation

Partitioning a large file tree can be useful for:

- Load balancing: Distributing files across multiple storage devices.
- Parallel processing: Running batch jobs on subsets of files.
- Dataset management: Splitting datasets for easier handling.
- Backup & transfer: Organizing files before migration.

By using symlinks, treecut ensures that no file duplication occurs, saving disk space.

## Features

- Partition files by count or size
- Creates symlinks instead of duplicating files
- Optimized file traversal using WalkDir & goroutines
- Simple API for integration into Go applications
- Command-line tool for quick use

## Installation

### Go Library

```bash
go get github.com/ezrantn/treecut
```
