# goMP

![](https://img.shields.io/badge/status-alpha-red)

 MPD client inspired by ncmpcpp written in GO

https://user-images.githubusercontent.com/51816057/140478368-5b724b2f-2499-4150-9569-c54734b3d4c1.mp4

# Roadmap

- [ ] Add Functionality to Sort out most played songs
- [x] Add a config parser
- [x] Image Previews
- [ ] Fuzzy Searching
- [ ] Visual Mode (like vim) for updating playlists
- [ ] Music Visualizer

# Prerequisites

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

### Tested on following terminals:

- Alacritty
- ST
- URXVT
