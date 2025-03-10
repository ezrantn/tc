package tc

import (
	"errors"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
)

type fileInfo struct {
	path string
	size int64
}

// Reserved characters and words for filename validation
// See: https://en.wikipedia.org/wiki/Filename#Reserved_characters_and_words
const (
	characterFilter = `[\x00-\x1F\\/:*?"<>|@!]` // Some FAT systems don't allow @ and ! in filenames
	defaultName     = "file"
	maxLength       = 255
)

var (
	characterFilterRegex = regexp.MustCompile(characterFilter)

	dosReservedNames = map[string]struct{}{
		"CON":     {},
		"PRN":     {},
		"AUX":     {},
		"NUL":     {},
		"CLOCK$":  {},
		"CONFIG$": {},
		"SCREEN$": {},
		"$IDLE$":  {},
		"COM0":    {},
		"COM1":    {},
		"COM2":    {},
		"COM3":    {},
		"COM4":    {},
		"COM5":    {},
		"COM6":    {},
		"COM7":    {},
		"COM8":    {},
		"COM9":    {},
		"LPT0":    {},
		"LPT1":    {},
		"LPT2":    {},
		"LPT3":    {},
		"LPT4":    {},
		"LPT5":    {},
		"LPT6":    {},
		"LPT7":    {},
		"LPT8":    {},
		"LPT9":    {},
	}

	windowsReservedNames = map[string]struct{}{
		"$MFT":     {},
		"$MFTMIRR": {},
		"$LOGFILE": {},
		"$VOLUME":  {},
		"$ATTRDEF": {},
		"$BITMAP":  {},
		"$BOOT":    {},
		"$BADCLUS": {},
		"$SECURE":  {},
		"$UPCASE":  {},
		"$EXTEND":  {},
		"$QUOTA":   {},
		"$OBJID":   {},
		"$REPARSE": {},
	}
)

// isValidFileName ensures a file is not using reserved names and meets system constraints
func isValidFileName(filename string) (bool, error) {
	original := filename
	filename = strings.TrimSpace(filename)
	filename = strings.ToUpper(filename)

	if filename == "" {
		return false, errors.New("filename cannot be empty")
	}

	if len(filename) > maxLength {
		return false, errors.New("filename exceeds maximum length")
	}

	if characterFilterRegex.MatchString(filename) {
		return false, errors.New("filename contains invalid characters")
	}

	baseName := filenameWithoutExtension(filename)

	if _, exists := dosReservedNames[baseName]; exists {
		return false, errors.New("filename is a reserved DOS name")
	}

	if _, exists := windowsReservedNames[baseName]; exists {
		return false, errors.New("filename is a reserved windows name")
	}

	if hasTrailingDotOrSpace(original) {
		return false, errors.New("filename has trailing dots or spaces")
	}

	return true, nil
}

func filenameWithoutExtension(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func hasTrailingDotOrSpace(filename string) bool {
	return regexp.MustCompile(`[.\s]+$`).MatchString(filename)
}

// collectFiles walks through the source directory and returns all file paths
func collectFiles(sourceDir string) ([]string, error) {
	filesChan := make(chan string, 100) // Buffered channel to reduce contention
	errChan := make(chan error, 1)
	// Walk the directory in a separate goroutine
	go func() {
		defer close(filesChan)
		errChan <- filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				baseName := filepath.Base(path)
				if ok, _ := isValidFileName(baseName); !ok {
					return errors.New("invalid file name: " + baseName)
				}

				filesChan <- path
			}
			return nil
		})

		close(errChan)
	}()

	// Collect results
	var files []string
	for file := range filesChan {
		files = append(files, file)
	}

	// Return any errors encountered
	if err, ok := <-errChan; ok && err != nil {
		return nil, err
	}

	return files, nil
}

// collectFilesWithSize collects file paths and sizes
func collectFilesWithSize(sourceDir string) ([]fileInfo, error) {
	var files []fileInfo
	err := filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}

			// Validate filename before adding it
			if ok, _ := isValidFileName(filepath.Base(path)); !ok {
				return errors.New("invalid file name: " + path)
			}

			files = append(files, fileInfo{path: path, size: info.Size()})
		}

		return nil
	})

	return files, err
}
