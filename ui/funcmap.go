package ui

import (
	"strconv"
	"time"

	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/gomp/config"
	"github.com/aditya-K2/gomp/utils"
	"github.com/fhs/gompd/v2/mpd"
)

func GenerateFuncMap(Conn *mpd.Client) map[string]func() {
	var Volume int64
	var Random, Repeat bool
	var SeekOffset = config.Config.SeekOffset
	var SeekFunc = func(back bool) {
		if status, err := Conn.Status(); err != nil {
			SendNotification("Could not get MPD Status")
		} else {
			if status["state"] == "play" {
				var stime time.Duration
				if back {
					stime = -1 * time.Second * time.Duration(SeekOffset)
				} else {
					stime = time.Second * time.Duration(SeekOffset)
				}
				if err := Conn.SeekCur(stime, true); err != nil {
					SendNotification("Could Not Seek Forward in the Song")
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

	// Function Maps is used For Mapping Keys According to the Value mapped to the Key the respective Function is called
	// For e.g. in the config if the User Maps T to togglePlayBack then whenever in the input handler the T is received
	// the respective function in this case togglePlayBack is called.
	var funcMap = map[string]func(){
		"showChildrenContent": func() {
			GetCurrentView().ShowChildrenContent()
		},
		"togglePlayBack": func() {
			if err := client.TogglePlayBack(); err != nil {
				SendNotification("Could not Toggle Play Back")
			}
		},
		"showParentContent": func() {
			GetCurrentView().ShowParentContent()
		},
		"nextSong": func() {
			if err := Conn.Next(); err != nil {
				SendNotification("Could not Select the Next Song")
			}
		},
		"clearPlaylist": func() {
			if err := Conn.Clear(); err != nil {
				SendNotification("Could not Clear the Playlist")
			} else {
				if PView.Playlist, err = client.Conn.PlaylistInfo(-1, -1); err != nil {
					utils.Print("RED", "Couldn't get the current Playlist.\n")
				} else {
					SendNotification("Playlist Cleared!")
				}
			}
		},
		"previousSong": func() {
			if err := Conn.Previous(); err != nil {
				SendNotification("Could Not Select the Previous Song")
			}
		},
		"addToPlaylist": func() {
			GetCurrentView().AddToPlaylist()
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
				SendNotification("Could Not Decrease the Volume")
			}
		},
		"increaseVolume": func() {
			if Volume >= 100 {
				Volume = 100
			} else {
				Volume += 10
			}
			if err := Conn.SetVolume(int(Volume)); err != nil {
				SendNotification("Could Not Increase the Volume")
			}
		},
		"navigateToFiles": func() {
			SetCurrentView(FView)
			Ui.Navbar.Select(1, 0)
			FView.Update(Ui.MainS)
		},
		"navigateToPlaylist": func() {
			SetCurrentView(&PView)
			Ui.Navbar.Select(0, 0)
			PView.Update(Ui.MainS)
		},
		"navigateToMostPlayed": func() {
			SetCurrentView(&MPView)
			Ui.Navbar.Select(2, 0)
			MPView.Update(Ui.MainS)
		},
		"navigateToSearch": func() {
			SetCurrentView(SView)
			Ui.Navbar.Select(3, 0)
			SView.Update(Ui.MainS)
		},
		"quit": func() {
			GetCurrentView().Quit()
		},
		"stop": func() {
			if err := Conn.Stop(); err != nil {
				SendNotification("Could not Stop the Playback")
			} else {
				SendNotification("Playback Stopped")
			}
		},
		"updateDB": func() {
			if _, err := Conn.Update(""); err != nil {
				SendNotification("Could Not Update the Database")
			} else {
				SendNotification("Database Updated")
			}
		},
		"deleteSongFromPlaylist": func() {
			GetCurrentView().DeleteSongFromPlaylist()
		},
		"FocusSearch": func() {
			Ui.App.SetFocus(Ui.SearchBar)
		},
		"FocusBuffSearch": func() {
			GetCurrentView().FocusBuffSearchView()
		},
		"SeekForward": func() {
			SeekFunc(false)
		},
		"SeekBackward": func() {
			SeekFunc(true)
		},
	}
	return funcMap
}
