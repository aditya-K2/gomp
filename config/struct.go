// ![Sep11(Sun)11:2026PM](https://user-images.githubusercontent.com/51816057/189541853-282716f1-0515-4ee6-a19a-4989b9de5daf.png)
// # Configuration
// Configuration of gomp is done through a `config.yml` file in `$HOME/.config/gomp/`
package config

//
// ## Note
// #### About Music Directory and MPD_PORT
//
// `gomp` reads your mpd.conf file for setting the defaults.
//
// first it checks for user-wide configuration @ `$HOME/mpd/mpd.conf`
// if that doesn't exists it checks for `/etc/mpd.conf`
// if there it doesn't find anything (Which shouldn't happen) it uses the
// `MUSIC-DIRECTORY` from config.yml
// #### Default Image
//
// If there is no `DEFAULT_IMAGE_PATH` key in the `config.yml` it downloads a
// default image @ ~/$USER_CACHE_DIR/gomp_default.jpg

import "github.com/aditya-K2/gomp/utils"

type ConfigS struct {

	// ## Network Type
	//
	// By Default gomp assumes that you connect to MPD Server through tcp, But if your MPD Server is configured to expose a unix socket rather than a port, then you can specify network type to "unix"
	// Defaults to `tcp` if not provided.
	//
	// ```yml
	// NETWORK_TYPE : "unix"
	// ```
	//
	// Read More about it in the [sample_config](https://github.com/aditya-K2/gomp/blob/master/sample_config.yml)
	NetworkType string `mapstructure:"NETWORK_TYPE"`

	//
	// ## Network Address
	//
	// The Address of the Host for e.g `"localhost"` or `"/path/to/unix/socket/"` if you are using unix sockets
	// Defaults to `localhost` if not provided.
	//
	// ```yml
	// NETWORK_ADDRESS : "/home/$USER/.mpd/socket"
	// ```
	//
	// Read More about it in the [sample_config](https://github.com/aditya-K2/gomp/blob/master/sample_config.yml)
	NetworkAddress string `mapstructure:"NETWORK_ADDRESS"`

	// ### Default Image Path
	//
	// This is the Fallback Image that will be rendered if gomp doesn't find the embedded cover art or LastFM Cover Art.
	//
	// ```yml
	// DEFAULT_IMAGE_PATH : "/path/to/default/image"
	// ```
	DefaultImagePath string `mapstructure:"DEFAULT_IMAGE_PATH"`

	// ## Cache Directory
	//
	// By Default Images are cached to avoid re-extracting images and making redundant API Calls. The Extracted Images are copied to the `CACHE_DIR`.
	//
	// ```yml
	// CACHE_DIR : "/path/to/the/cache/Directory/"
	// ```
	//
	// Read More about Caching in the [sample_config](https://github.com/aditya-K2/gomp/blob/master/sample_config.yml)
	CacheDir string `mapstructure:"CACHE_DIR"`

	// ## Seek Offset
	// The amount of seconds by which the current song should seek forward or backward.
	//
	// ```yml
	// SEEK_OFFSET : 5
	// ```
	SeekOffset int `mapstructure:"SEEK_OFFSET"`

	// ## Redraw Interval
	//
	// The amount of milliseconds after which the cover art should be redrawn if there is a event.
	//
	// ```yml
	// REDRAW_INTERVAL : 500
	// ```
	RedrawInterval int `mapstructure:"REDRAW_INTERVAL"`

	// ## Database Path
	//
	// The path where the database of playtime of all the songs is stored.
	//
	// ```yml
	// DB_PATH : "~/.cache/gompDB"
	// ```
	DBPath string `mapstructure:"DB_PATH"`

	// ## MPD Port
	// This is the port where your Music Player Daemon Is Running.
	//
	// ##### If not provided gomp looks for the port in mpd.conf in `~/.config/mpd/`
	//
	// ```yml
	// MPD_PORT : "6600"
	// ```
	Port string `mapstructure:"MPD_PORT"`

	// ### Music Directory
	// It is the path to your Music Folder that you have provided to mpd in the `mpd.conf` file.
	// If you do not provide the path to the `MUSIC_DIRECTORY` then gomp parses the mpd.conf file for
	// the `music_directory` key.
	//
	// ```yml
	// MUSIC_DIRECTORY : "~/Music"
	// ```
	//
	// The reason why you need to provide `MUSIC_DIRECTORY` is because the paths
	// to the songs received from mpd are relative the `MUSIC_DIRECTORY`.
	MusicDirectory string  `mapstructure:"MUSIC_DIRECTORY"`
	Colors         *Colors `mapstructure:"COLORS"`

	AdditionalPaddingX    int     `mapstructure:"ADDITIONAL_PADDING_X"`
	AdditionalPaddingY    int     `mapstructure:"ADDITIONAL_PADDING_Y"`
	ExtraImageWidthX      float64 `mapstructure:"IMAGE_WIDTH_EXTRA_X"`
	ExtraImageWidthY      float64 `mapstructure:"IMAGE_WIDTH_EXTRA_Y"`
	LastFmAPIKey          string  `mapstructure:"LASTFM_API_KEY"`
	LastFmAPISecret       string  `mapstructure:"LASTFM_API_SECRET"`
	LastFmAPIAutoCorrect  int     `mapstructure:"LASTFM_AUTO_CORRECT"`
	GetCoverArtFromLastFm bool    `mapstructure:"GET_COVER_ART_FROM_LAST_FM"`
}

func NewConfigS() *ConfigS {
	return &ConfigS{
		AdditionalPaddingX:    12,
		AdditionalPaddingY:    16,
		ExtraImageWidthX:      -1.5,
		ExtraImageWidthY:      -3.75,
		NetworkType:           "tcp",
		NetworkAddress:        "localhost",
		DefaultImagePath:      DefaultImagePath,
		CacheDir:              utils.CheckDirectoryFmt(UserCacheDir),
		SeekOffset:            1,
		RedrawInterval:        500,
		DBPath:                (UserCacheDir + "/gompDB"),
		Colors:                NewColors(),
		GetCoverArtFromLastFm: false,
	}
}
