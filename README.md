# goMP

![](https://img.shields.io/badge/status-usable-blue)

 MPD client inspired by ncmpcpp written in GO

https://user-images.githubusercontent.com/51816057/137694586-199e8c0c-aa5b-473f-9657-ea399bb582a8.mp4

# Roadmap

- [ ] Add Functionality to Sort out most played songs
- [ ] Add a config parser ( preferably ***YAML*** )
- [ ] Ueberzug Integration
- [ ] Fuzzy Searching
- [ ] Visual Mode (like vim) for updating playlists

# Prerequisites

- Music Player Daemon must be setup
- Go Should Be Installed ( for building )
- Set the Path to your mpd DATABASE in progressbar.go

```go
	var DBDIR string = "PATH TO YOUR MPD DATABASE HERE"
```

# Installing / Building

```bash
	git clone https://github.com/aditya-K2/goMP &&
	cd goMP &&
	go build
```

## Most of the Key-bindings are same as ncmpcpp
