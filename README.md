# goMP

![](https://img.shields.io/badge/status-alpha-red)

 MPD client inspired by ncmpcpp written in GO

https://user-images.githubusercontent.com/51816057/140478368-5b724b2f-2499-4150-9569-c54734b3d4c1.mp4

# Roadmap

- [ ] Add Functionality to Sort out most played songs
- [x] Add a config parser
- [x] Image Previews
- [x] User Key Mappings
- [x] Querying LastFM API for getting Album Art
- [ ] Fuzzy Searching
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
|     quit                           |
|     stop                           |
|     updateDB                       |
|     deleteSongFromPlaylist         |

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
- URXVT
