package render

import (
	"errors"
	"fmt"
	"image"
	"os"

	"github.com/aditya-K2/gomp/cache"
	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/gomp/ui"
	"github.com/aditya-K2/gomp/ui/notify"
	"github.com/aditya-K2/gomp/utils"
	"github.com/dhowden/tag"
	"github.com/nfnt/resize"
	"github.com/spf13/viper"
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

// GetImagePath This Function returns the path to the image that is to be
// rendered it checks first for the image in the cache
// else it adds the image to the cache and then extracts it and renders it.
func GetImagePath(path string) string {
	a, err := client.Conn.ListInfo(path)
	var extractedImage string = viper.GetString("DEFAULT_IMAGE_PATH")
	if err == nil && len(a) != 0 {
		if cache.Exists(a[0]["artist"], a[0]["album"]) {
			extractedImage = cache.GenerateName(a[0]["artist"], a[0]["album"])
		} else {
			imagePath := cache.GenerateName(a[0]["artist"], a[0]["album"])
			absPath := utils.CheckDirectoryFmt(viper.GetString("MUSIC_DIRECTORY")) + path
			if _eimg, exErr := ExtractImage(absPath, imagePath); exErr != nil {
				if viper.GetString("GET_COVER_ART_FROM_LAST_FM") == "TRUE" {
					downloadedImage, lFmErr := getImageFromLastFM(a[0]["artist"], a[0]["album"], imagePath)
					if lFmErr == nil {
						notify.Send("Image From LastFM")
						extractedImage = downloadedImage
					} else {
						notify.Send(exErr.Error())
					}
				}
			} else {
				notify.Send("Image Extracted Succesfully!")
				extractedImage = _eimg
			}
		}
	} else {
		notify.Send(fmt.Sprintf("Couldn't Get Attributes for %s", path))
	}
	return extractedImage
}

// Gets the Image Struct from the provided path
func GetImg(uri string) (image.Image, error) {
	f, err := os.Open(uri)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	fw, fh := utils.GetFontWidth()
	img = resize.Resize(
		uint(float32(ui.ImgW)*(fw+float32(viper.GetFloat64("IMAGE_WIDTH_EXTRA_X")))),
		uint(float32(ui.ImgH)*(fh+float32(viper.GetFloat64("IMAGE_WIDTH_EXTRA_Y")))),
		img,
		resize.Bilinear,
	)
	return img, nil
}
