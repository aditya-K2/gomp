package render

import (
	"errors"
	"fmt"
	"image"
	"os"

	"github.com/aditya-K2/gomp/cache"
	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/gomp/config"
	"github.com/aditya-K2/gomp/ui"
	"github.com/aditya-K2/utils"
	"github.com/dhowden/tag"
	"github.com/nfnt/resize"
)

var (
	ExtractionErr    = errors.New("No Image Extracted")
	ImageNotFoundErr = errors.New("Image Not Found")
	ImageWriteErr    = errors.New("Error Writing Image to File")
)

func CreateFile(path string, data []byte) error {
	if imHd, err := os.Create(path); err != nil {
		return ImageWriteErr
	} else {
		imHd.Write(data)
		return nil
	}
}

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
		if err := CreateFile(imagePath, picture.Data); err != nil {
			return "", ImageNotFoundErr
		}
	}

	return imagePath, nil
}

/*
GetImagePath This Function returns the path to the image that is to be
rendered it checks first for the image in the cache
else it adds the image to the cache and then extracts it and renders it.
*/
func GetImagePath(path string) (imagePath string) {
	imagePath = config.Config.DefaultImagePath
	attrs, attrErr := client.Conn.ListInfo(path)

	if attrErr != nil || len(attrs) <= 0 {
		ui.SendNotification(fmt.Sprintf("Error getting image from LastFM: Attribute Error(%v)", attrErr))
		return
	}

	artist := attrs[0]["artist"]
	album := attrs[0]["album"]
	cachedPath := cache.GenerateName(artist, album)

	if cache.Exists(artist, album) {
		imagePath = cachedPath
		return
	}

	absPath := utils.CheckDirectoryFmt(config.Config.MusicDirectory) + path
	extractedImage, xErr := ExtractImage(absPath, cachedPath)

	if xErr == nil {
		ui.SendNotification("Covert art metadata extracted from file.")
		return extractedImage
	}

	// Query MPD for Album Art.
	// See: https://mpd.readthedocs.io/en/latest/protocol.html#the-music-database
	albumArt, cErr := client.Conn.AlbumArt(path)
	if cErr == nil {
		if err := CreateFile(cachedPath, albumArt); err == nil {
			ui.SendNotification("Cover art retrieved from MPD")
			return cachedPath
		}
	}

	if !config.Config.GetCoverArtFromLastFm {
		// Display the Extraction Error and fallback to default image
		ui.SendNotification(xErr.Error())
		return
	}

	downloadedImage, lFmErr := getImageFromLastFM(artist, album, cachedPath)

	if lFmErr != nil {
		ui.SendNotification("Error Downloading Image from LastFM: " + xErr.Error())
		return
	}

	ui.SendNotification("Image From LastFM")
	imagePath = downloadedImage

	return
}

/* Gets the Image Struct from the provided path */
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
		uint(float32(ui.ImgW)*(fw+float32(config.Config.ExtraImageWidthX))),
		uint(float32(ui.ImgH)*(fh+float32(config.Config.ExtraImageWidthY))),
		img,
		resize.Bilinear,
	)
	return img, nil
}
