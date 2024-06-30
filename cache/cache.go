package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"

	"github.com/aditya-K2/utils"
)

var (
	CACHE_DIR string
)

func SetCacheDir(path string) {
	CACHE_DIR = utils.CheckDirectoryFmt(path)
}

func Exists(artist, album string) bool {
	if _, err := os.Stat(GenerateName(artist, album)); errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}

func GenerateName(artist, album string) string {
	h := sha256.New()
	h.Write([]byte(artist + album))
	return CACHE_DIR + hex.EncodeToString(h.Sum(nil)) + ".jpg"
}
