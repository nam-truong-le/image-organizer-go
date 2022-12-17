package image_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nam-truong-le/image-organizer-go/image"
)

func TestDestinationDir(t *testing.T) {
	dir, err := image.DestinationDir("/Volumes/photo/iphone_nhung/2022.12.13.2/IMG_5636.MOV", false)
	assert.NoError(t, err)
	fmt.Println(dir)
}

func TestDestinationDirAnalyze(t *testing.T) {
	dir, err := image.DestinationDir("/Volumes/photo/__date_time_not_found__/2022/12/17/yFUZrsFJzl_IMG_3039.JPG", true)
	assert.Error(t, err)
	assert.Empty(t, dir)
}
