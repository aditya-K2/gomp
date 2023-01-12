# gomp

<div class="top">
    <style>
        .info {
            display: flex;
            flex-direction: row;
            align-items: center;
        }
        .links{
            padding: 4px;
            padding-top: 0px;
            padding-left: 12px;
        }
        .top {
            display: flex;
            justify-content: center;
        }
    </style>
    <div class="info">
        <img src="docs/assets/logo.png" alt="mascot" width="200" class="mascot"/>
        <div class="links">
            MPD client inspired by ncmpcpp <br>
            builtin <b>cover-art view</b> and <b>LastFM integration.</b> <br>
            <a href="https://aditya-K2.github.io/gomp/"> Documentation </a> |
            <a href="https://github.com/aditya-K2/gomp/discussions">Discussion</a>
        </div>
    </div>
</div>

----

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

### AUR Package [<img src="https://img.shields.io/aur/version/gomp-git">](https://aur.archlinux.org/packages/gomp-git/)

```bash
yay -S gomp-git
```

## Manually

- Install `go` and `git` through your package manager and then

```bash
git clone https://github.com/aditya-K2/gomp && cd gomp && GOFLAGS="-buildmode=pie -trimpath -mod=readonly -modcacherw" go build && sudo install -D gomp -t "/usr/bin/"
```
---
