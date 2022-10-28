package render

import (
	"errors"
	"os"

	"github.com/dhowden/tag"
)

var (
	ExtractionErr    = errors.New("No Image Extracted")
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
		if imHd, err := os.Create(imagePath); err != nil {
			return "", ImageWriteErr
		} else {
			imHd.Write(picture.Data)
		}
	}

	return imagePath, nil
}
