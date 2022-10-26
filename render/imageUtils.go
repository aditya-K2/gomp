package render

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/dhowden/tag"
)

var (
	ExtractionErr    = errors.New("Empty Image Extracted")
	ImageNotFoundErr = errors.New("Image Not Found")
	ImageWriteErr    = errors.New("Error Writing Image to File")
)

func ExtractImage(path string, imagePath string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", ExtractionErr
	}
	defer f.Close()

	meta, err := tag.ReadFrom(f)
	if err != nil {
		return "", ExtractionErr
	}

	if picture := meta.Picture(); picture == nil {
		return "", ImageNotFoundErr
	} else {
		if err := ioutil.WriteFile(imagePath, picture.Data, fs.ModeIrregular); err != nil {
			return "", ImageWriteErr
		}
	}

	return imagePath, nil
}
