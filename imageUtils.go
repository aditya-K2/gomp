package main

import (
	"image"
	"os"

	"github.com/dhowden/tag"
	"github.com/nfnt/resize"
)

/*
	Gets the Image Path from the uri to the string passed
	if embedded image is found the path to that Image is returned else
	path to default image is sent.
*/
func getAlbumArt(uri string) string {
	var path string = "default.jpg"
	f, err := os.Open(uri)
	if err != nil {
		panic(err)
	}
	m, err := tag.ReadFrom(f)
	if err != nil {
		panic(err)
	}
	albumCover := m.Picture()
	if albumCover != nil {
		b, err := os.Create("hello.jpg")
		if err != nil {
			panic(err)
		}
		defer b.Close()
		b.Write(albumCover.Data)
		path = "hello.jpg"
		b.Close()
	}
	f.Close()
	return path
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

	img = resize.Thumbnail(
		uint(IMG_W*22), uint(IMG_H*15),
		img,
		resize.Bilinear,
	)

	return img, nil
}
