package config

import (
	"errors"
)

var KMAP = map[string]int{
	"TAB":    9,
	"RETURN": 13,
	"ENTER":  13,
	"SPACE":  32,
}

func GetAsciiValue(s string) (int, error) {
	if len([]rune(s)) == 1 {
		char := []rune(s)[0]
		if (int(char) >= 65 && int(char) <= 90) || (int(char) >= 97 && int(char) <= 122) {
			return int(char), nil
		} else {
			return -1, errors.New("Not Found")
		}
	} else if val, ok := KMAP[s]; ok {
		return val, nil
	} else {
		return -1, errors.New("Not Found")
	}
}
