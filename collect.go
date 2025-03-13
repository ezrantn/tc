package trc

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gabriel-vasile/mimetype"
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

func collectFiles(sourceDir string) ([]string, error) {
	filesChan := make(chan string, 100)
	errChan := make(chan error, 1)

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

func collectFilesWithSize(sourceDir string) ([]fileInfo, error) {
	filesChan := make(chan fileInfo, 100)
	errChan := make(chan error, 1)
	var files []fileInfo

	go func() {
		defer close(filesChan)
		err := filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				info, err := d.Info()
				if err != nil {
					return err
				}

				if ok, _ := isValidFileName(filepath.Base(path)); !ok {
					return errors.New("invalid file name: " + path)
				}

				filesChan <- fileInfo{path: path, size: info.Size()}
			}

			return nil
		})
		errChan <- err
		close(errChan)
	}()

	// Read from filesChan
	for file := range filesChan {
		files = append(files, file)
	}

	// Return error if any
	if err := <-errChan; err != nil {
		return nil, err
	}

	return files, nil
}

// collectFilesWithMimeType collects files from the source directory and categorizes them by MIME type
func collectFilesWithMimeType(sourceDir string) (map[string][]string, error) {
	mimeMap := make(map[string][]string)

	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || info.Size() == 0 {
			return nil
		}

		// Detect MIME type using third-party library
		mtype, err := mimetype.DetectFile(path)
		if err != nil {
			fmt.Printf("Failed to detect MIME for %s: %v\n", path, err)
			return nil
		}

		// Extract the category (e.g., "image", "video", etc.)
		mainType := mtype.String()
		category := mainType[:strings.Index(mainType, "/")]
		mimeMap[category] = append(mimeMap[category], path)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return mimeMap, nil
}
