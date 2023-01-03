package render

import (
	_ "image/jpeg"

	"github.com/aditya-K2/gomp/config"
	"github.com/aditya-K2/gomp/utils"
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
		return utils.DownloadImage(v.Images[len(v.Images)-1].Url, imagePath)
	}
}
