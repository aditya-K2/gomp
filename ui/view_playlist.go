package ui

import (
	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/utils"
	"github.com/aditya-K2/tview"
	"github.com/fhs/gompd/v2/mpd"
)

type PlaylistView struct {
	Playlist []mpd.Attrs
}

func (s PlaylistView) Name() string {
	return "PlaylistView"
}

func (p PlaylistView) ShowChildrenContent() {
	CONN := client.Conn
	r, _ := Ui.MainS.GetSelection()
	if err := CONN.Play(r); err != nil {
		SendNotification("Could Not Play the Song")
		return
	}
}

func (s PlaylistView) ShowParentContent() {
	SendNotification("Not Allowed in this View")
	return
}
func (p PlaylistView) AddToPlaylist() {}

func (p PlaylistView) Quit() {
	Ui.App.Stop()
}

func (p PlaylistView) FocusBuffSearchView() {}

func (p *PlaylistView) DeleteSongFromPlaylist() {
	CONN := client.Conn
	r, _ := Ui.MainS.GetSelection()
	if err := CONN.Delete(r, -1); err != nil {
		SendNotification("Could not Remove the Song from Playlist")
	} else {
		if p.Playlist, err = client.Conn.PlaylistInfo(-1, -1); err != nil {
			utils.Print("RED", "Couldn't get the current Playlist.\n")
			panic(err)
		}
	}

}

func (p PlaylistView) Update(inputTable *tview.Table) {
	inputTable.Clear()
	for i, j := range p.Playlist {
		_, _, w, _ := inputTable.GetInnerRect()
		if j["Title"] == "" || j["Artist"] == "" || j["Album"] == "" {
			inputTable.SetCell(i, 0,
				GetCell(
					utils.GetFormattedString(j["file"], w/3), clr.File))

		} else {
			inputTable.SetCell(i, 0,
				GetCell(
					utils.GetFormattedString(j["Title"], w/3), clr.Track))
			inputTable.SetCell(i, 1,
				GetCell(
					utils.GetFormattedString(j["Artist"], w/3), clr.Artist))
			inputTable.SetCell(i, 2,
				GetCell(j["Album"], clr.Album))
		}
	}
}
