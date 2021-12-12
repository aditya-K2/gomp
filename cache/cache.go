package cache

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aditya-K2/goMP/config"
)

var (
	CACHE_DIR   string
	DEFAULT_IMG string
)

func SetCacheDir(path string) {
	CACHE_DIR = config.CheckDirectoryFmt(path)
}

func Exists(artist, album string) bool {
	if _, err := os.Stat(GenerateName(artist, album)); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}

func GenerateName(artist, album string) string {
	if (artist == "" && album == "") || (artist == " " && album == " ") {
		return CACHE_DIR + "UnknownArtist-UnknownAlbum.jpg"
	}
	return CACHE_DIR + strings.Replace(strings.Replace(fmt.Sprintf("%s-%s.jpg", artist, album), " ", "_", -1), "/", "_", -1)
}
