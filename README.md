# goMP

![](https://img.shields.io/badge/status-alpha-red)

 MPD client inspired by ncmpcpp written in GO

https://user-images.githubusercontent.com/51816057/138585868-92aff5bd-dd7e-46af-bf06-28b83115120b.mp4

# Roadmap

- [ ] Add Functionality to Sort out most played songs
- [ ] Add a config parser ( preferably ***YAML*** )
- [x] Image Previews
	- The Image Previews are working but the placement is very hacky and needs some work 
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
