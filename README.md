# goMP

![](https://img.shields.io/badge/status-alpha-red)

 MPD client inspired by ncmpcpp written in GO

https://user-images.githubusercontent.com/51816057/138585868-92aff5bd-dd7e-46af-bf06-28b83115120b.mp4

# Roadmap

- [ ] Add Functionality to Sort out most played songs
- [ ] Add a config parser ( preferably ***YAML*** )
- [x] Image Previews
- [ ] Fuzzy Searching
- [ ] Visual Mode (like vim) for updating playlists
- [ ] Music Visualizer

# Prerequisites

- Music Player Daemon must be setup
- Go Should Be Installed ( for building )
- Your mpd.conf file in ``~/.config/mpd/mpd.conf`` is automatically parsed to get ``music_directory``
- In [globals.go](https://github.com/aditya-K2/goMP/blob/master/globals.go) set the path for your default Image

```go
var path string = "YOUR DEFAULT IMAGE PATH HERE"
```

- Change the default additional Image padding according to your terminal in [globals.go](https://github.com/aditya-K2/goMP/blob/master/globals.go)

```go
// Change according to your needs

var ADDITIONAL_PADDING_X int = 16
var ADDITIONAL_PADDING_Y int = 24
```
- Increment/Decrement Image Pixels according to your terminal in [globals.go](https://github.com/aditya-K2/goMP/blob/master/globals.go)

```go
// Change according to your needs

var IMAGE_WIDTH_EXTRA_X float32 = -1.5 // Decrementing
var IMAGE_WIDTH_EXTRA_Y float32 = -3.75 // Incrementing

```

# Installing / Building

```bash
git clone https://github.com/aditya-K2/goMP &&
cd goMP &&
go build
```

### Tested on following terminals:

- Alacritty
- ST
- URXVT
