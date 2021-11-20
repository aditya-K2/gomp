package main

import (
	"errors"
	_ "image/jpeg"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/shkh/lastfm-go/lastfm"
	"github.com/spf13/viper"
)

func getImageFromLastFM(artist, album, imagePath string) (string, error) {
	api := lastfm.New(viper.GetString("LASTFM_API_KEY"), viper.GetString("LASTFM_API_SECRET"))
	v, err := api.Album.GetInfo(map[string]interface{}{
		"artist":      artist,
		"album":       album,
		"autocorrect": viper.GetInt("LASTFM_AUTO_CORRECT"),
	})
	if err != nil {
		return "", err
	} else {
		return downloadImage(v.Images[len(v.Images)-1].Url, imagePath)
	}
}

func downloadImage(url string, imagePath string) (string, error) {
	var reader io.Reader
	if strings.HasPrefix(url, "http") {
		r, err := http.Get(url)
		if err != nil {
			return "", err
		}
		defer r.Body.Close()
		reader = r.Body
		v, err := io.ReadAll(reader)
		if err == nil {
			b, err := os.Create(imagePath)
			if err == nil {
				b.Write(v)
				return imagePath, nil
			} else {
				b.Close()
				return "", err
			}
		} else {
			return "", err
		}
	}
	return "", errors.New("Image Not Received")
}
