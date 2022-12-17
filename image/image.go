package image

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/samber/lo"

	"github.com/nam-truong-le/image-organizer-go/log"
)

var searchKeys = []string{
	"Date/Time Original",
	"Create Date",
	"Date Created",
	"Modify Date",
}

var ignoredKey = []string{
	"File Modification Date/Time",
	"File Access Date/Time",
	"File Inode Change Date/Time",
	"Profile Date Time",
	//"Modify Date", // can we add this to searchKeys?
}

var invalidValues = []string{
	"0000:00:00",
}

func DestinationDir(file string, analyzeMode bool) (string, error) {
	logger := log.Logger()
	logger.Infof("Read created date of [%s] [analze=%t]", file, analyzeMode)
	out, err := exec.Command("exiftool", file).Output()
	if err != nil {
		logger.Errorf("Failed to run exiftool: %s\n%s", err, string(out))
		return "", err
	}
	items := getExif(string(out))
	exifCreatedDate, foundCreatedDate := getCreatedDate(items)

	if analyzeMode {
		if foundCreatedDate {
			logger.Warnf("Found registered key: %+v", exifCreatedDate)
			return "", nil
		}

		// analyze
		for _, item := range items {
			invalid := false
			for _, invalidValue := range invalidValues {
				if strings.Contains(item.Value, invalidValue) {
					invalid = true
					break
				}
			}
			if invalid {
				continue
			}
			if lo.Contains(ignoredKey, item.Key) {
				continue
			}

			if strings.Contains(strings.ToLower(item.Key), "date") {
				logger.Errorf("%+v", item)
				logger.Errorf(string(out))
				logger.Fatalf("Found unregistered key has date: %+v", item)
			}
		}

		return "", fmt.Errorf("no content created date found")
	}

	// normal mode
	if !foundCreatedDate {
		return "", fmt.Errorf("no content created date found")
	}
	t, err := time.Parse("2006:01:02 15:04:05", exifCreatedDate.Value)
	if err != nil {
		logger.Fatalf("Invalid date format: %s", exifCreatedDate.Value)
	}
	dir := t.Format("2006/1/2")
	return dir, nil
}

type exif struct {
	Key   string
	Value string
}

func getExif(exifOutput string) []exif {
	items := make([]exif, 0)
	lines := strings.Split(exifOutput, "\n")
	for _, line := range lines {
		before, after, found := strings.Cut(line, ":")
		if !found {
			continue
		}

		key := strings.TrimSpace(before)
		value := strings.TrimSpace(after)

		items = append(items, exif{key, value})
	}
	return items
}

func getCreatedDate(items []exif) (*exif, bool) {
	for _, searchKey := range searchKeys {
		for _, item := range items {
			if searchKey == item.Key {
				invalid := false
				for _, invalidValue := range invalidValues {
					if strings.Contains(item.Value, invalidValue) {
						invalid = true
						break
					}
				}
				if invalid {
					continue
				}

				return &item, true
			}
		}
	}
	return nil, false
}
