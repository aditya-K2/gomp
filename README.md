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
- [ ] Music Visualizer

# Prerequisites

- Music Player Daemon must be setup
- Go Should Be Installed ( for building )
- Set the Path to your mpd DATABASE in progressbar.go

```go
var DBDIR string = "PATH TO YOUR MPD DATABASE HERE"
```

- In imageUtils.go set the path for your default Image

```go
var path string = "YOUR DEFAULT IMAGE PATH HERE"
```

# Installing / Building

```bash
git clone https://github.com/aditya-K2/goMP &&
cd goMP &&
go build
```

## Image Placement

The image is rendered by calculating pixels by multiplying the rows and columns with the font-width which is calculated by dividing the terminal width and height ( Please let me know if there is a better way to do this ) this let's to uneven placement so it is better that you disable the borders for the imagePreview holder.

in [App.go](https://github.com/aditya-K2/goMP/blob/master/App.go)

```go
	imagePreviewer.SetBorder(false)
```

#### With Borders ( Sometimes \[Mostly Maximized Terminals\] )


![Oct28(Thu)02:4134PM](https://user-images.githubusercontent.com/51816057/139225915-b3e30742-65a8-4482-ad38-753646b5082f.png)

#### Without Borders 

![Oct28(Thu)02:4513PM](https://user-images.githubusercontent.com/51816057/139226138-b68ebc22-204c-40f7-a7f2-0dd92b88f72b.png)

Note: Your terminal window padding also affects the tui.
