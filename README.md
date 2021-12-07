# goMP

![](https://img.shields.io/badge/status-alpha-red)

 MPD client inspired by ncmpcpp written in GO



https://user-images.githubusercontent.com/51816057/144759799-b9eecf9e-93ad-43a7-a076-0ae47f03d03c.mp4



# Table Of Contents

- [goMP](#gomp)
- [Roadmap](#roadmap)
- [Setting Up](#setting-up)
- [Installing / Building](#installing--building)
- [Configuration](#configuration)
  - [Image Rendering :](#image-rendering-)
  - [Key Mappings](#key-mappings)
  - [Getting Album Art from LastFm API](#getting-album-art-from-lastfm-api)
    - [Tested on following terminals:](#tested-on-following-terminals)

# Roadmap


- [ ] Add Functionality to Sort out most played songs
- [x] Add a config parser
- [x] Image Previews
- [x] User Key Mappings
- [x] Querying LastFM API for getting Album Art
- [x] Fuzzy Searching
- [ ] Visual Mode (like vim) for updating playlists

# Setting Up

- Music Player Daemon must be setup
- Go Should Be Installed ( for building )
- Make a YAML/TOML file in ``$HOME/.config/goMP`` named config.yml / config.toml
- Read the sample_config.yml for config options

# Installing / Building

```bash
git clone https://github.com/aditya-K2/goMP &&
cd goMP &&
go build
```

# Configuration

## Image Rendering :

The Default Position of the Image without any configuration assumes that you have no padding or margin so Image will
be rendered at different places in different terminals, Also the TUI calculates positions with the help of rows and columns
and the image is rendered at pixel positions so the exact position can't be defined [ the app tries its best by calculating
the font width and then predicting the position but it is best that you define extra padding and own image width ratio
in config.yml. Please Read more about it in the [sample_config](https://github.com/aditya-K2/goMP/blob/master/sample_config.yml)

# Setting Up Cache Directory and Cache File

The Images that are extracted are cached.
In the `config.yml` file add the following

```yml
CACHE_DIR : "/path/to/the/cache/Directory"
CACHE_FILE : "/path/to/the/cache/file" # Cache file stores the path to all the images ( think of it like a database. )
```

## Key Mappings

Following Keys can be used for Mappings

| Keys            | Using them in Config  |
|-----------------|-----------------------|
| a-z             | a-z                   |
| A-Z             | A-z                   |
| {,},(,),[,],<,> | {,},(,),[,],<,>       |
| Enter(Return)   | ENTER/RETURN          |
| Tab             | TAB                   |
| Space           | SPACE                 |

See config/kMap.go for more information

For mapping a key to some function use the following format:


```yml
Function: [ firstMapping, secondMapping, thirdMapping]
```
for.eg


```yml
togglePlayBack : [ "p", "TAB", "[" ] # using the quotes is neccessary.
```

Following functions are provided :

|          Functions                 |
|------------------------------------|
|     showChildrenContent            |
|     togglePlayBack                 |
|     showParentContent              |
|     nextSong                       |
|     clearPlaylist                  |
|     previousSong                   |
|     addToPlaylist                  |
|     toggleRandom                   |
|     toggleRepeat                   |
|     decreaseVolume                 |
|     increaseVolume                 |
|     navigateToFiles                |
|     navigateToPlaylist             |
|     navigateToMostPlayed           |
|     navigateToSearch               |
|     quit                           |
|     stop                           |
|     updateDB                       |
|     deleteSongFromPlaylist         |
|     FocusSearch                    |

## Getting Album Art from [LastFm API](https://www.last.fm/api)

1. Generate API Key [here](https://www.last.fm/login?next=%2Fapi%2Faccount%2Fcreate%3F_pjax%3D%2523content)

   ![Screenshot from 2021-11-13 21-54-54](https://user-images.githubusercontent.com/51816057/141651276-f76a5c7f-65fe-4a1a-b130-18cdf67dd471.png)

2. Add the api key and api secret to config.yml

```yml

GET_COVER_ART_FROM_LAST_FM : "TRUE" # Turn On Getting Album art from lastfm api
LASTFM_API_KEY: "YOUR API KEY HERE"
LASTFM_API_SECRET: "YOUR API SECRET HERE"
```
3. Auto correct

![Screenshot from 2021-11-13 21-59-46](https://user-images.githubusercontent.com/51816057/141651414-1586577a-cab2-48e2-a24b-1053f8634fbe.png)
:

```yml

LASTFM_AUTO_CORRECT: 1  # 0 means it is turned off

```

### Tested on following terminals:

- Alacritty
- ST
- URXVT ( URXVT sometimes crashes when downloading coverart from lastfm )
