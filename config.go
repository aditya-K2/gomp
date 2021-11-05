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
		"MUSIC_DIRECTORY":      HOME_DIR + "/Music",
		"PORT":                 "6600",
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
