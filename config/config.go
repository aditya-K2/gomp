package config

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/aditya-K2/gomp/utils"
	"github.com/spf13/viper"
)

var (
	HOME_DIR, _       = os.UserHomeDir()
	USER_CACHE_DIR, _ = os.UserCacheDir()
	defaults          = map[string]interface{}{
		"ADDITIONAL_PADDING_X": 12,
		"ADDITIONAL_PADDING_Y": 16,
		"IMAGE_WIDTH_EXTRA_X":  -1.5,
		"IMAGE_WIDTH_EXTRA_Y":  -3.75,
		"MUSIC_DIRECTORY":      utils.CheckDirectoryFmt(getMusicDirectory()),
		"PORT":                 "6600",
		"DEFAULT_IMAGE_PATH":   "default.jpg",
		"CACHE_DIR":            utils.CheckDirectoryFmt(USER_CACHE_DIR),
	}
)

func ReadConfig() {
	for k, v := range defaults {
		viper.SetDefault(k, v)
	}
	viper.SetConfigName("config")
	viper.AddConfigPath(HOME_DIR + "/.config/gomp")
	err := viper.ReadInConfig()
	if err != nil {
		utils.Print("RED", "Could Not Read Config file.\n")
		utils.Print("GREEN", "Make a config file $HOME/.config/gomp/config.yml\n")
		os.Exit(1)
	}
}

func GenerateKeyMap(funcMap map[string]func()) {
	for k := range funcMap {
		mappingsForFunctionK := viper.GetStringSlice(k)
		if len(mappingsForFunctionK) != 0 {
			for _, i := range mappingsForFunctionK {
				aV, err := GetAsciiValue(i)
				if err == nil {
					KEY_MAP[aV] = k
				}
			}
		}
	}
}

func getMusicDirectory() string {
	content, err := ioutil.ReadFile(HOME_DIR + "/.config/mpd/mpd.conf")
	if err != nil {
		utils.Print("RED", "There is No Config File for MPD at $HOME/.config/mpd/mpd.conf\n")
		utils.Print("CYAN", "Make Sure there is mpd.conf file or Mention the ")
		utils.Print("GREEN", "music_directory")
		utils.Print("CYAN", " in the config file at $HOME/.config/gomp/config.yml\n")
		os.Exit(1)
	}
	ab := string(content)
	maps := strings.Split(ab, "\n")
	for _, j := range maps {
		if strings.Contains(j, "music_directory") {
			s := strings.SplitAfter(strings.ReplaceAll(j, " ", ""), "y")[1]
			s = strings.ReplaceAll(s, "\t", "")
			d := ""
			for z, m := range s {
				if (z != 0) && (z != (len(s) - 1)) {
					d += string(m)
				}
			}
			return d
		}
	}
	return ""
}
