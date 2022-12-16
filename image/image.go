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
	"Date/Time Original", "Create Date",
}

func DestinationDir(file string) (*string, error) {
	logger := log.Logger()
	logger.Infof("Read created date of [%s]", file)
	out, err := exec.Command("exiftool", file).Output()
	if err != nil {
		logger.Errorf("Failed to run mdls: %s\n%s", err, string(out))
		return nil, err
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		before, after, found := strings.Cut(line, ":")
		if !found {
			continue
		}

		key := strings.TrimSpace(before)
		value := strings.TrimSpace(after)
		if !lo.Contains(searchKeys, key) {
			continue
		}

		t, err := time.Parse("2006:01:02 15:04:05", value)
		if err != nil {
			continue
		}
		dir := t.Format("2006/1/2")
		return &dir, nil
	}

	return nil, fmt.Errorf("no content created date found")
}
