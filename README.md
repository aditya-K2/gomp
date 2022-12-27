# gomp

[Documentation](https://aditya-K2.github.io/gomp/) | [Discussion](https://github.com/aditya-K2/gomp/discussions)

![](https://img.shields.io/badge/status-beta-yellow) &nbsp;&nbsp;  [<img src="https://img.shields.io/aur/version/gomp-git">](https://aur.archlinux.org/packages/gomp-git/)

 MPD client inspired by ncmpcpp with builtin cover-art view and LastFM integration.

![Sep11(Sun)11:2026PM](https://user-images.githubusercontent.com/51816057/189541853-282716f1-0515-4ee6-a19a-4989b9de5daf.png)

# Roadmap


- [x] Add Functionality to Sort out most played songs
    - [ ] daemonising `gomp` or a `headless` flag
- [x] Add a config parser
- [x] Image Previews
- [x] User Key Mappings
- [x] Querying LastFM API for getting Album Art
- [x] Fuzzy Searching
- Config Related Changes
    - [x] Support Live Changes in the Config File
    - [ ] Aesthetic Customisability ***(colors etc.)***
- [ ] Visual Mode (like vim) for updating playlists

# Installing / Building

## AUR Package

```bash
yay -S gomp-git
```

## Manually

```bash
git clone https://github.com/aditya-K2/gomp && cd gomp && GOFLAGS="-buildmode=pie -trimpath -mod=readonly -modcacherw" go build
```
---
