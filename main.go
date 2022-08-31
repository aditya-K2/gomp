package main

import (
	"strconv"
	"time"

	"github.com/aditya-K2/gomp/cache"
	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/gomp/config"
	"github.com/aditya-K2/gomp/globals"
	"github.com/aditya-K2/gomp/notify"
	"github.com/aditya-K2/gomp/render"
	"github.com/aditya-K2/gomp/ui"
	"github.com/aditya-K2/gomp/utils"
	"github.com/aditya-K2/gomp/views"

	"github.com/aditya-K2/fuzzy"
	"github.com/fhs/gompd/mpd"
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
	globals.Conn, mpdConnectionError = mpd.Dial(nt,
		viper.GetString("NETWORK_ADDRESS")+del+port)
	if mpdConnectionError != nil {
		utils.Print("RED", "Could Not Connect to MPD Server\n")
		utils.Print("GREEN", "Make Sure You Mention the Correct MPD Port in the config file.\n")
		panic(mpdConnectionError)
	}
	CONN := globals.Conn
	defer CONN.Close()

	client.SetConnection(CONN)
	ui.SetConnection(CONN)
	render.SetConnection(CONN)

	cache.SetCacheDir(viper.GetString("CACHE_DIR"))

	globals.Renderer = render.NewRenderer()
	// Connecting the Renderer to the Main UI
	ui.ConnectRenderer(globals.Renderer)

	globals.Ui = ui.NewApplication()

	// Connecting the Notification Server to the Main UI
	notify.ConnectUI(globals.Ui)

	fileMap, err := CONN.ListAllInfo("/")
	if err != nil {
		utils.Print("RED", "Could Not Generate the File Map\n")
		utils.Print("GREEN", "Make Sure You Mention the Correct MPD Port in the config file.\n")
		panic(err)
	}

	// Generating the Directory Tree for File Navigation.
	globals.DirTree = client.GenerateDirectoryTree(fileMap)

	// Default View upon Opening is of Playlist.
	views.PView.Update(globals.Ui.ExpandedView)

	var Volume int64
	var Random, Repeat bool

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

	Notify := notify.NewNotificationServer()
	Notify.Start()

	if c, err := CONN.CurrentSong(); err != nil {
		utils.Print("RED", "Could Not Retrieve the Current Song\n")
		panic(err)
	} else {
		if len(c) != 0 {
			globals.Renderer.Start(c["file"])
		} else {
			globals.Renderer.Start("stop")
		}
	}

	// Connecting Notification Server to Client and Rendering Module so that they can send Notifications
	client.SetNotificationServer(Notify)
	render.SetNotificationServer(Notify)

	// This Function Is Responsible for Changing the Focus it uses the Focus Map and Based on it Chooses
	// the Draw Function
	views.SetCurrentView(views.PView)
	globals.Ui.ExpandedView.SetDrawFunc(func(s tcell.Screen, x, y, width, height int) (int, int, int, int) {
		views.GetCurrentView().Update(globals.Ui.ExpandedView)
		return globals.Ui.ExpandedView.GetInnerRect()
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
				Notify.Send("Could not Toggle Play Back")
			}
		},
		"showParentContent": func() {
			views.GetCurrentView().ShowParentContent()
		},
		"nextSong": func() {
			if err := CONN.Next(); err != nil {
				Notify.Send("Could not Select the Next Song")
			}
		},
		"clearPlaylist": func() {
			if err := CONN.Clear(); err != nil {
				Notify.Send("Could not Clear the Playlist")
			} else {
				Notify.Send("Playlist Cleared")
			}
		},
		"previousSong": func() {
			if err := CONN.Previous(); err != nil {
				Notify.Send("Could Not Select the Previous Song")
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
				Notify.Send("Could Not Decrease the Volume")
			}
		},
		"increaseVolume": func() {
			if Volume >= 100 {
				Volume = 100
			} else {
				Volume += 10
			}
			if err := CONN.SetVolume(int(Volume)); err != nil {
				Notify.Send("Could Not Increase the Volume")
			}
		},
		"navigateToFiles": func() {
			views.SetCurrentView(views.FView)
			globals.Ui.Navbar.Select(1, 0)
			views.FView.Update(globals.Ui.ExpandedView)
		},
		"navigateToPlaylist": func() {
			views.SetCurrentView(views.PView)
			globals.Ui.Navbar.Select(0, 0)
			views.PView.Update(globals.Ui.ExpandedView)
		},
		"navigateToMostPlayed": func() {
			globals.Ui.Navbar.Select(2, 0)
		},
		"navigateToSearch": func() {
			views.SetCurrentView(views.SView)
			globals.Ui.Navbar.Select(3, 0)
			views.SView.Update(globals.Ui.ExpandedView)
		},
		"quit": func() {
			views.GetCurrentView().Quit()
		},
		"stop": func() {
			if err := CONN.Stop(); err != nil {
				Notify.Send("Could not Stop the Playback")
			} else {
				Notify.Send("Playback Stopped")
			}
		},
		"updateDB": func() {
			_, err = CONN.Update("")
			if err != nil {
				Notify.Send("Could Not Update the Database")
			} else {
				Notify.Send("Database Updated")
			}
		},
		"deleteSongFromPlaylist": func() {
			views.GetCurrentView().DeleteSongFromPlaylist()
		},
		"FocusSearch": func() {
			globals.Ui.App.SetFocus(globals.Ui.SearchBar)
		},
		"FocusBuffSearch": func() {
			views.GetCurrentView().FocusBuffSearchView()
		},
	}

	// Generating the Key Map Based on the Function Map Here Basically the Values will be flipped
	// In the config if togglePlayBack is mapped to [ T , P, SPACE ] then here Basically we will receive a map
	// for each event T, P, SPACE mapped to the same function togglePlayBack
	config.GenerateKeyMap(FuncMap)

	globals.Ui.SearchBar.SetAutocompleteFunc(func(c string) []string {
		if views.GetCurrentView().GetViewName() == "BuffSearchView" {
			return nil
		} else {
			if c != "" && c != " " && c != "  " {
				_, _, w, _ := globals.Ui.SearchBar.GetRect()
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
	globals.Ui.ExpandedView.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if val, ok := config.KEY_MAP[int(e.Rune())]; ok {
			FuncMap[val]()
			return nil
		} else {
			if views.GetCurrentView().GetViewName() == "PlaylistView" {
				if e.Rune() == 'j' || e.Rune() == 'k' {
					if p, err := CONN.PlaylistInfo(-1, -1); err != nil {
						Notify.Send("Error Getting PlaylistInfo")
					} else {
						if len(p) == 0 {
							Notify.Send("Empty Playlist")
							return nil
						}
					}
				}
			}
			return e
		}
	})

	globals.Ui.SearchBar.SetDoneFunc(func(e tcell.Key) {
		if e == tcell.KeyEnter {
			globals.Ui.ExpandedView.Select(0, 0)
			if views.GetCurrentView().GetViewName() == "BuffSearchView" {
				globals.Ui.App.SetFocus(globals.Ui.ExpandedView)
			} else {
				views.SetCurrentView(views.SView)
				globals.SearchContentSlice = nil
				globals.SearchContentSlice, err = client.GenerateContentSlice(globals.Ui.SearchBar.GetText())
				if err != nil {
					Notify.Send("Could Not Retrieve the Results")
				} else {
					globals.Ui.SearchBar.SetText("")
					globals.Ui.App.SetFocus(globals.Ui.ExpandedView)
					globals.Ui.Navbar.Select(3, 0)
				}
			}
		}
		if e == tcell.KeyEscape {
			if views.GetCurrentView().GetViewName() == "SearchView" {
			} else if views.GetCurrentView().GetViewName() == "BuffSearchView" {
				views.SetCurrentView(views.FView)
				globals.Matches = nil
			}
			globals.Ui.SearchBar.SetText("")
			globals.Ui.App.SetFocus(globals.Ui.ExpandedView)
		}
	})

	globals.Ui.ExpandedView.SetDoneFunc(func(e tcell.Key) {
		if e == tcell.KeyEscape {
			if views.GetCurrentView().GetViewName() == "BuffSearchView" {
				views.SetCurrentView(views.FView)
				globals.Ui.SearchBar.SetText("")
				globals.Matches = nil
			}
		}
	})

	globals.Ui.SearchBar.SetChangedFunc(func(text string) {
		if views.GetCurrentView().GetViewName() == "BuffSearchView" {
			var f client.FileNodes = globals.DirTree.Children
			globals.Matches = fuzzy.FindFrom(text, f)
			views.BuffSView.Update(globals.Ui.ExpandedView)
		}
	})

	go func() {
		for {
			globals.Ui.App.Draw()
			time.Sleep(time.Second)
		}
	}()

	if err := globals.Ui.App.Run(); err != nil {
		panic(err)
	}
}
