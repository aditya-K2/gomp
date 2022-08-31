package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aditya-K2/gomp/cache"
	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/gomp/config"
	"github.com/aditya-K2/gomp/notify"
	"github.com/aditya-K2/gomp/render"
	"github.com/aditya-K2/gomp/ui"
	"github.com/aditya-K2/gomp/utils"

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
	CONN, mpdConnectionError := mpd.Dial(nt,
		viper.GetString("NETWORK_ADDRESS")+del+port)
	if mpdConnectionError != nil {
		utils.Print("RED", "Could Not Connect to MPD Server\n")
		utils.Print("GREEN", "Make Sure You Mention the Correct MPD Port in the config file.\n")
		panic(mpdConnectionError)
	}
	defer CONN.Close()

	ui.GenerateFocusMap()

	client.SetConnection(CONN)
	ui.SetConnection(CONN)
	render.SetConnection(CONN)

	cache.SetCacheDir(viper.GetString("CACHE_DIR"))

	Renderer := render.NewRenderer()
	// Connecting the Renderer to the Main UI
	ui.ConnectRenderer(Renderer)

	UI := ui.NewApplication()

	// Connecting the Notification Server to the Main UI
	notify.ConnectUI(UI)

	fileMap, err := CONN.ListAllInfo("/")
	if err != nil {
		utils.Print("RED", "Could Not Generate the File Map\n")
		utils.Print("GREEN", "Make Sure You Mention the Correct MPD Port in the config file.\n")
		panic(err)
	}

	// Generating the Directory Tree for File Navigation.
	dirTree := client.GenerateDirectoryTree(fileMap)

	// Default View upon Opening is of Playlist.
	client.UpdatePlaylist(UI.ExpandedView)

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
			Renderer.Start(c["file"])
		} else {
			Renderer.Start("stop")
		}
	}

	// Connecting Notification Server to Client and Rendering Module so that they can send Notifications
	client.SetNotificationServer(Notify)
	render.SetNotificationServer(Notify)

	// This is the Slice that is used to Display Content in the SearchView
	var SearchContentSlice []interface{}
	var Matches fuzzy.Matches

	// This Function Is Responsible for Changing the Focus it uses the Focus Map and Based on it Chooses
	// the Draw Function
	UI.ExpandedView.SetDrawFunc(func(s tcell.Screen, x, y, width, height int) (int, int, int, int) {
		if ui.HasFocus("Playlist") {
			client.UpdatePlaylist(UI.ExpandedView)
		} else if ui.HasFocus("SearchView") {
			client.UpdateSearchView(UI.ExpandedView, SearchContentSlice)
		} else if ui.HasFocus("FileBrowser") {
			client.Update(dirTree.Children, UI.ExpandedView)
		} else if ui.HasFocus("BuffSearchView") {
			client.UpdateBuffSearchView(UI.ExpandedView, Matches, dirTree.Children)
		}
		return UI.ExpandedView.GetInnerRect()
	})

	// Function Maps is used For Mapping Keys According to the Value mapped to the Key the respective Function is called
	// For e.g. in the config if the User Maps T to togglePlayBack then whenever in the input handler the T is received
	// the respective function in this case togglePlayBack is called.
	var FuncMap = map[string]func(){
		"showChildrenContent": func() {
			if ui.HasFocus("FileBrowser") {
				r, _ := UI.ExpandedView.GetSelection()
				ui.SetFocus("FileBrowser")
				if len(dirTree.Children[r].Children) == 0 {
					if id, err := CONN.AddId(dirTree.Children[r].AbsolutePath, -1); err != nil {
						Notify.Send(fmt.Sprintf("Could not Add Song %s",
							dirTree.Children[r].Path))
					} else {
						if err := CONN.PlayId(id); err != nil {
							Notify.Send(fmt.Sprintf("Could Not Play Song %s",
								dirTree.Children[r].Path))
						}
					}
				} else {
					client.Update(dirTree.Children[r].Children, UI.ExpandedView)
					dirTree = &dirTree.Children[r]
					UI.ExpandedView.Select(0, 0)
				}
			} else if ui.HasFocus("Playlist") {
				r, _ := UI.ExpandedView.GetSelection()
				if err := CONN.Play(r); err != nil {
					Notify.Send("Could Not Play the Song")
				}
			} else if ui.HasFocus("SearchView") {
				r, _ := UI.ExpandedView.GetSelection()
				client.AddToPlaylist(SearchContentSlice[r], true)
			} else if ui.HasFocus("BuffSearchView") {
				r, _ := UI.ExpandedView.GetSelection()
				ui.SetFocus("FileBrowser")
				if len(dirTree.Children[r].Children) == 0 {
					if id, err := CONN.AddId(dirTree.Children[Matches[r].Index].AbsolutePath, -1); err != nil {
						Notify.Send(fmt.Sprintf("Could Not add the Song %s to the Playlist",
							dirTree.Children[Matches[r].Index].AbsolutePath))
					} else {
						if err := CONN.PlayId(id); err != nil {
							Notify.Send("Could not Play the Song")
						}
					}
				} else {
					client.Update(dirTree.Children[Matches[r].Index].Children, UI.ExpandedView)
					dirTree = &dirTree.Children[Matches[r].Index]
				}
				UI.SearchBar.SetText("")
				// Resetting Matches
				Matches = nil
			}
		},
		"togglePlayBack": func() {
			if err := client.TogglePlayBack(); err != nil {
				Notify.Send("Could not Toggle Play Back")
			}
		},
		"showParentContent": func() {
			if ui.HasFocus("FileBrowser") {
				if dirTree.Parent != nil {
					client.Update(dirTree.Parent.Children, UI.ExpandedView)
					dirTree = dirTree.Parent
				}
			} else {
				Notify.Send("Not Allowed in this View")
				return
			}
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
			if ui.HasFocus("FileBrowser") {
				r, _ := UI.ExpandedView.GetSelection()
				if err := CONN.Add(dirTree.Children[r].AbsolutePath); err != nil {
					Notify.Send(fmt.Sprintf("Could not add %s to the Playlist",
						dirTree.Children[r].Path))
				}
			} else if ui.HasFocus("SearchView") {
				r, _ := UI.ExpandedView.GetSelection()
				client.AddToPlaylist(SearchContentSlice[r], false)
			} else if ui.HasFocus("BuffSearchView") {
				r, _ := UI.ExpandedView.GetSelection()
				if err := CONN.Add(dirTree.Children[Matches[r].Index].AbsolutePath); err != nil {
					Notify.Send(fmt.Sprintf("Could Not Add URI %s to the Playlist",
						dirTree.Children[Matches[r].Index].Path))
				} else {
					ui.SetFocus("FileBrowser")
					Notify.Send(fmt.Sprintf("URI Added %s to the Playlist",
						dirTree.Children[Matches[r].Index].Path))
					ui.SetFocus("BuffSearchView")
				}
			}
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
			ui.SetFocus("FileBrowser")
			UI.Navbar.Select(1, 0)
			client.Update(dirTree.Children, UI.ExpandedView)
		},
		"navigateToPlaylist": func() {
			ui.SetFocus("Playlist")
			UI.Navbar.Select(0, 0)
			client.UpdatePlaylist(UI.ExpandedView)
		},
		"navigateToMostPlayed": func() {
			UI.Navbar.Select(2, 0)
		},
		"navigateToSearch": func() {
			ui.SetFocus("SearchView")
			UI.Navbar.Select(3, 0)
		},
		"quit": func() {
			if ui.HasFocus("BuffSearchView") {
				ui.SetFocus("FileBrowser")
				UI.SearchBar.SetText("")
				Matches = nil
			} else {
				UI.App.Stop()
			}
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
			if ui.HasFocus("Playlist") {
				r, _ := UI.ExpandedView.GetSelection()
				if err := CONN.Delete(r, -1); err != nil {
					Notify.Send("Could not Remove the Song from Playlist")
				}
			}
		},
		"FocusSearch": func() {
			UI.App.SetFocus(UI.SearchBar)
		},
		"FocusBuffSearch": func() {
			if ui.HasFocus("FileBrowser") || ui.HasFocus("BuffSearchView") {
				ui.SetFocus("BuffSearchView")
				UI.App.SetFocus(UI.SearchBar)
			}
		},
	}

	// Generating the Key Map Based on the Function Map Here Basically the Values will be flipped
	// In the config if togglePlayBack is mapped to [ T , P, SPACE ] then here Basically we will receive a map
	// for each event T, P, SPACE mapped to the same function togglePlayBack
	config.GenerateKeyMap(FuncMap)

	UI.SearchBar.SetAutocompleteFunc(func(c string) []string {
		if ui.HasFocus("BuffSearchView") {
			return nil
		} else {
			if c != "" && c != " " && c != "  " {
				_, _, w, _ := UI.SearchBar.GetRect()
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
	UI.ExpandedView.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		if val, ok := config.KEY_MAP[int(e.Rune())]; ok {
			FuncMap[val]()
			return nil
		} else {
			if ui.HasFocus("Playlist") {
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

	UI.SearchBar.SetDoneFunc(func(e tcell.Key) {
		if e == tcell.KeyEnter {
			UI.ExpandedView.Select(0, 0)
			if ui.HasFocus("BuffSearchView") {
				UI.App.SetFocus(UI.ExpandedView)
			} else {
				ui.SetFocus("SearchView")
				SearchContentSlice = nil
				SearchContentSlice, err = client.GenerateContentSlice(UI.SearchBar.GetText())
				if err != nil {
					Notify.Send("Could Not Retrieve the Results")
				} else {
					UI.SearchBar.SetText("")
					UI.App.SetFocus(UI.ExpandedView)
					UI.Navbar.Select(3, 0)
				}
			}
		}
		if e == tcell.KeyEscape {
			if ui.HasFocus("SearchView") {
				ui.FocusMap["SearchView"] = false
			} else if ui.HasFocus("BuffSearchView") {
				ui.SetFocus("FileBrowser")
				Matches = nil
			}
			UI.SearchBar.SetText("")
			UI.App.SetFocus(UI.ExpandedView)
		}
	})

	UI.ExpandedView.SetDoneFunc(func(e tcell.Key) {
		if e == tcell.KeyEscape {
			if ui.HasFocus("BuffSearchView") {
				ui.SetFocus("FileBrowser")
				UI.SearchBar.SetText("")
				Matches = nil
			}
		}
	})

	UI.SearchBar.SetChangedFunc(func(text string) {
		if ui.HasFocus("BuffSearchView") {
			var f client.FileNodes = dirTree.Children
			Matches = fuzzy.FindFrom(text, f)
			client.UpdateBuffSearchView(UI.ExpandedView, Matches, dirTree.Children)
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
