package render

import (
	"github.com/aditya-K2/gomp/utils"
	"os"
	"strings"

	"github.com/bogem/id3v2"
	"github.com/mewkiz/flac"
	"github.com/mewkiz/flac/meta"
	"github.com/spf13/viper"
)

func GetMp3Image(songPath, imagePath string) string {
	tag, err := id3v2.Open(songPath, id3v2.Options{Parse: true})
	if err != nil {
		return ""
	}
	defer tag.Close()
	if err != nil {
		return ""
	}
	// Read tags.
	Frames := tag.GetFrames(tag.CommonID("Attached picture"))
	var ImageData []byte
	for _, er := range Frames {
		pic, ok := er.(id3v2.PictureFrame)
		if ok {
			for _, i := range pic.Picture {
				ImageData = append(ImageData, byte(i))
			}
			imageHandler, err := os.Create(imagePath)
			if err != nil {
				return ""
			} else {
				imageHandler.Write(ImageData)
				return imagePath
			}
		}
	}
	return ""
}

func GetFlacImage(songPath, imagePath string) string {
	stream, err := flac.ParseFile(songPath)
	if err != nil {
		return ""
	}
	defer stream.Close()
	for _, block := range stream.Blocks {
		if block.Type == meta.TypePicture {
			pic := block.Body.(*meta.Picture)
			if pic.Type == 3 {
				imageHandler, err := os.Create(imagePath)
				if err != nil {
					return ""
				}
				imageHandler.Write(pic.Data)
				return imagePath
			}
		}
	}
	return ""
}

func ExtractImageFromFile(uri string, imagePath string) string {
	_i := imagePath
	if strings.HasSuffix(uri, ".mp3") {
		imagePath := GetMp3Image(uri, imagePath)
		if imagePath == "" {
			utils.Copy(viper.GetString("DEFAULT_IMAGE_PATH"), _i)
			return viper.GetString("DEFAULT_IMAGE_PATH")
		}
	} else if strings.HasSuffix(uri, ".flac") {
		imagePath := GetFlacImage(uri, imagePath)
		if imagePath == "" {
			utils.Copy(viper.GetString("DEFAULT_IMAGE_PATH"), _i)
			return viper.GetString("DEFAULT_IMAGE_PATH")
		}
	} else {
		utils.Copy(viper.GetString("DEFAULT_IMAGE_PATH"), _i)
		return viper.GetString("DEFAULT_IMAGE_PATH")
	}
	return imagePath
}
