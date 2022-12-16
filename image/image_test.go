package image_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nam-truong-le/image-organizer-go/image"
)

func TestDestinationDir(t *testing.T) {
	dir, err := image.DestinationDir("/Volumes/photo/iphone_nhung/2022.12.13.2/IMG_5636.MOV")
	assert.NoError(t, err)
	if dir != nil {
		fmt.Println(*dir)
	}
}
