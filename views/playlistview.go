package views

import (
	"github.com/aditya-K2/gomp/client"
	"github.com/aditya-K2/gomp/notify"
	"github.com/aditya-K2/gomp/ui"
	"github.com/aditya-K2/gomp/utils"
	"github.com/aditya-K2/tview"
	"github.com/gdamore/tcell/v2"
)

type PlaylistView struct {
}

func (s PlaylistView) GetViewName() string {
	return "PlaylistView"
}

func GetCell(text string, foreground tcell.Color, bold bool) *tview.TableCell {
	return tview.NewTableCell(text).
		SetAlign(tview.AlignLeft).
		SetStyle(tcell.StyleDefault.
			Foreground(foreground).
			Background(tcell.ColorBlack).
			Bold(bold))
}

func (p PlaylistView) ShowChildrenContent() {
	UI := ui.Ui
	CONN := client.Conn
	r, _ := UI.ExpandedView.GetSelection()
	if err := CONN.Play(r); err != nil {
		notify.Notify.Send("Could Not Play the Song")
		return
	}
}

func (s PlaylistView) ShowParentContent() {
	notify.Notify.Send("Not Allowed in this View")
	return
}
func (p PlaylistView) AddToPlaylist() {}

func (p PlaylistView) Quit() {
	ui.Ui.App.Stop()
}

func (p PlaylistView) FocusBuffSearchView() {}

func (p PlaylistView) DeleteSongFromPlaylist() {
	UI := ui.Ui
	CONN := client.Conn
	r, _ := UI.ExpandedView.GetSelection()
	if err := CONN.Delete(r, -1); err != nil {
		notify.Notify.Send("Could not Remove the Song from Playlist")
	}
}

func (p PlaylistView) Update(inputTable *tview.Table) {
	CONN := client.Conn
	_playlistAttr, _ := CONN.PlaylistInfo(-1, -1)

	inputTable.Clear()
	for i, j := range _playlistAttr {
		_, _, w, _ := inputTable.GetInnerRect()
		if j["Title"] == "" || j["Artist"] == "" || j["Album"] == "" {
			inputTable.SetCell(i, 0,
				GetCell(
					utils.GetFormattedString(j["file"], w/3), tcell.ColorBlue, true))

		} else {
			inputTable.SetCell(i, 0,
				GetCell(
					utils.GetFormattedString(j["Title"], w/3), tcell.ColorGreen, false))
			inputTable.SetCell(i, 1,
				GetCell(
					utils.GetFormattedString(j["Artist"], w/3), tcell.ColorPurple, false))
			inputTable.SetCell(i, 2,
				GetCell(j["Album"], tcell.ColorYellow, false))
		}
	}
}
