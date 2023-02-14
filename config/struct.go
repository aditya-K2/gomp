// ![cover](../docs/assets/readme.png)
// # Configuration
// Configuration of gomp is done through a `config.yml` file in `$HOME/.config/gomp/`
package config

//
// ## Note
// - About Music Directory and MPD_PORT
//
// `gomp` reads your mpd.conf file for setting the defaults.
//
// first it checks for user-wide configuration @ `$HOME/mpd/mpd.conf`
// if that doesn't exists it checks for `/etc/mpd.conf`.
//
// The `MUSIC_DIRECTORY` from config.yml overrides the default values.
// - Default Image
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

	// # Image Rendering
	//
	// The Default Position of the Image without any configuration assumes that you have no font or terminal padding or margin so Image will
	// be rendered at different places in different terminals, Also the TUIs calculates positions with the respect to rows and columns
	// and the image is rendered at pixel positions so the exact position can't be defined [ the app tries its best by calculating
	// the font width and then predicting the position but it is best that you define some extra padding and your own image width ratio
	// in config.yml. Please Read more about it in the [sample_config](https://github.com/aditya-K2/gomp/blob/master/sample_config.yml#L4)

	// ## Additional Padding
	// Extra padding for the placement of the image
	//
	// ![Additional Padding](./assets/info_padding.png)
	//
	// Note: Padding can be negative or positive
	// `ADDITIONAL_PADDING_X` on decrementing will move the image to right
	// and on incrementing will move the image to left
	// similarly `ADDITIONAL_PADDING_Y` on decrementing will move the image UP
	// and on incrementing will move the image DOWN
	AdditionalPaddingX int `mapstructure:"ADDITIONAL_PADDING_X"`
	AdditionalPaddingY int `mapstructure:"ADDITIONAL_PADDING_Y"`

	// ## Image Width
	// Add extra `IMAGE_WIDTH` to the image so that it fits perfectly
	// in the image preview box
	//
	// ![Image Width](./assets/info_width.png)
	//
	// Note: IMAGE_WIDTH_EXTRA_X and IMAGE_WIDTH_EXTRA_Y can be positive or negative
	ExtraImageWidthX float64 `mapstructure:"IMAGE_WIDTH_EXTRA_X"`
	ExtraImageWidthY float64 `mapstructure:"IMAGE_WIDTH_EXTRA_Y"`
	// ## Tutorial
	// for e.g
	//
	// ![Cover Art Position](./assets/default.png)
	//
	// Let's say upon opening gomp for the first time and your image is rendered this way.
	//
	// Here the `Y` Position is too low hence we have to decrease the `ADDITIONAL_PADDING_Y` so that image will be rendered
	// in a better position so we decrement the  `ADDITIONAL_PADDING_Y` by `9`
	//
	// Now the configuration becomes
	// ```yml
	// ADDITIONAL_PADDING_Y : 9
	// ```
	//
	// and the image appears like this:
	//
	// ![Padding](./assets/padding.png)
	//
	// One might be happy the way things turn out but for a perfectionist like me this is not enough.
	// You can notice that the Height of the image is a little bit more than the box height and hence the image is flowing outside the box. Now it's  time to change the `WIDTH_Y`.
	//
	// Width can be changed by defining the `IMAGE_WIDTH_EXTRA_Y` and `IMAGE_WIDTH_EXTRA_X` it act's a little differently think of it like a chunk which is either added or subtracted from the image's original width. We can look at the configuration and realize that the chunk `IMAGE_WIDTH_EXTRA_Y` when subtracted from the original `IMAGE_WIDTH` doesn't get us the proper result and is a little to low. We need to subtract a more bigger chunk, Hence we will increase the magnitude of `IMAGE_WIDTH_EXTRA_Y` or decrease `IMAGE_WIDTH_EXTRA_Y`
	//
	// Now the Configuration becomes:
	// ```yml
	// IMAGE_WIDTH_EXTRA_Y : - 3.2
	// ```
	// and the image appears like this:
	//
	// ![Width](./assets/width.png)
	//
	// Which looks perfect. 🎉
	//
	// Read More about Additional Padding and Image Width in the [sample_config](https://github.com/aditya-K2/gomp/blob/master/sample_config.yml)
	//
	// Please change the configuration according to your needs.

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
	GetCoverArtFromLastFm bool `mapstructure:"GET_COVER_ART_FROM_LAST_FM"`

	// LASTFM_API_KEY: "YOUR API KEY HERE"
	LastFmAPIKey string `mapstructure:"LASTFM_API_KEY"`

	// LASTFM_API_SECRET: "YOUR API SECRET HERE"
	LastFmAPISecret string `mapstructure:"LASTFM_API_SECRET"`
	// ```

	// 3. Auto correct
	//
	// ![Screenshot from 2021-11-13 21-59-46](https://user-images.githubusercontent.com/51816057/141651414-1586577a-cab2-48e2-a24b-1053f8634fbe.png)
	//
	// ```yml
	// LASTFM_AUTO_CORRECT: 1  # 0 means it is turned off
	// ```
	LastFmAPIAutoCorrect int `mapstructure:"LASTFM_AUTO_CORRECT"`
}

func NewConfigS() *ConfigS {
	return &ConfigS{
		AdditionalPaddingX:    12,
		AdditionalPaddingY:    16,
		ExtraImageWidthX:      -1.5,
		ExtraImageWidthY:      -3.75,
		NetworkType:           "tcp",
		NetworkAddress:        "localhost",
		DefaultImagePath:      defaultImagePath,
		CacheDir:              utils.CheckDirectoryFmt(userCacheDir),
		SeekOffset:            1,
		RedrawInterval:        500,
		Colors:                NewColors(),
		GetCoverArtFromLastFm: false,
	}
}
