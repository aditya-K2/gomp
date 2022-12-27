package config

import (
	"os"

	"github.com/aditya-K2/gomp/config/conf"
	"github.com/aditya-K2/gomp/utils"
	"github.com/spf13/viper"
)

type ConfigS struct {
	AdditionalPaddingX    int     `mapstructure:"ADDITIONAL_PADDING_X"`
	AdditionalPaddingY    int     `mapstructure:"ADDITIONAL_PADDING_Y"`
	ExtraImageWidthX      float64 `mapstructure:"IMAGE_WIDTH_EXTRA_X"`
	ExtraImageWidthY      float64 `mapstructure:"IMAGE_WIDTH_EXTRA_Y"`
	NetworkType           string  `mapstructure:"NETWORK_TYPE"`
	NetworkAddress        string  `mapstructure:"NETWORK_ADDRESS"`
	DefaultImagePath      string  `mapstructure:"DEFAULT_IMAGE_PATH"`
	CacheDir              string  `mapstructure:"CACHE_DIR"`
	SeekOffset            int     `mapstructure:"SEEK_OFFSET"`
	RedrawInterval        int     `mapstructure:"REDRAW_INTERVAL"`
	DBPath                string  `mapstructure:"DB_PATH"`
	LastFmAPIKey          string  `mapstructure:"LASTFM_API_KEY"`
	LastFmAPISecret       string  `mapstructure:"LASTFM_API_SECRET"`
	LastFmAPIAutoCorrect  bool    `mapstructure:"LASTFM_AUTO_CORRECT"`
	GetCoverArtFromLastFm string  `mapstructure:"GET_COVER_ART_FROM_LAST_FM"`
	Port                  string  `mapstructure:"MPD_PORT"`
	MusicDirectory        string  `mapstructure:"MUSIC_DIRECTORY"`
}

var (
	ConfigDir, configErr   = os.UserConfigDir()
	UserCacheDir, cacheErr = os.UserCacheDir()
	Config                 = NewConfigS()
)

func NewConfigS() *ConfigS {
	return &ConfigS{
		AdditionalPaddingX: 12,
		AdditionalPaddingY: 16,
		ExtraImageWidthX:   -1.5,
		ExtraImageWidthY:   -3.75,
		NetworkType:        "tcp",
		NetworkAddress:     "localhost",
		DefaultImagePath:   "default.jpg",
		CacheDir:           utils.CheckDirectoryFmt(UserCacheDir),
		SeekOffset:         1,
		RedrawInterval:     500,
		DBPath:             (UserCacheDir + "/gompDB"),
	}
}

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

	err := viper.ReadInConfig()
	if err != nil {
		utils.Print("RED", "Could Not Read Config file.\n")
	}
	viper.Unmarshal(Config)

	// Expanding ~ to the User's Home Directory
	Config.MusicDirectory = utils.ExpandHomeDir(Config.MusicDirectory)
	Config.DefaultImagePath = utils.ExpandHomeDir(Config.DefaultImagePath)
	Config.CacheDir = utils.ExpandHomeDir(Config.CacheDir)
	Config.NetworkAddress = utils.ExpandHomeDir(Config.NetworkAddress)
	Config.DBPath = utils.ExpandHomeDir(Config.DBPath)
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
