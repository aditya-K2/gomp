# gomp

<div class="info" align="center">
    <br><img src="docs/assets/logo.png" alt="mascot" width="200" class="mascot"/><br>
    MPD client inspired by ncmpcpp <br>
    with builtin <b>cover-art view</b> and <b>LastFM integration.</b> <br>
    <a href="https://aditya-K2.github.io/gomp/"> Documentation </a> |
    <a href="https://github.com/aditya-K2/gomp/discussions">Discussion</a>
</div>

<!-- NEW CHANGE -->

----

![Cover](./docs/assets/readme.png)

# Features

- **Live Config Changes**
- **Cover Art View**
- **LastFM Integration** *(Cover Art)*
- **Fuzzy Searching** *(Global and Buffer Specific)*

# Roadmap

- [ ] Vim-like Visual Mode

# Installing / Building

### AUR Package [<img src="https://img.shields.io/aur/version/gomp-git">](https://aur.archlinux.org/packages/gomp-git/)

```bash
yay -S gomp-git
```

## Manually

- Install `go` and `git` through your package manager and then

```bash
git clone https://github.com/aditya-K2/gomp && cd gomp && GOFLAGS="-buildmode=pie -trimpath -mod=readonly -modcacherw" go build -v && sudo install -D gomp -t "/usr/bin/"
```
---
