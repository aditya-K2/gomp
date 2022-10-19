package watchers

import (
	"time"

	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/gomp/notify"
	"github.com/aditya-K2/gomp/render"
	"github.com/aditya-K2/gomp/ui"
	"github.com/aditya-K2/gomp/utils"
	"github.com/aditya-K2/gomp/views"
	"github.com/fhs/gompd/v2/mpd"
)

var (
	currentSong mpd.Attrs
)

func Init() {
	if c, err := client.Conn.CurrentSong(); err != nil {
		notify.Notify.Send("Couldn't get current song from MPD")
	} else {
		currentSong = c
	}
}

func DrawCover(c mpd.Attrs) {
	if len(c) == 0 {
		render.Rendr.Send("stop")
	} else {
		render.Rendr.Send(c["file"])
	}
}

func StartRectWatcher() {
	go func() {
		for {
			ImgX, ImgY, ImgW, ImgH := ui.Ui.ImagePreviewer.GetRect()
			if ImgX != ui.ImgX || ImgY != ui.ImgY ||
				ImgW != ui.ImgW || ImgH != ui.ImgH {
				ui.ImgX = ImgX
				ui.ImgY = ImgY
				ui.ImgW = ImgW
				ui.ImgH = ImgH
				DrawCover(currentSong)
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
}

func StartPlaylistWatcher() {
	var err error
	if views.PView.Playlist == nil {
		if views.PView.Playlist, err = client.Conn.PlaylistInfo(-1, -1); err != nil {
			utils.Print("RED", "Watcher couldn't get the current Playlist.\n")
			panic(err)
		}
	}

	nt, addr := utils.GetNetwork()
	w, err := mpd.NewWatcher(nt, addr, "", "playlist")
	if err != nil {
		utils.Print("RED", "Could Not Start Watcher.\n")
		utils.Print("GREEN", "Please check your MPD Info in config File.\n")
		panic(err)
	}

	go func() {
		for err := range w.Error {
			notify.Notify.Send(err.Error())
		}
	}()

	go func() {
		for subsystem := range w.Event {
			if subsystem == "playlist" {
				if views.PView.Playlist, err = client.Conn.PlaylistInfo(-1, -1); err != nil {
					utils.Print("RED", "Watcher couldn't get the current Playlist.\n")
					panic(err)
				}
			} else if subsystem == "player" {
				if c, cerr := client.Conn.CurrentSong(); cerr != nil {
					utils.Print("RED", "Watcher couldn't get the current Playlist.\n")
					panic(err)
				} else {
					currentSong = c
				}
			}
		}
	}()
}
