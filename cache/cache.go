package cache

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	CACHE_LIST map[[2]string]string = make(map[[2]string]string)
)

func LoadCache(path string) {
	cacheFileContent, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Could Not Read From Cache File")
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
}

func WriteToCache(path string) {
	b, err := os.Create(path)
	if err == nil {
		for k, v := range CACHE_LIST {
			b.Write([]byte(fmt.Sprintf("%s\t%s\t%s\n", k[0], k[1], v)))
		}
	}
}

func GenerateName(artist, album string) string {
	return fmt.Sprintf("%s-%s.jpg", artist, album)
}
