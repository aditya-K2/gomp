package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var (
	HOME_DIR, _ = os.UserHomeDir()
	defaults    = map[string]interface{}{
		"ADDITIONAL_PADDING_X": 12,
		"ADDITIONAL_PADDING_Y": 16,
		"IMAGE_WIDTH_EXTRA_X":  -1.5,
		"IMAGE_WIDTH_EXTRA_Y":  -3.75,
		"MUSIC_DIRECTORY":      getMusicDirectory(),
		"PORT":                 "6600",
		"DEFAULT_IMAGE_PATH":   "default.jpg",
		"COVER_IMAGE_PATH":     "cover.jpg",
	}
)

func readConfig() {
	for k, v := range defaults {
		viper.SetDefault(k, v)
	}
	viper.SetConfigName("config")
	viper.AddConfigPath(HOME_DIR + "/.config/goMP")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Could Not Read Config file.")
	}
}
