package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/samber/lo"

	"github.com/nam-truong-le/image-organizer-go/image"
	"github.com/nam-truong-le/image-organizer-go/log"
)

var (
	currentTotal = 0
	processed    = 0
	moved        = 0
	dateNotFound = 0
	moveFailed   = 0
)

var (
	fSource      *string
	fDestination *string
	fAnalyzeDate *bool
)

const (
	notFoundDir = "__date_time_not_found__"
)

func main() {
	logger := log.Logger()

	fSource = flag.String("source", "", "Input directory")
	fDestination = flag.String("destination", "", "Output directory")
	fAnalyzeDate = flag.Bool("analyzeDate", false, "This enables date analyze mode")

	flag.Parse()

	if *fSource == "" {
		logger.Fatalf("Missing parameter(s)")
	}

	if !*fAnalyzeDate && *fDestination == "" {
		logger.Fatalf("Missing parameter(s)")
	}

	fmt.Printf("Source = [%s] | Destination = [%s] | Analyze = [%t] ... continue? (y/n)\n", *fSource, *fDestination, *fAnalyzeDate)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		logger.Fatalf("%s", err)
	}
	if strings.ToLower(strings.TrimSpace(scanner.Text())) != "y" {
		logger.Fatalf("STOPPED")
	}

	processDir(*fSource, *fDestination)

	if *fAnalyzeDate && len(falseNegativeFiles) > 0 {
		logger.Warnf("FALSE NEGATIVE!!!")
		lo.ForEach(falseNegativeFiles, func(f falseNegativeFile, i int) {
			logger.Warnf("%+v", f)
		})

		fmt.Printf("Move these files to (empty value will skip): \n")
		scanner = bufio.NewScanner(os.Stdin)
		scanner.Scan()
		err := scanner.Err()
		if err != nil {
			logger.Fatalf("%s", err)
		}
		dest := strings.TrimSpace(scanner.Text())
		if dest == "" {
			logger.Fatalf("Empty destination")
		}

		for _, f := range falseNegativeFiles {
			err := os.Rename(f.FullPath, fmt.Sprintf("%s/%s", dest, f.Name))
			if err != nil {
				logger.Fatalf("%s", err)
			}
		}
	}
}

func processDir(source string, destination string) {
	logger := log.Logger()
	stat()

	logger.Infof("Source [%s] | Destination [%s]", source, destination)
	entries, err := os.ReadDir(source)
	if err != nil {
		logger.Fatalf("Failed to read dir [%s]: %s", source, err)
	}

	currentTotal = currentTotal + len(entries)
	for _, entry := range entries {
		processDirEntry(source, entry, destination)
		processed++
		stat()
	}
}

func processDirEntry(source string, entry os.DirEntry, destination string) {
	logger := log.Logger()
	logger.Infof("Dir entry [%s]", entry.Name())
	entryPath := fmt.Sprintf("%s/%s", source, entry.Name())
	if entry.IsDir() {
		processDir(entryPath, destination)
		return
	}

	processFile(entryPath, destination, entry.Name())
}

var ignoredFiles = []string{".DS_Store"}

type falseNegativeFile struct {
	FullPath string
	Name     string
}

var falseNegativeFiles = make([]falseNegativeFile, 0)

func processFile(fileFullPath string, destination string, fileName string) {
	logger := log.Logger()
	logger.Infof("Process file [%s]", fileFullPath)

	if lo.Contains(ignoredFiles, fileName) {
		logger.Warnf("Ignored file: %s", fileName)
		return
	}

	// Analyze mode
	if *fAnalyzeDate {
		_, err := image.DestinationDir(fileFullPath, true)
		if err == nil {
			falseNegativeFiles = append(falseNegativeFiles, falseNegativeFile{
				FullPath: fileFullPath,
				Name:     fileName,
			})
			logger.Warnf("False negative files: %d", len(falseNegativeFiles))
		}
		return
	}

	// Normal mode
	destinationDir, err := image.DestinationDir(fileFullPath, false)
	if err != nil {
		logger.Warnf("Failed to calculate destination dir for [%s]: %s", fileName, err)
		destinationDir = fmt.Sprintf("%s/%s", notFoundDir, time.Now().Format("2006/1/2"))
		dateNotFound++
	}
	fullDestinationDir := fmt.Sprintf("%s/%s", destination, destinationDir)
	err = os.MkdirAll(fullDestinationDir, 0777)
	if err != nil {
		logger.Fatalf("Failed to create dir [%s]: %s", fullDestinationDir, err)
	}

	destinationName := fmt.Sprintf("%s_%s", lo.RandomString(10, lo.LettersCharset), fileName)
	destinationFile := fmt.Sprintf("%s/%s", fullDestinationDir, destinationName)
	logger.Infof("File will be moved to: %s", destinationFile)

	_, err = os.Stat(destinationFile)
	if err == nil {
		logger.Fatalf("File [%s] exists", destinationFile)
	}

	if !errors.Is(err, os.ErrNotExist) {
		logger.Fatalf("Unexpected error: %s", err)
	}

	err = os.Rename(fileFullPath, destinationFile)
	if err != nil {
		logger.Warnf("Failed to move file [%s] to [%s]: %s", fileFullPath, destinationFile, err)
		moveFailed++
		stat()
		return
	}
	logger.Infof("File [%s] moved to [%s]", fileFullPath, destinationFile)
	moved++
	stat()
}

func stat() {
	logger := log.Logger()
	logger.Infof("Current total = %d | Processed = %d | Moved = %d | Move failed = %d | Date not found = %d | False negative = %d",
		currentTotal, processed, moved, moveFailed, dateNotFound, len(falseNegativeFiles))
}
