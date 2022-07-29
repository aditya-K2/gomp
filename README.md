# gomp

![](https://img.shields.io/badge/status-beta-yellow) &nbsp;&nbsp;  [<img src="https://img.shields.io/aur/version/gomp-git">](https://aur.archlinux.org/packages/gomp-git/)

 MPD client inspired by ncmpcpp with builtin cover-art view and LastFM integration.

![screenshot](https://user-images.githubusercontent.com/51816057/147781035-69eeeb1c-cd62-4e42-8e71-3b07538704e8.png)

## Table of Contents


* [Roadmap](#roadmap)
* [Installing / Building](#installing--building)
    * [AUR Package](#aur-package)
    * [Manually](#manually)
* [Setting Up](#setting-up)

#### Drop your Feedback or Questions about the Documentation [Here](https://github.com/aditya-K2/gomp/issues/25)


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
---

# Setting Up

- [Configuring](https://aditya-K2.github.io/gomp/)
- [Video Showcase](https://github.com/aditya-K2/gomp/tree/master/extras/showcase.md)
