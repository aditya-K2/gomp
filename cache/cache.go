package cache

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	USER_CACHE_DIR, err                      = os.UserCacheDir()
	CACHE_LIST          map[[2]string]string = make(map[[2]string]string)
	CACHE_DIR           string               = USER_CACHE_DIR
)

func SetCacheDir(path string) {
	CACHE_DIR = path
}

func LoadCache(path string) error {
	cacheFileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.New("Could Not Read From Cache File")
	}
	lineSlice := strings.Split(string(cacheFileContent), "\n")
	for _, line := range lineSlice {
		if len(line) != 0 {
			param := strings.Split(line, "\t")
			if len(param) == 3 {
				CACHE_LIST[[2]string{param[0], param[1]}] = param[2]
			}
		}
	}
	return nil
}

func GetFromCache(artist, album string) (string, error) {
	if val, ok := CACHE_LIST[[2]string{artist, album}]; ok {
		return val, nil
	} else {
		return "", errors.New("Element Not In Cache")
	}
}

func AddToCache(artist, album string) string {
	fileName := CACHE_DIR + GenerateName(artist, album)
	CACHE_LIST[[2]string{artist, album}] = fileName
	return fileName
}

func WriteCache(path string) {
	b, err := os.Create(path)
	if err == nil {
		for k, v := range CACHE_LIST {
			b.Write([]byte(fmt.Sprintf("%s\t%s\t%s\n", k[0], k[1], v)))
		}
	}
}

func GenerateName(artist, album string) string {
	return strings.Replace(strings.Replace(fmt.Sprintf("%s-%s.jpg", artist, album), " ", "_", -1), "/", "_", -1)
}
