package main

import (
	"strconv"
	"time"

	"github.com/aditya-K2/gomp/cache"
	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/gomp/config"
	"github.com/aditya-K2/gomp/notify"
	"github.com/aditya-K2/gomp/render"
	"github.com/aditya-K2/gomp/ui"
	"github.com/aditya-K2/gomp/utils"
	"github.com/aditya-K2/gomp/views"

	"github.com/aditya-K2/fuzzy"
	"github.com/fhs/gompd/v2/mpd"
	"github.com/gdamore/tcell/v2"
	"github.com/spf13/viper"
)

func main() {
	config.ReadConfig()
	var mpdConnectionError error
	del := ""
	nt := viper.GetString("NETWORK_TYPE")
	port := viper.GetString("MPD_PORT")
	if nt == "tcp" {
		del = ":"
	} else if nt == "unix" && port != "" {
		port = ""
	}
	client.Conn, mpdConnectionError = mpd.Dial(nt,
		viper.GetString("NETWORK_ADDRESS")+del+port)
	if mpdConnectionError != nil {
		utils.Print("RED", "Could Not Connect to MPD Server\n")
		utils.Print("GREEN", "Make Sure You Mention the Correct MPD Port in the config file.\n")
		panic(mpdConnectionError)
	}
	CONN := client.Conn
	defer CONN.Close()

	ui.SetConnection(CONN)

	cache.SetCacheDir(viper.GetString("CACHE_DIR"))

	render.Rendr = render.NewRenderer()
	// Connecting the Renderer to the Main UI
	ui.ConnectRenderer(render.Rendr)

	ui.Ui = ui.NewApplication()

	fileMap, err := CONN.ListAllInfo("/")
	if err != nil {
		utils.Print("RED", "Could Not Generate the File Map\n")
		utils.Print("GREEN", "Make Sure You Mention the Correct MPD Port in the config file.\n")
		panic(err)
	}

	// Generating the Directory Tree for File Navigation.
	client.DirTree = client.GenerateDirectoryTree(fileMap)

	// Default View upon Opening is of Playlist.
	views.PView.Update(ui.Ui.ExpandedView)

	var Volume int64
	var Random, Repeat bool
	var SeekOffset = viper.GetInt("SEEK_OFFSET")

	if _v, err := CONN.Status(); err != nil {
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

	notify.Notify = notify.NewNotificationServer()
	notify.Notify.Start()

	if c, err := CONN.CurrentSong(); err != nil {
		utils.Print("RED", "Could Not Retrieve the Current Song\n")
		panic(err)
	} else {
		if len(c) != 0 {
			render.Rendr.Start(c["file"])
		} else {
			render.Rendr.Start("stop")
		}
	}

	// This Function Is Responsible for Changing the Focus it uses the Focus Map and Based on it Chooses
	// the Draw Function
	views.SetCurrentView(views.PView)
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
				notify.Notify.Send("Could not Toggle Play Back")
			}
		},
		"showParentContent": func() {
			views.GetCurrentView().ShowParentContent()
		},
		"nextSong": func() {
			if err := CONN.Next(); err != nil {
				notify.Notify.Send("Could not Select the Next Song")
			}
		},
		"clearPlaylist": func() {
			if err := CONN.Clear(); err != nil {
				notify.Notify.Send("Could not Clear the Playlist")
			} else {
				notify.Notify.Send("Playlist Cleared")
			}
		},
		"previousSong": func() {
			if err := CONN.Previous(); err != nil {
				notify.Notify.Send("Could Not Select the Previous Song")
			}
		},
		"addToPlaylist": func() {
			views.GetCurrentView().AddToPlaylist()
		},
		"toggleRandom": func() {
			if err := CONN.Random(!Random); err == nil {
				Random = !Random
			}
		},
		"toggleRepeat": func() {
			if err := CONN.Repeat(!Repeat); err == nil {
				Repeat = !Repeat
			}
		},
		"decreaseVolume": func() {
			if Volume <= 0 {
				Volume = 0
			} else {
				Volume -= 10
			}
			if err := CONN.SetVolume(int(Volume)); err != nil {
				notify.Notify.Send("Could Not Decrease the Volume")
			}
		},
		"increaseVolume": func() {
			if Volume >= 100 {
				Volume = 100
			} else {
				Volume += 10
			}
			if err := CONN.SetVolume(int(Volume)); err != nil {
				notify.Notify.Send("Could Not Increase the Volume")
			}
		},
		"navigateToFiles": func() {
			views.SetCurrentView(views.FView)
			ui.Ui.Navbar.Select(1, 0)
			views.FView.Update(ui.Ui.ExpandedView)
		},
		"navigateToPlaylist": func() {
			views.SetCurrentView(views.PView)
			ui.Ui.Navbar.Select(0, 0)
			views.PView.Update(ui.Ui.ExpandedView)
		},
		"navigateToMostPlayed": func() {
			ui.Ui.Navbar.Select(2, 0)
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
			if err := CONN.Stop(); err != nil {
				notify.Notify.Send("Could not Stop the Playback")
			} else {
				notify.Notify.Send("Playback Stopped")
			}
		},
		"updateDB": func() {
			_, err = CONN.Update("")
			if err != nil {
				notify.Notify.Send("Could Not Update the Database")
			} else {
				notify.Notify.Send("Database Updated")
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
			if err := CONN.SeekCur(time.Second*time.Duration(SeekOffset), true); err != nil {
				notify.Notify.Send("Could Not Seek Forward in the Song")
			}
		},
		"SeekBackward": func() {
			if err := CONN.SeekCur(-1*time.Second*time.Duration(SeekOffset), true); err != nil {
				notify.Notify.Send("Could Not Seek Backward in the Song")
			}
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
					if p, err := CONN.PlaylistInfo(-1, -1); err != nil {
						notify.Notify.Send("Error Getting PlaylistInfo")
					} else {
						if len(p) == 0 {
							notify.Notify.Send("Empty Playlist")
							return nil
						}
					}
				}
			} else if views.GetCurrentView().GetViewName() == "SearchView" {
				if e.Rune() == 'j' || e.Rune() == 'k' {
					if client.SearchContentSlice == nil || len(client.SearchContentSlice) == 0 {
						notify.Notify.Send("No Search Results")
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
					notify.Notify.Send("Could Not Retrieve the Results")
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
}
