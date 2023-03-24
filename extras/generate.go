package main

import (
	"fmt"
	"os"
	"strings"
)

func ReadComments(path string) (string, error) {
	if _c, err := os.ReadFile(path); err != nil {
		return "", err
	} else {
		content := ""
		for _, v := range strings.Split(string(_c), "\n") {
			v = strings.Trim(v, "\t ")
			if strings.HasPrefix(v, "//") {
				v = strings.ReplaceAll(v, "\t", "  ")
				content += strings.TrimPrefix(
					strings.TrimPrefix(v, "//"), " ") + "\n"
			}
		}
		return content, nil
	}
}

func main() {
	path := os.Args[1]
	dpath := os.Args[2]
	if content, err := ReadComments(path); err != nil {
		fmt.Printf("Error Reading %s\n", path)
		panic(err)
	} else {
		if file, ferr := os.OpenFile(dpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); ferr != nil {
			fmt.Printf("Error Reading %s\n", dpath)
			panic(ferr)
		} else {
			if _, werr := file.WriteString(content); werr != nil {
				fmt.Printf("Error Writing to %s\n", dpath)
				panic(werr)
			}
		}
	}
}
