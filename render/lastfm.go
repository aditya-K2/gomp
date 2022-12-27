package render

import (
	"errors"
	_ "image/jpeg"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/aditya-K2/gomp/config"
	"github.com/shkh/lastfm-go/lastfm"
)

func getImageFromLastFM(artist, album, imagePath string) (string, error) {
	api := lastfm.New(config.Config.LastFmAPIKey, config.Config.LastFmAPISecret)
	v, err := api.Album.GetInfo(map[string]interface{}{
		"artist":      artist,
		"album":       album,
		"autocorrect": config.Config.LastFmAPIAutoCorrect,
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
				if _, err := b.Write(v); err == nil {
					return imagePath, nil
				} else {
					return "", errors.New("could Not Write Image")
				}
			} else {
				b.Close()
				return "", err
			}
		} else {
			return "", err
		}
	}
	return "", errors.New("image Not Received")
}
