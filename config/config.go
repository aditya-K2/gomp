package config

import (
	"os"

	"github.com/aditya-K2/gomp/config/conf"
	"github.com/aditya-K2/gomp/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	ConfigDir, configErr   = os.UserConfigDir()
	UserCacheDir, cacheErr = os.UserCacheDir()
	Config                 = NewConfigS()
	OnConfigChange         func()
	DefaultImageLink       = "https://raw.githubusercontent.com/aditya-K2/gomp/master/docs/assets/logo.png"
	DefaultImagePath       = UserCacheDir + "/gomp_default.jpg"
)

func ReadConfig() {
	// Parse mpd.conf to set default values.
	ParseMPDConfig()

	if configErr != nil {
		utils.Print("RED", "Couldn't get $XDG_CONFIG!")
		panic(configErr)
	}

	if cacheErr != nil {
		utils.Print("RED", "Couldn't get CACHE DIR!")
		panic(cacheErr)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(ConfigDir + "/gomp")

	if err := viper.ReadInConfig(); err != nil {
		utils.Print("RED", "Could Not Read Config file.\n")
	} else {
		viper.Unmarshal(Config)
	}

	// Expanding ~ to the User's Home Directory
	expandHome := func() {
		Config.MusicDirectory = utils.ExpandHomeDir(Config.MusicDirectory)
		Config.DefaultImagePath = utils.ExpandHomeDir(Config.DefaultImagePath)
		Config.CacheDir = utils.ExpandHomeDir(Config.CacheDir)
		Config.NetworkAddress = utils.ExpandHomeDir(Config.NetworkAddress)
		Config.DBPath = utils.ExpandHomeDir(Config.DBPath)
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		viper.Unmarshal(Config)
		expandHome()
		OnConfigChange()
	})
	viper.WatchConfig()

	if Config.DefaultImagePath == DefaultImagePath {
		if !utils.FileExists(DefaultImagePath) {
			utils.Print("BLUE", "Default Image Not Provided Downloading Default Image From: ")
			utils.Print("YELLOW", DefaultImageLink+"\n")
			if _, err := utils.DownloadImage(DefaultImageLink, DefaultImagePath); err != nil {
				utils.Print("RED", "Couldn't Download Default Image!\n")
				os.Exit(-1)
			} else {
				utils.Print("CYAN", "Downloaded @ ")
				utils.Print("PURPLE", DefaultImagePath+"\n")
			}
		}
	}

	expandHome()
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
	uwconf := ConfigDir + "/mpd/mpd.conf"
	swconf := "/etc/mpd.conf"
	set_defaults := func(path string) {
		m := conf.GenerateMap(path)
		if val, ok := m["music_directory"]; ok {
			Config.MusicDirectory = utils.CheckDirectoryFmt(val.(string))
		}
		if val, ok := m["port"]; ok {
			Config.Port = val.(string)
		}
	}
	if utils.FileExists(uwconf) {
		set_defaults(uwconf)
	} else {
		set_defaults(swconf)
	}
}
