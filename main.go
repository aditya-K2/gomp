package main

import (
	"log"
	"strconv"
	"time"

	"github.com/aditya-K2/goMP/config"
	"github.com/fhs/gompd/mpd"
	"github.com/gdamore/tcell/v2"
	"github.com/spf13/viper"
)

var Volume int64
var Random bool
var Repeat bool
var InsidePlaylist bool = true

func main() {
	config.ReadConfig()
	// Connect to MPD server
	conn, err := mpd.Dial("tcp", "localhost:"+viper.GetString("MPD_PORT"))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	r := newRenderer()
	c, _ := conn.CurrentSong()
	if len(c) != 0 {
		r.Start(viper.GetString("MUSIC_DIRECTORY") + c["file"])
	} else {
		r.Start("stop")
	}

	UI := newApplication(*conn, r)

	fileMap, err := conn.GetFiles()
	dirTree := generateDirectoryTree(fileMap)

	UpdatePlaylist(*conn, UI.expandedView)

	_v, _ := conn.Status()
	Volume, _ = strconv.ParseInt(_v["volume"], 10, 64)
	Random, _ = strconv.ParseBool(_v["random"])
	Repeat, _ = strconv.ParseBool(_v["repeat"])

	UI.expandedView.SetDrawFunc(func(s tcell.Screen, x, y, width, height int) (int, int, int, int) {
		if InsidePlaylist {
			UpdatePlaylist(*conn, UI.expandedView)
		} else {
			Update(*conn, dirTree.children, UI.expandedView)
		}
		return UI.expandedView.GetInnerRect()
	})

	var kMap = map[string]func(){
		"showChildrenContent": func() {
			r, _ := UI.expandedView.GetSelection()
			if !InsidePlaylist {
				if len(dirTree.children[r].children) == 0 {
					id, _ := conn.AddId(dirTree.children[r].absolutePath, -1)
					conn.PlayId(id)
				} else {
					Update(*conn, dirTree.children[r].children, UI.expandedView)
					dirTree = &dirTree.children[r]
				}
			} else {
				conn.Play(r)
			}
		},
		"togglePlayBack": func() {
			togglePlayBack(*conn)
		},
		"showParentContent": func() {
			if !InsidePlaylist {
				if dirTree.parent != nil {
					Update(*conn, dirTree.parent.children, UI.expandedView)
					dirTree = dirTree.parent
				}
			}
		},
		"nextSong": func() {
			conn.Next()
		},
		"clearPlaylist": func() {
			conn.Clear()
			if InsidePlaylist {
				UpdatePlaylist(*conn, UI.expandedView)
			}
		},
		"previousSong": func() {
			conn.Previous()
		},
		"addToPlaylist": func() {
			if !InsidePlaylist {
				r, _ := UI.expandedView.GetSelection()
				conn.Add(dirTree.children[r].absolutePath)
			}
		},
		"toggleRandom": func() {
			err := conn.Random(!Random)
			if err == nil {
				Random = !Random
			}
		},
		"toggleRepeat": func() {
			err := conn.Repeat(!Repeat)
			if err == nil {
				Repeat = !Repeat
			}
		},
		"decreaseVolume": func() {
			if Volume <= 0 {
				Volume = 0
			} else {
				Volume -= 10
			}
			conn.SetVolume(int(Volume))
		},
		"increaseVolume": func() {
			if Volume >= 100 {
				Volume = 100
			} else {
				Volume += 10
			}
			conn.SetVolume(int(Volume))
		},
		"navigateToFiles": func() {
			InsidePlaylist = false
			UI.Navbar.Select(1, 0)
			Update(*conn, dirTree.children, UI.expandedView)
		},
		"navigateToPlaylist": func() {
			InsidePlaylist = true
			UI.Navbar.Select(0, 0)
			UpdatePlaylist(*conn, UI.expandedView)
		},
		"navigateToMostPlayed": func() {
			InsidePlaylist = false
			UI.Navbar.Select(2, 0)
		},
		"quit": func() {
			UI.App.Stop()
		},
		"stop": func() {
			conn.Stop()
		},
		"updateDB": func() {
			_, err = conn.Update("")
			if err != nil {
				panic(err)
			}
		},
		"deleteSongFromPlaylist": func() {
			if InsidePlaylist {
				r, _ := UI.expandedView.GetSelection()
				conn.Delete(r, -1)
			}
		},
	}

	config.GenerateKeyMap(kMap)

	UI.expandedView.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if val, ok := config.KEY_MAP[int(e.Rune())]; ok {
			kMap[val]()
			return nil
		} else {
			return e
		}
	})

	go func() {
		for {
			UI.App.Draw()
			time.Sleep(time.Second)
		}
	}()

	if err := UI.App.Run(); err != nil {
		panic(err)
	}
}
