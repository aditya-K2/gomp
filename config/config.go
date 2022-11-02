package config

import (
	"os"

	"github.com/aditya-K2/gomp/config/conf"
	"github.com/aditya-K2/gomp/utils"
	"github.com/spf13/viper"
)

var (
	CONFIG_DIR, CONFIG_ERR    = os.UserConfigDir()
	USER_CACHE_DIR, CACHE_ERR = os.UserCacheDir()
	defaults                  = map[string]interface{}{
		"ADDITIONAL_PADDING_X": 12,
		"ADDITIONAL_PADDING_Y": 16,
		"IMAGE_WIDTH_EXTRA_X":  -1.5,
		"IMAGE_WIDTH_EXTRA_Y":  -3.75,
		"NETWORK_TYPE":         "tcp",
		"NETWORK_ADDRESS":      "localhost",
		"DEFAULT_IMAGE_PATH":   "default.jpg",
		"CACHE_DIR":            utils.CheckDirectoryFmt(USER_CACHE_DIR),
		"SEEK_OFFSET":          1,
		"REDRAW_INTERVAL":      500,
		"DB_PATH":              utils.CheckDirectoryFmt(USER_CACHE_DIR + "/gompDB"),
	}
)

func ReadConfig() {
	// Parse mpd.conf to set default values.
	ParseMPDConfig()

	for k, v := range defaults {
		viper.SetDefault(k, v)
	}

	if CONFIG_ERR != nil {
		utils.Print("RED", "Couldn't get XDG_CONFIG!")
		panic(CONFIG_ERR)
	}

	if CACHE_ERR != nil {
		utils.Print("RED", "Couldn't get CACHE DIR!")
		panic(CACHE_ERR)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(CONFIG_DIR + "/gomp")

	err := viper.ReadInConfig()
	if err != nil {
		utils.Print("RED", "Could Not Read Config file.\n")
	}

	// Expanding ~ to the User's Home Directory
	viper.Set("MUSIC_DIRECTORY", utils.ExpandHomeDir(viper.GetString("MUSIC_DIRECTORY")))
	viper.Set("DEFAULT_IMAGE_PATH", utils.ExpandHomeDir(viper.GetString("DEFAULT_IMAGE_PATH")))
	viper.Set("CACHE_DIR", utils.ExpandHomeDir(viper.GetString("CACHE_DIR")))
	viper.Set("NETWORK_ADDRESS", utils.ExpandHomeDir(viper.GetString("NETWORK_ADDRESS")))
	viper.Set("DB_PATH", utils.ExpandHomeDir(viper.GetString("DB_PATH")))
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

func ParseMPDConfig() {
	uwconf := CONFIG_DIR + "/mpd/mpd.conf"
	swconf := "/etc/mpd.conf"
	set_defaults := func(path string) {
		m := conf.GenerateMap(path)
		if val, ok := m["music_directory"]; ok {
			defaults["MUSIC_DIRECTORY"] = utils.CheckDirectoryFmt(val.(string))
		}
		if val, ok := m["port"]; ok {
			defaults["MPD_PORT"] = val.(string)
		}
	}
	if utils.FileExists(uwconf) {
		set_defaults(uwconf)
	} else {
		set_defaults(swconf)
	}
}
