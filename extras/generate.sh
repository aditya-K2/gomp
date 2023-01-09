#!/bin/sh

go build -v ./extras/generate.go
./generate ./config/struct.go ./docs/index.md
./generate ./render/imageUtils.go ./docs/index.md
./generate ./config/kmap.go ./docs/index.md
./generate ./render/lastfm.go ./docs/index.md
./generate ./config/colors.go ./docs/index.md
