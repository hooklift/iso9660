package isofs

import (
	"fmt"
	"os"
	"testing"

	"github.com/hooklift/assert"
)

func TestNewReader(t *testing.T) {
	image, err := os.Open("./fixtures/test.iso")
	defer image.Close()
	reader, err := NewReader(image)
	assert.Ok(t, err)
	fmt.Printf("%+v", reader.primaryVolume)
}
