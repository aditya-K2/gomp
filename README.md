# gomp

![](https://img.shields.io/badge/status-beta-yellow) &nbsp;&nbsp;  [<img src="https://img.shields.io/aur/version/gomp-git">](https://aur.archlinux.org/packages/gomp-git/)

 MPD client inspired by ncmpcpp with builtin cover-art view and LastFM integration.

https://user-images.githubusercontent.com/51816057/144759799-b9eecf9e-93ad-43a7-a076-0ae47f03d03c.mp4

## Table of Contents

* [Roadmap](#roadmap)
* [Installing / Building](#installing--building)
	* [AUR Package](#aur-package)
	* [Manually](#manually)
* [Documentation](#documentation)

# Roadmap


- [ ] Add Functionality to Sort out most played songs
- [x] Add a config parser
- [x] Image Previews
- [x] User Key Mappings
- [x] Querying LastFM API for getting Album Art
- [x] Fuzzy Searching
- [ ] Visual Mode (like vim) for updating playlists

# Installing / Building

## AUR Package

```bash
yay -S gomp-git
```

## Manually

```bash
git clone https://github.com/aditya-K2/gomp &&
cd gomp &&
go build
```

# Documentation

The Documentation is hosted [here](https://aditya-K2.github.io/gomp/)

