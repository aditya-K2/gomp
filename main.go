package main

import (
	"strconv"
	"time"

	"github.com/aditya-K2/gomp/cache"
	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/gomp/config"
	"github.com/aditya-K2/gomp/database"
	"github.com/aditya-K2/gomp/ui"
	"github.com/aditya-K2/gomp/ui/notify"
	"github.com/aditya-K2/gomp/utils"
	"github.com/aditya-K2/gomp/views"
	"github.com/aditya-K2/gomp/watchers"

	"github.com/aditya-K2/fuzzy"
	"github.com/fhs/gompd/v2/mpd"
	"github.com/gdamore/tcell/v2"
)

func main() {
	config.OnConfigChange = watchers.OnConfigChange
	config.ReadConfig()
	var mpdConnectionError error
	client.Conn, mpdConnectionError = mpd.Dial(
		utils.GetNetwork(config.Config.NetworkType,
			config.Config.Port,
			config.Config.NetworkAddress))
	if mpdConnectionError != nil {
		utils.Print("RED", "Could Not Connect to MPD Server\n")
		utils.Print("GREEN", "Make Sure You Mention the Correct MPD Port in the config file.\n")
		panic(mpdConnectionError)
	}
	Conn := client.Conn
	defer Conn.Close()

	cache.SetCacheDir(config.Config.CacheDir)

	watchers.Init()
	ui.Ui = ui.NewApplication()
	ui.Ui.ProgressBar.SetProgressFunc(watchers.ProgressFunction)

	fileMap, err := Conn.ListAllInfo("/")
	if err != nil {
		utils.Print("RED", "Could Not Generate the File Map\n")
		utils.Print("GREEN", "Make Sure You Mention the Correct MPD Port in the config file.\n")
		panic(err)
	}

	// Generating the Directory Tree for File Navigation.
	client.DirTree = client.GenerateDirectoryTree(fileMap)

	var Volume int64
	var Random, Repeat bool
	var SeekOffset = config.Config.SeekOffset
	var SeekFunc = func(back bool) {
		if status, err := Conn.Status(); err != nil {
			notify.Send("Could not get MPD Status")
		} else {
			if status["state"] == "play" {
				var stime time.Duration
				if back {
					stime = -1 * time.Second * time.Duration(SeekOffset)
				} else {
					stime = time.Second * time.Duration(SeekOffset)
				}
				if err := Conn.SeekCur(stime, true); err != nil {
					notify.Send("Could Not Seek Forward in the Song")
				}
			}
		}
	}

	if _v, err := Conn.Status(); err != nil {
		utils.Print("RED", "Could Not Get the MPD Status\n")
		panic(err)
	} else {
		// Setting Volume, Random and Repeat Values
		Volume, _ = strconv.ParseInt(_v["volume"], 10, 64)
		Random, _ = strconv.ParseBool(_v["random"])
		Repeat, _ = strconv.ParseBool(_v["repeat"])
	}

	ArtistTree, err := client.GenerateArtistTree()
	if err != nil {
		utils.Print("RED", "Could Not Generate the ArtistTree\n")
		utils.Print("GREEN", "Make Sure You Mention the Correct MPD Port in the config file.\n")
		panic(err)
	}

	// Used for Fuzzy Searching
	ArtistTreeContent := utils.ConvertToArray(ArtistTree)

	notify.Init()

	// This Function Is Responsible for Changing the Focus it uses the Focus Map and Based on it Chooses
	// the Draw Function
	watchers.StartPlaylistWatcher()
	watchers.StartMPListener()
	watchers.StartRectWatcher()
	database.Publish()

	views.SetCurrentView(&views.PView)
	ui.Ui.ExpandedView.SetDrawFunc(func(s tcell.Screen, x, y, width, height int) (int, int, int, int) {
		views.GetCurrentView().Update(ui.Ui.ExpandedView)
		return ui.Ui.ExpandedView.GetInnerRect()
	})

	// Function Maps is used For Mapping Keys According to the Value mapped to the Key the respective Function is called
	// For e.g. in the config if the User Maps T to togglePlayBack then whenever in the input handler the T is received
	// the respective function in this case togglePlayBack is called.
	var FuncMap = map[string]func(){
		"showChildrenContent": func() {
			views.GetCurrentView().ShowChildrenContent()
		},
		"togglePlayBack": func() {
			if err := client.TogglePlayBack(); err != nil {
				notify.Send("Could not Toggle Play Back")
			}
		},
		"showParentContent": func() {
			views.GetCurrentView().ShowParentContent()
		},
		"nextSong": func() {
			if err := Conn.Next(); err != nil {
				notify.Send("Could not Select the Next Song")
			}
		},
		"clearPlaylist": func() {
			if err := Conn.Clear(); err != nil {
				notify.Send("Could not Clear the Playlist")
			} else {
				if views.PView.Playlist, err = client.Conn.PlaylistInfo(-1, -1); err != nil {
					utils.Print("RED", "Couldn't get the current Playlist.\n")
				} else {
					notify.Send("Playlist Cleared!")
				}
			}
		},
		"previousSong": func() {
			if err := Conn.Previous(); err != nil {
				notify.Send("Could Not Select the Previous Song")
			}
		},
		"addToPlaylist": func() {
			views.GetCurrentView().AddToPlaylist()
		},
		"toggleRandom": func() {
			if err := Conn.Random(!Random); err == nil {
				Random = !Random
			}
		},
		"toggleRepeat": func() {
			if err := Conn.Repeat(!Repeat); err == nil {
				Repeat = !Repeat
			}
		},
		"decreaseVolume": func() {
			if Volume <= 0 {
				Volume = 0
			} else {
				Volume -= 10
			}
			if err := Conn.SetVolume(int(Volume)); err != nil {
				notify.Send("Could Not Decrease the Volume")
			}
		},
		"increaseVolume": func() {
			if Volume >= 100 {
				Volume = 100
			} else {
				Volume += 10
			}
			if err := Conn.SetVolume(int(Volume)); err != nil {
				notify.Send("Could Not Increase the Volume")
			}
		},
		"navigateToFiles": func() {
			views.SetCurrentView(views.FView)
			ui.Ui.Navbar.Select(1, 0)
			views.FView.Update(ui.Ui.ExpandedView)
		},
		"navigateToPlaylist": func() {
			views.SetCurrentView(&views.PView)
			ui.Ui.Navbar.Select(0, 0)
			views.PView.Update(ui.Ui.ExpandedView)
		},
		"navigateToMostPlayed": func() {
			views.SetCurrentView(&views.MPView)
			ui.Ui.Navbar.Select(2, 0)
			views.MPView.Update(ui.Ui.ExpandedView)
		},
		"navigateToSearch": func() {
			views.SetCurrentView(views.SView)
			ui.Ui.Navbar.Select(3, 0)
			views.SView.Update(ui.Ui.ExpandedView)
		},
		"quit": func() {
			views.GetCurrentView().Quit()
		},
		"stop": func() {
			if err := Conn.Stop(); err != nil {
				notify.Send("Could not Stop the Playback")
			} else {
				notify.Send("Playback Stopped")
			}
		},
		"updateDB": func() {
			_, err = Conn.Update("")
			if err != nil {
				notify.Send("Could Not Update the Database")
			} else {
				notify.Send("Database Updated")
			}
		},
		"deleteSongFromPlaylist": func() {
			views.GetCurrentView().DeleteSongFromPlaylist()
		},
		"FocusSearch": func() {
			ui.Ui.App.SetFocus(ui.Ui.SearchBar)
		},
		"FocusBuffSearch": func() {
			views.GetCurrentView().FocusBuffSearchView()
		},
		"SeekForward": func() {
			SeekFunc(false)
		},
		"SeekBackward": func() {
			SeekFunc(true)
		},
	}

	// Generating the Key Map Based on the Function Map Here Basically the Values will be flipped
	// In the config if togglePlayBack is mapped to [ T , P, SPACE ] then here Basically we will receive a map
	// for each event T, P, SPACE mapped to the same function togglePlayBack
	config.GenerateKeyMap(FuncMap)

	ui.Ui.SearchBar.SetAutocompleteFunc(func(c string) []string {
		if views.GetCurrentView().GetViewName() == "BuffSearchView" {
			return nil
		} else {
			if c != "" && c != " " && c != "  " {
				_, _, w, _ := ui.Ui.SearchBar.GetRect()
				matches := fuzzy.Find(c, ArtistTreeContent)
				var suggestions []string
				for i, match := range matches {
					if i == 10 {
						break
					}
					suggestions = append(suggestions, utils.GetFormattedString(match.Str, w-2))
				}
				return suggestions
			} else {
				return make([]string, 0)
			}
		}
	})

	// Input Handler
	ui.Ui.ExpandedView.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if val, ok := config.KEY_MAP[int(e.Rune())]; ok {
			FuncMap[val]()
			return nil
		} else {
			if views.GetCurrentView().GetViewName() == "PlaylistView" {
				if e.Rune() == 'j' || e.Rune() == 'k' {
					if len(views.PView.Playlist) == 0 {
						notify.Send("Empty Playlist")
						return nil
					}
				}
			} else if views.GetCurrentView().GetViewName() == "SearchView" {
				if e.Rune() == 'j' || e.Rune() == 'k' {
					if client.SearchContentSlice == nil || len(client.SearchContentSlice) == 0 {
						notify.Send("No Search Results")
						return nil
					}
				}
			}
			return e
		}
	})

	ui.Ui.SearchBar.SetDoneFunc(func(e tcell.Key) {
		if e == tcell.KeyEnter {
			ui.Ui.ExpandedView.Select(0, 0)
			if views.GetCurrentView().GetViewName() == "BuffSearchView" {
				ui.Ui.App.SetFocus(ui.Ui.ExpandedView)
			} else {
				views.SetCurrentView(views.SView)
				client.SearchContentSlice = nil
				client.SearchContentSlice, err = client.GenerateContentSlice(ui.Ui.SearchBar.GetText())
				if err != nil {
					notify.Send("Could Not Retrieve the Results")
				} else {
					ui.Ui.SearchBar.SetText("")
					ui.Ui.App.SetFocus(ui.Ui.ExpandedView)
					ui.Ui.Navbar.Select(3, 0)
				}
			}
		}
		if e == tcell.KeyEscape {
			if views.GetCurrentView().GetViewName() == "BuffSearchView" {
				client.Matches = nil
			}
			ui.Ui.SearchBar.SetText("")
			ui.Ui.App.SetFocus(ui.Ui.ExpandedView)
			views.SetCurrentView(views.FView)
		}
	})

	ui.Ui.ExpandedView.SetDoneFunc(func(e tcell.Key) {
		if e == tcell.KeyEscape {
			if views.GetCurrentView().GetViewName() == "BuffSearchView" {
				views.SetCurrentView(views.FView)
				ui.Ui.SearchBar.SetText("")
				client.Matches = nil
			}
		}
	})

	ui.Ui.SearchBar.SetChangedFunc(func(text string) {
		if views.GetCurrentView().GetViewName() == "BuffSearchView" {
			var f client.FileNodes = client.DirTree.Children
			client.Matches = fuzzy.FindFrom(text, f)
			views.BuffSView.Update(ui.Ui.ExpandedView)
		}
	})

	go func() {
		for {
			ui.Ui.App.Draw()
			time.Sleep(time.Second)
		}
	}()

	if err := ui.Ui.App.Run(); err != nil {
		panic(err)
	}
	defer watchers.DBCheck()

}
