# gomp

[Documentation](https://aditya-K2.github.io/gomp/) | [Discussion](https://github.com/aditya-K2/gomp/discussions)

![](https://img.shields.io/badge/status-beta-yellow) &nbsp;&nbsp;  [<img src="https://img.shields.io/aur/version/gomp-git">](https://aur.archlinux.org/packages/gomp-git/)

 MPD client inspired by ncmpcpp with builtin cover-art view and LastFM integration.

![Sep11(Sun)11:2026PM](https://user-images.githubusercontent.com/51816057/189541853-282716f1-0515-4ee6-a19a-4989b9de5daf.png)

# Features

- **Functionality to Sort out most played songs**
- **Live Config Changes**
- **Cover Art View**
- **LastFM Integration** *(Cover Art)*
- **Fuzzy Searching** *(Global and Buffer Specific)*

# Roadmap

- [ ] daemonising `gomp` or a `headless` flag
- [ ] Vim-like Visual Mode

# Installing / Building

## AUR Package

```bash
yay -S gomp-git
```

## Manually

- Install `go` and `git` through your package manager and then

```bash
git clone https://github.com/aditya-K2/gomp && cd gomp && GOFLAGS="-buildmode=pie -trimpath -mod=readonly -modcacherw" go build && sudo install -D gomp -t "/usr/bin/"
```
---
