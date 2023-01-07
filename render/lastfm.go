package render

//
// # Getting Album Art from [LastFm API](https://www.last.fm/api)
//
// 1. Generate API Key [here](https://www.last.fm/login?next=%2Fapi%2Faccount%2Fcreate%3F_pjax%3D%2523content)
//
//    ![Screenshot from 2021-11-13 21-54-54](https://user-images.githubusercontent.com/51816057/141651276-f76a5c7f-65fe-4a1a-b130-18cdf67dd471.png)
//
// 2. Add the api key and api secret to config.yml
//
// ```yml
// GET_COVER_ART_FROM_LAST_FM : True # Turn On Getting Album art from lastfm api
// LASTFM_API_KEY: "YOUR API KEY HERE"
// LASTFM_API_SECRET: "YOUR API SECRET HERE"
// ```
// 3. Auto correct
//
// ![Screenshot from 2021-11-13 21-59-46](https://user-images.githubusercontent.com/51816057/141651414-1586577a-cab2-48e2-a24b-1053f8634fbe.png)
//
// ```yml
// LASTFM_AUTO_CORRECT: 1  # 0 means it is turned off
// ```

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
