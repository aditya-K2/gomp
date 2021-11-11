package main

import (
	"image"
	"os"
	"strings"

	"github.com/bogem/id3v2"
	"github.com/mewkiz/flac"
	"github.com/mewkiz/flac/meta"
	"github.com/nfnt/resize"
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

func extractImageFromFile(uri string) string {
	if strings.HasSuffix(uri, ".mp3") {
		imagePath := GetMp3Image(uri, viper.GetString("COVER_IMAGE_PATH"))
		if imagePath == "" {
			return viper.GetString("DEFAULT_IMAGE_PATH")
		} else {
			return imagePath
		}
	} else if strings.HasSuffix(uri, ".flac") {
		imagePath := GetFlacImage(uri, viper.GetString("COVER_IMAGE_PATH"))
		if imagePath == "" {
			return viper.GetString("DEFAULT_IMAGE_PATH")
		} else {
			return imagePath
		}
	} else {
		return viper.GetString("DEFAULT_IMAGE_PATH")
	}
}

func getImg(uri string) (image.Image, error) {
	f, err := os.Open(uri)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	fw, fh := getFontWidth()
	img = resize.Resize(
		uint(float32(IMG_W)*(fw+float32(viper.GetFloat64("IMAGE_WIDTH_EXTRA_X")))), uint(float32(IMG_H)*(fh+float32(viper.GetFloat64("IMAGE_WIDTH_EXTRA_Y")))),
		img,
		resize.Bilinear,
	)
	return img, nil
}
