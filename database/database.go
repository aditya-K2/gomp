package database

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"sort"
	"strings"
	"time"

	"github.com/aditya-K2/gomp/utils"
	"github.com/fhs/gompd/v2/mpd"
)

var (
	fmap      = make(map[string]time.Duration)
	fch       = make(chan Payload)
	dbPath    string
	dbSep     = "#"
	fchannels = []chan SubPayload{}
)

type Payload struct {
	path     string
	duration time.Duration
}

type SubPayload struct {
	Fmap  map[string]time.Duration
	Slice []string
}

func SetDB(path string) {
	dbPath = path
}

func ValidRow(row []string) bool {
	if len(row) != 2 {
		return false
	}

	if row[0] == "" || row[1] == "" {
		return false
	}

	if strings.HasPrefix(row[0], "#") ||
		strings.HasPrefix(row[1], "#") {
		return false
	}

	return true
}

func Read() {
	if !utils.FileExists(dbPath) {
		if err := ioutil.WriteFile(dbPath, []byte{}, fs.ModeAppend); err != nil {
			utils.Print("RED", fmt.Sprintf("Error Creating Database: %s\n", dbPath))
			panic(err)
		}
	}
	var (
		readErr error
		content []byte
	)
	if content, readErr = ioutil.ReadFile(dbPath); readErr != nil {
		utils.Print("RED", fmt.Sprintf("Error Reading Database: %s\n", dbPath))
		panic(readErr)
	} else {
		for _, v := range strings.Split(string(content), "\n") {
			row := strings.Split(v, dbSep)
			for _, v := range row {
				v = strings.TrimSpace(v)
			}
			if ValidRow(row) {
				if du, err := time.ParseDuration(row[1]); err == nil {
					fmap[row[0]] = du
				}
			}
		}
	}
}

func Write() {
	var content string
	for k, v := range fmap {
		content += fmt.Sprintf("%s%s%s\n", k, dbSep, v.String())
	}
	ioutil.WriteFile(dbPath, []byte(content), fs.ModeAppend)
}

func Start() {
	var (
		path = "SKIP"
		dur  = time.Second
	)
	go DBRoutine(Payload{path, dur})
}

func Send(c mpd.Attrs, dur time.Duration) {
	var path string
	if len(c) == 0 {
		path = "SKIP"
	} else {
		path = c["file"]
	}
	fch <- Payload{path, dur}
}

func Add(path string, dur time.Duration) {
	if _, ok := fmap[path]; ok {
		fmap[path] += dur
	} else {
		fmap[path] = dur
	}
}

func DBRoutine(payload Payload) {
	if payload.path != "SKIP" {
		Add(payload.path, payload.duration)
		Publish()
	}
	payload = <-fch
	DBRoutine(payload)
}

func GetSlice() []string {
	_fmap := fmap
	keys := []string{}
	for k := range _fmap {
		keys = append(keys, k)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return fmap[keys[i]] > fmap[keys[j]]
	})
	return keys
}

func Subscribe(ch chan SubPayload) {
	fchannels = append(fchannels, ch)
}

func Publish() {
	for _, v := range fchannels {
		go func(v chan SubPayload) {
			v <- SubPayload{fmap, GetSlice()}
		}(v)
	}
}
