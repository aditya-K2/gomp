package views

import (
	"github.com/aditya-K2/gomp/globals"
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
	UI := globals.Ui
	CONN := globals.Conn
	r, _ := UI.ExpandedView.GetSelection()
	if err := CONN.Play(r); err != nil {
		globals.Notify.Send("Could Not Play the Song")
		return
	}
}

func (s PlaylistView) ShowParentContent() {
	globals.Notify.Send("Not Allowed in this View")
	return
}
func (p PlaylistView) AddToPlaylist() {}

func (p PlaylistView) Quit() {
	globals.Ui.App.Stop()
}

func (p PlaylistView) FocusBuffSearchView() {}

func (p PlaylistView) DeleteSongFromPlaylist() {
	UI := globals.Ui
	CONN := globals.Conn
	r, _ := UI.ExpandedView.GetSelection()
	if err := CONN.Delete(r, -1); err != nil {
		globals.Notify.Send("Could not Remove the Song from Playlist")
	}
}

func (p PlaylistView) Update(inputTable *tview.Table) {
	CONN := globals.Conn
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
