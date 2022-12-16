package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/samber/lo"

	"github.com/nam-truong-le/image-organizer-go/image"
	"github.com/nam-truong-le/image-organizer-go/log"
)

var (
	currentTotal = 0
	processed    = 0
	moved        = 0
	moveFailed   = 0
)

func main() {
	logger := log.Logger()

	source := flag.String("source", "", "Input directory")
	destination := flag.String("destination", "", "Output directory")

	flag.Parse()

	if *source == "" || *destination == "" {
		logger.Errorf("Missing parameter(s)")
		panic("missing parameters(s)")
	}

	moveImages(*source, *destination)
}

func moveImages(source string, destination string) {
	logger := log.Logger()
	stat()

	logger.Infof("Move images from [%s] to [%s]", source, destination)
	entries, err := os.ReadDir(source)
	if err != nil {
		logger.Errorf("Failed to read dir [%s]: %s", source, err)
		panic("failed to read source directory")
	}

	currentTotal = currentTotal + len(entries)
	for _, entry := range entries {
		moveDirEntry(source, entry, destination)
		processed++
		stat()
	}
}

func moveDirEntry(source string, entry os.DirEntry, destination string) {
	logger := log.Logger()
	logger.Infof("Move dir entry [%s]", entry.Name())
	entryPath := fmt.Sprintf("%s/%s", source, entry.Name())
	if entry.IsDir() {
		moveImages(entryPath, destination)
		return
	}

	moveFile(entryPath, destination, entry.Name())
}

func moveFile(fileFullPath string, destination string, fileName string) {
	logger := log.Logger()
	logger.Infof("Move file [%s]", fileFullPath)

	destinationDir, err := image.DestinationDir(fileFullPath)
	if err != nil {
		logger.Warnf("Failed to calculate destination dir for [%s]: %s", fileName, err)
		return
	}
	fullDestinationDir := fmt.Sprintf("%s/%s", destination, *destinationDir)
	err = os.MkdirAll(fullDestinationDir, 0777)
	if err != nil {
		logger.Errorf("Failed to create dir [%s]: %s", fullDestinationDir, err)
		panic("failed to create dir")
	}

	destinationName := fmt.Sprintf("%s_%s", lo.RandomString(10, lo.LettersCharset), fileName)
	destinationFile := fmt.Sprintf("%s/%s", fullDestinationDir, destinationName)
	logger.Infof("File will be moved to: %s", destinationFile)

	_, err = os.Stat(destinationFile)
	if err == nil {
		logger.Errorf("File [%s] exists", destinationFile)
		panic("file exists")
	}

	if !errors.Is(err, os.ErrNotExist) {
		logger.Errorf("Unexpected error: %s", err)
		panic("unexpected error")
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
	logger.Infof("Current total = %d | Processed = %d | Moved = %d | Move failed = %d", currentTotal, processed, moved, moveFailed)
}
