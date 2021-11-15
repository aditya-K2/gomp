package main

import (
	"sort"
	"strconv"
	"time"

	"github.com/aditya-K2/goMP/config"
	"github.com/aditya-K2/goMP/search"
	"github.com/fhs/gompd/mpd"
	"github.com/gdamore/tcell/v2"
	"github.com/spf13/viper"
)

var CONN *mpd.Client
var UI *Application
var NOTIFICATION_SERVER *NotificationServer
var Volume int64
var Random bool
var Repeat bool
var InsidePlaylist bool = true

func main() {
	config.ReadConfig()
	// Connect to MPD server
	var mpdConnectionError error
	CONN, mpdConnectionError = mpd.Dial("tcp", "localhost:"+viper.GetString("MPD_PORT"))
	if mpdConnectionError != nil {
		panic(mpdConnectionError)
	}
	defer CONN.Close()

	r := newRenderer()
	c, _ := CONN.CurrentSong()
	if len(c) != 0 {
		r.Start(c["file"])
	} else {
		r.Start("stop")
	}

	UI = newApplication(r)

	fileMap, err := CONN.GetFiles()
	dirTree := generateDirectoryTree(fileMap)

	UpdatePlaylist(UI.ExpandedView)

	_v, _ := CONN.Status()
	Volume, _ = strconv.ParseInt(_v["volume"], 10, 64)
	Random, _ = strconv.ParseBool(_v["random"])
	Repeat, _ = strconv.ParseBool(_v["repeat"])

	UI.ExpandedView.SetDrawFunc(func(s tcell.Screen, x, y, width, height int) (int, int, int, int) {
		if InsidePlaylist {
			UpdatePlaylist(UI.ExpandedView)
		} else {
			Update(dirTree.children, UI.ExpandedView)
		}
		return UI.ExpandedView.GetInnerRect()
	})

	NOTIFICATION_SERVER = NewNotificationServer()
	NOTIFICATION_SERVER.Start()

	var FUNC_MAP = map[string]func(){
		"showChildrenContent": func() {
			r, _ := UI.ExpandedView.GetSelection()
			if !InsidePlaylist {
				if len(dirTree.children[r].children) == 0 {
					id, _ := CONN.AddId(dirTree.children[r].absolutePath, -1)
					CONN.PlayId(id)
				} else {
					Update(dirTree.children[r].children, UI.ExpandedView)
					dirTree = &dirTree.children[r]
				}
			} else {
				CONN.Play(r)
			}
		},
		"togglePlayBack": func() {
			togglePlayBack()
		},
		"showParentContent": func() {
			if !InsidePlaylist {
				if dirTree.parent != nil {
					Update(dirTree.parent.children, UI.ExpandedView)
					dirTree = dirTree.parent
				}
			}
		},
		"nextSong": func() {
			CONN.Next()
		},
		"clearPlaylist": func() {
			CONN.Clear()
			if InsidePlaylist {
				UpdatePlaylist(UI.ExpandedView)
			}
			NOTIFICATION_SERVER.Send("PlayList Cleared")
		},
		"previousSong": func() {
			CONN.Previous()
		},
		"addToPlaylist": func() {
			if !InsidePlaylist {
				r, _ := UI.ExpandedView.GetSelection()
				CONN.Add(dirTree.children[r].absolutePath)
			}
		},
		"toggleRandom": func() {
			err := CONN.Random(!Random)
			if err == nil {
				Random = !Random
			}
		},
		"toggleRepeat": func() {
			err := CONN.Repeat(!Repeat)
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
			CONN.SetVolume(int(Volume))
		},
		"increaseVolume": func() {
			if Volume >= 100 {
				Volume = 100
			} else {
				Volume += 10
			}
			CONN.SetVolume(int(Volume))
		},
		"navigateToFiles": func() {
			InsidePlaylist = false
			UI.Navbar.Select(1, 0)
			Update(dirTree.children, UI.ExpandedView)
		},
		"navigateToPlaylist": func() {
			InsidePlaylist = true
			UI.Navbar.Select(0, 0)
			UpdatePlaylist(UI.ExpandedView)
		},
		"navigateToMostPlayed": func() {
			InsidePlaylist = false
			UI.Navbar.Select(2, 0)
		},
		"quit": func() {
			UI.App.Stop()
		},
		"stop": func() {
			CONN.Stop()
			NOTIFICATION_SERVER.Send("Playback Stopped")
		},
		"updateDB": func() {
			_, err = CONN.Update("")
			if err != nil {
				panic(err)
			}
			NOTIFICATION_SERVER.Send("Database Updated")
		},
		"deleteSongFromPlaylist": func() {
			if InsidePlaylist {
				r, _ := UI.ExpandedView.GetSelection()
				CONN.Delete(r, -1)
			}
		},
		"FocusSearch": func() {
			UI.App.SetFocus(UI.SearchBar)
		},
	}

	config.GenerateKeyMap(FUNC_MAP)
	a, _ := GenerateArtistTree()
	UI.SearchBar.SetAutocompleteFunc(func(c string) []string {
		if c != "" && c != " " && c != "  " {
			var p search.PairList
			for k2, v := range a {
				p = append(p, search.Pair{k2, search.GetLevenshteinDistance(c, k2)})
				for k1, v1 := range v {
					p = append(p, search.Pair{k1, search.GetLevenshteinDistance(c, k1)})
					for k := range v1 {
						p = append(p, search.Pair{k, search.GetLevenshteinDistance(c, k)})
					}
				}
			}

			sort.Sort(p)
			var suggestions []string
			i := 0
			for _, k := range p {
				if i == 10 {
					break
				}
				_, _, w, _ := UI.SearchBar.GetRect()
				suggestions = append(suggestions, getFormattedString(k.Key, w-2))
				i++
			}
			return suggestions
		} else {
			return make([]string, 0)
		}
	})
	UI.ExpandedView.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if val, ok := config.KEY_MAP[int(e.Rune())]; ok {
			FUNC_MAP[val]()
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
